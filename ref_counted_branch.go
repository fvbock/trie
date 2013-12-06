package trie

import (
	"fmt"
	// "log"
	"strings"
	"sync"
	// "time"
)

// BRANCH

type MemberInfo struct {
	Value string
	Count int64
}

func (m *MemberInfo) String() string {
	return fmt.Sprintf("%s(%v)", m.Value, m.Count)
}

type RefCountBranch struct {
	sync.RWMutex
	Branches  map[byte]*RefCountBranch
	LeafValue []byte // tail end
	End       bool
	Count     int64
}

/*
NewRefCountBranch returns a new initialezed *RefCountBranch
*/
func (b *RefCountBranch) NewBranch() *RefCountBranch {
	return &RefCountBranch{
		Branches: make(map[byte]*RefCountBranch),
		Count:    0,
	}
}

/*
Add adds an entry to the Branch
*/
func (b *RefCountBranch) add(entry []byte) (addedBranch *RefCountBranch) {
	if b.LeafValue == nil && len(b.Branches) == 0 {
		b.LeafValue = entry
		b.setEnd(true)
		addedBranch = b
		return
	}

	// something came in but we already have branches for it
	// so the tail was the current branches index but no value
	// to push. just mark the current idx position as End
	if len(b.LeafValue) == 0 && len(entry) == 0 {
		b.setEnd(true)
		addedBranch = b
		return
	}

	// check the overlap between the current LeafValue and the new entry
	newLeaf := func(LeafValue, newEntry []byte) (leaf []byte) {
		for li, b := range LeafValue {
			if li > len(newEntry)-1 {
				break
			}
			if b == newEntry[li] {
				leaf = append(leaf, b)
			} else {
				break
			}
		}
		return
	}(b.LeafValue, entry)

	newLeafLen := len(newLeaf)

	// the new leaf is smaller than the current leaf.
	// we will push the old leaf down the branch
	if newLeafLen < len(b.LeafValue) {
		tail := b.LeafValue[newLeafLen:]
		idx := tail[0]
		newBranch := b.NewBranch()
		newBranch.LeafValue = tail[1:]

		b.LeafValue = newLeaf
		newBranch.Branches, b.Branches = b.Branches, newBranch.Branches
		newBranch.End, b.End = b.End, newBranch.End
		if newBranch.End {
			if b.Count > 0 {
				newBranch.Count = b.Count
			} else {
				newBranch.Count = 1
			}
		} else {
			newBranch.Count = 0
		}
		if b.End {
			b.Count = 1
		} else {
			b.Count = 0
		}
		b.Branches[idx] = newBranch
	}

	// new leaf is smaller than the entry, which means there will be more stuff
	// that we need to push down
	if newLeafLen < len(entry) {
		tail := entry[newLeafLen:]
		idx := tail[0]

		// create new branch at idx if it does not exists yet
		if _, notPresent := b.Branches[idx]; !notPresent {
			b.Branches[idx] = b.NewBranch()
		}
		// check whether the idx itself marks an End $. if so add a new idx
		addedBranch = b.Branches[idx].add(tail[1:])
	} else {
		// if there is nothing else to be pushed down we just have to mark the
		// current branch as an end. this happens when you add a value that already
		// is covered by the index but this particular end had not been marked.
		// eg. you already have 'food' and 'foot' (shared LeafValue of 'foo') in
		// your index and now add 'foo'.
		b.setEnd(true)
		addedBranch = b
	}
	return addedBranch
}

/*
Members returns slice of all Members of the Branch prepended with `branchPrefix`
*/
func (b *RefCountBranch) members(branchPrefix []byte) (members []*MemberInfo) {
	if b.End {
		members = append(members, &MemberInfo{string(append(branchPrefix, b.LeafValue...)), b.Count})
	}
	for idx, br := range b.Branches {
		newPrefix := append(append(branchPrefix, b.LeafValue...), idx)
		members = append(members, br.members(newPrefix)...)
	}
	return
}

/*
prefixMembers returns a slice of all Members of the Branch matching the given prefix. The values returned are prepended with `branchPrefix`
*/
func (b *RefCountBranch) prefixMembers(branchPrefix []byte, searchPrefix []byte) (members []*MemberInfo) {
	leafLen := len(b.LeafValue)
	searchPrefixLen := len(searchPrefix)

	// if the searchPrefix is empty we want all members
	if searchPrefixLen == 0 {
		members = append(members, b.members(branchPrefix)...)
		return
	}

	// if the searchPrefix is shorter than the leaf we will add the LeafValue
	// if it is an End and a the searchPrefix matches
	// if searchPrefixLen < leafLen {
	if searchPrefixLen > leafLen {
		for idx, br := range b.Branches {
			// does it match the next byte?
			if idx == searchPrefix[leafLen] {
				newSearchPrefix := searchPrefix[leafLen+1:]
				members = append(members, br.prefixMembers(append(append(branchPrefix, b.LeafValue...), idx), newSearchPrefix)...)
			}
		}
	} else if searchPrefixLen == leafLen {
		for i, sb := range searchPrefix {
			if sb != b.LeafValue[i] {
				return
			}
		}
		members = append(members, b.members(branchPrefix)...)
	} else {
		if b.End {
			for i, sb := range searchPrefix {
				if sb != b.LeafValue[i] {
					return
				}
			}
			members = append(members, &MemberInfo{string(append(branchPrefix, b.LeafValue...)), b.Count})
		}
	}
	return
}

/*
 */
func (b *RefCountBranch) delete(entry []byte) (deleted bool) {
	leafLen := len(b.LeafValue)
	entryLen := len(entry)
	// does the leafValue match?
	if leafLen > 0 {
		if entryLen >= leafLen {
			for i, lb := range b.LeafValue {
				if entry[i] != lb {
					return false
				}
			}
		} else {
			return false
		}
	}

	// entry matches leaf. zero+ length
	// if there are branches there cant be End == true with a LeafValue.
	// if there are NO branches there MUST be End == true with either a LeafValue or not

	// we are at the leafend
	// log.Println("entryLen-leafLen", entryLen, leafLen, entryLen-leafLen)
	if b.End && ((entryLen - leafLen) == 0) {
		b.setEnd(false)
		if len(b.Branches) == 0 && b.Count == 0 {
			b.LeafValue = nil
		} else if len(b.Branches) == 1 && b.Count == 0 {
			b = b.pullUp()
		}
		return true
	}

	// if End == true and there are no Branches we can delete the branch because either the idx or the LeafValue mark the end - if it is matched it can be deleted
	// this is being checked in the branch above
	// prefix is matched. check for branches
	if leafLen < entryLen && b.hasBranch(entry[leafLen]) {
		// next branch matches. check the leaf/branches again
		nextBranch := b.Branches[entry[leafLen]]
		if len(nextBranch.Branches) == 0 && nextBranch.Count == 0 {
			delete(b.Branches, entry[leafLen])
			return true
		} else {
			deleted := nextBranch.delete(entry[leafLen+1:])
			if deleted && len(nextBranch.Branches) == 0 && !nextBranch.End {
				delete(b.Branches, entry[leafLen])
				// dangling leaf value?
				if len(b.Branches) == 0 && b.Count == 0 {
					b.LeafValue = nil
				}
			}
			return deleted
		}
	}

	return false
}

/*
 */
func (b *RefCountBranch) has(entry []byte) bool {
	exists, _ := b.hasCount(entry)
	return exists
}

func (b *RefCountBranch) hasCount(entry []byte) (exists bool, count int64) {
	leafLen := len(b.LeafValue)
	entryLen := len(entry)

	if entryLen >= leafLen {
		for i, pb := range b.LeafValue {
			if pb != entry[i] {
				return false, 0
			}
		}
	} else {
		return false, 0
	}

	if entryLen > leafLen {
		if br, present := b.Branches[entry[leafLen]]; present {
			return br.hasCount(entry[leafLen+1:])
		} else {
			return false, 0
		}
	} else if entryLen == leafLen && b.End {
		return true, b.Count
	}
	return false, 0
}

/*
TODO: refactor has and hasCount with this one as base
*/
func (b *RefCountBranch) getBranch(entry []byte) (be *RefCountBranch) {
	leafLen := len(b.LeafValue)
	entryLen := len(entry)

	if entryLen >= leafLen {
		for i, pb := range b.LeafValue {
			if pb != entry[i] {
				return
			}
		}
	} else {
		return
	}

	if entryLen > leafLen {
		if br, present := b.Branches[entry[leafLen]]; present {
			return br.getBranch(entry[leafLen+1:])
		} else {
			return
		}
	} else if entryLen == leafLen && b.End {
		be = b
	}
	return
}

/*
 */
func (b *RefCountBranch) hasPrefix(prefix []byte) bool {
	exists, _ := b.hasPrefixCount(prefix)
	return exists
}

func (b *RefCountBranch) hasPrefixCount(prefix []byte) (exists bool, count int64) {
	leafLen := len(b.LeafValue)
	prefixLen := len(prefix)

	if leafLen > 0 {
		if prefixLen <= leafLen {
			for i, pb := range prefix {
				if pb != b.LeafValue[i] {
					return false, b.Count
				}
			}
		} else {
			for i, lb := range b.LeafValue {
				if prefix[i] != lb {
					return false, b.Count
				}
			}
		}
	}

	if prefixLen > leafLen {
		if br, present := b.Branches[prefix[leafLen]]; present {
			return br.hasPrefixCount(prefix[leafLen+1:])
		} else {
			return false, b.Count
		}
	}

	return true, b.Count
}

/*
 */
func (b *RefCountBranch) Dump(depth int) (out string) {
	if len(b.LeafValue) > 0 {
		if b.End {
			out += fmt.Sprintf("%s V:%v (%v)\n", strings.Repeat(PADDING_CHAR, depth), string(b.LeafValue), b.Count)
		} else {
			out += fmt.Sprintf("%s V:%v (%v)\n", strings.Repeat(PADDING_CHAR, depth), string(b.LeafValue), "-")
		}
	}

	if b.End {
		out += fmt.Sprintf("%s $\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)))
	}

	for idx, branch := range b.Branches {
		if branch.End && len(branch.LeafValue) == 0 {
			out += fmt.Sprintf("%s I:%v (%v)\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)), string(idx), branch.Count)
		} else {
			out += fmt.Sprintf("%s I:%v (%v)\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)), string(idx), "-")
		}
		out += branch.Dump(depth + len(b.LeafValue) + 1)
	}

	return
}

/*
 */
func (b *RefCountBranch) hasBranches() bool {
	return len(b.Branches) == 0
}

/*
 */
func (b *RefCountBranch) hasBranch(idx byte) bool {
	if _, present := b.Branches[idx]; present {
		return true
	}
	return false
}

/*
 */
func (b *RefCountBranch) matchesLeaf(entry []byte) bool {
	leafLen := len(b.LeafValue)
	entryLen := len(entry)

	if leafLen == 0 && entryLen == 0 {
		return true
	}

	if leafLen == entryLen {
		for i, lb := range b.LeafValue {
			if entry[i] != lb {
				return false
			}
		}
	}
	return true
}

/*
 */
func (b *RefCountBranch) pullUp() *RefCountBranch {
	if len(b.Branches) == 1 {
		for k, nextBranch := range b.Branches {
			if len(nextBranch.Branches) == 0 {
				b.LeafValue = append(b.LeafValue, append([]byte{k}, nextBranch.LeafValue...)...)
			} else {
				b.LeafValue = append(b.LeafValue, k)
			}
			b.End = nextBranch.End
			b.Branches = nextBranch.Branches
			b.Count = nextBranch.Count
		}
		return b.pullUp()
	}
	return b
}

func (b *RefCountBranch) setEnd(flag bool) {
	if flag {
		b.Count += 1
	} else {
		if b.End && b.Count > 0 {
			b.Count -= 1
			if b.Count > 0 {
				return
			}
		}
	}
	b.End = flag
	return
}

func (b *RefCountBranch) String() string {
	return b.Dump(0)
}

func (b *RefCountBranch) PrintDump() {
	fmt.Printf("\n\n%s\n\n", b.Dump(0))
}

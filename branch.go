package trie

import (
	"fmt"
	"strings"
	"sync"
)

type MemberInfo struct {
	Value string
	Count int64
}

func (m *MemberInfo) String() string {
	return fmt.Sprintf("%s(%v)", m.Value, m.Count)
}

type Branch struct {
	sync.RWMutex
	Branches  map[byte]*Branch
	LeafValue []byte
	End       bool
	Count     int64
}

/*
NewBranch returns a new initialezed *Branch
*/
func (b *Branch) NewBranch() *Branch {
	return &Branch{
		Branches: make(map[byte]*Branch),
		Count:    0,
	}
}

/*
Add adds an entry to the Branch
*/
func (b *Branch) add(entry []byte) (addedBranch *Branch) {
	if b.LeafValue == nil && len(b.Branches) == 0 {
		if len(entry) > 0 {
			b.LeafValue = entry
		} else {
			// something came in but we already have branches for it
			// so the tail was the current branches index but no value
			// to push. just mark the current idx position as End
		}
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
func (b *Branch) members(branchPrefix []byte) (members []*MemberInfo) {
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
func (b *Branch) prefixMembers(branchPrefix []byte, searchPrefix []byte) (members []*MemberInfo) {
	exists, br, matchedPrefix := b.hasPrefixBranch(searchPrefix)
	if exists {
		members = br.members(matchedPrefix)
	}
	return
}

// func (b *Branch) prefixMembers(branchPrefix []byte, searchPrefix []byte) (members []*MemberInfo) {
// 	leafLen := len(b.LeafValue)
// 	searchPrefixLen := len(searchPrefix)

// 	// if the searchPrefix is empty we want all members
// 	if searchPrefixLen == 0 {
// 		members = append(members, b.members(branchPrefix)...)
// 		return
// 	}

// 	// if the searchPrefix is shorter than the leaf we will add the LeafValue
// 	// if it is an End and a the searchPrefix matches
// 	// if searchPrefixLen < leafLen {
// 	if searchPrefixLen > leafLen {
// 		for idx, br := range b.Branches {
// 			// does it match the next byte?
// 			if idx == searchPrefix[leafLen] {
// 				newSearchPrefix := searchPrefix[leafLen+1:]
// 				members = append(members, br.prefixMembers(append(append(branchPrefix, b.LeafValue...), idx), newSearchPrefix)...)
// 			}
// 		}
// 	} else if searchPrefixLen == leafLen {
// 		for i, sb := range searchPrefix {
// 			if sb != b.LeafValue[i] {
// 				return
// 			}
// 		}
// 		members = append(members, b.members(branchPrefix)...)
// 	} else {
// 		if b.End {
// 			for i, sb := range searchPrefix {
// 				if sb != b.LeafValue[i] {
// 					return
// 				}
// 			}
// 			members = append(members, b.members(branchPrefix)...)
// 			// members = append(members, &MemberInfo{string(append(branchPrefix, b.LeafValue...)), b.Count})
// 		}
// 	}
// 	return
// }

/*
 */
func (b *Branch) delete(entry []byte) (deleted bool) {
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
func (b *Branch) has(entry []byte) bool {
	if b.getBranch(entry) != nil {
		return true
	}
	return false
}

func (b *Branch) hasCount(entry []byte) (bool, int64) {
	br := b.getBranch(entry)
	if br != nil {
		return true, br.Count
	}
	return false, 0
}

func (b *Branch) getBranch(entry []byte) (be *Branch) {
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
func (b *Branch) hasPrefix(prefix []byte) bool {
	exists, _, _ := b.hasPrefixBranch(prefix)
	return exists
}

func (b *Branch) hasPrefixCount(prefix []byte) (exists bool, count int64) {
	exists, br, _ := b.hasPrefixBranch(prefix)
	if exists {
		count = br.sumCount()
	}
	return
}

func (b *Branch) hasPrefixBranch(prefix []byte) (exists bool, branch *Branch, matchedPrefix []byte) {
	leafLen := len(b.LeafValue)
	prefixLen := len(prefix)
	exists = false
	var pref []byte

	if leafLen > 0 {
		if prefixLen <= leafLen {
			for i, pb := range prefix {
				if pb != b.LeafValue[i] {
					return
				}
			}
		} else {
			for i, lb := range b.LeafValue {
				if prefix[i] != lb {
					return
				}
			}
			matchedPrefix = append(matchedPrefix, prefix[:leafLen]...)
		}
	}

	if prefixLen > leafLen {
		if br, present := b.Branches[prefix[leafLen]]; present {
			matchedPrefix = append(matchedPrefix, prefix[leafLen])
			exists, branch, pref = br.hasPrefixBranch(prefix[leafLen+1:])
			matchedPrefix = append(matchedPrefix, pref...)
			return
		} else {
			return
		}
	}
	return true, b, matchedPrefix
}

func (b *Branch) sumCount() (count int64) {
	// leaf itself matches
	if b.End {
		count += b.Count
	}
	for _, br := range b.Branches {
		count += br.sumCount()
	}
	return
}

/*
 */
func (b *Branch) Dump(depth int) (out string) {
	if len(b.LeafValue) > 0 {
		if b.End {
			out += fmt.Sprintf("%s V:%v %v (%v)\n", strings.Repeat(PADDING_CHAR, depth), string(b.LeafValue), b.LeafValue, b.Count)
		} else {
			out += fmt.Sprintf("%s V:%v %v (%v)\n", strings.Repeat(PADDING_CHAR, depth), string(b.LeafValue), b.LeafValue, "-")
		}
	}

	if b.End {
		out += fmt.Sprintf("%s $\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)))
	}

	for idx, branch := range b.Branches {
		if branch.End && len(branch.LeafValue) == 0 {
			out += fmt.Sprintf("%s I:%v %v (%v)\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)), string(idx), idx, branch.Count)
		} else {
			out += fmt.Sprintf("%s I:%v %v (%v)\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)), string(idx), idx, "-")
		}
		out += branch.Dump(depth + len(b.LeafValue) + 1)
	}

	return
}

/*
 */
// func (b *Branch) hasBranches() bool {
// 	return len(b.Branches) == 0
// }

/*
 */
func (b *Branch) hasBranch(idx byte) bool {
	if _, present := b.Branches[idx]; present {
		return true
	}
	return false
}

/*
 */
// func (b *Branch) matchesLeaf(entry []byte) bool {
// 	leafLen := len(b.LeafValue)
// 	entryLen := len(entry)

// 	if leafLen == 0 && entryLen == 0 {
// 		return true
// 	}

// 	if leafLen == entryLen {
// 		for i, lb := range b.LeafValue {
// 			if entry[i] != lb {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

/*
 */
func (b *Branch) pullUp() *Branch {
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

func (b *Branch) setEnd(flag bool) {
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

func (b *Branch) String() string {
	return b.Dump(0)
}

func (b *Branch) PrintDump() {
	fmt.Printf("\n%s\n\n", b)
}

package trie

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	// "sort"
	"strings"
	"sync"
	"time"
)

func init() {
	// gob.Register(&trie.EnDecodeTrie{})
}

type Trie struct {
	Root *Branch
}

func NewTrie() *Trie {
	t := &Trie{}
	t.Root = NewBranch()
	return t
}

func (t *Trie) Add(entry string) {
	t.Root.Add([]byte(entry))
}

func (t *Trie) Delete(entry string) bool {
	return t.Root.Delete([]byte(entry))
}

func (t *Trie) Dump() {
	t.Root.Dump(0)
}

func (t *Trie) PrintDump() {
	t.Root.PrintDump()
}

func (t *Trie) Has(prefix string) bool {
	return t.Root.Has([]byte(prefix))
}

func (t *Trie) Members() []string {
	return t.Root.Members([]byte{})
}

func (t *Trie) PrefixMembers(prefix string) []string {
	return t.Root.PrefixMembers([]byte{}, []byte(prefix))
}

func (t *Trie) DumpToFile(fname string) (err error) {
	entries := t.Members()
	// sort.Sort(sort.Reverse(sort.StringSlice(entries)))
	// fmt.Println(entries)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(entries); err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	bl, err := w.Write(buf.Bytes())
	fmt.Printf("wrote %d bytes\n", bl)
	w.Flush()
	return
}

func LoadFromFile(fname string) (tr *Trie, err error) {
	f, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)

	var entries []string
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&entries); err != nil {
		fmt.Println("decoding error:", err)
	}

	// fmt.Println(entries)

	startTime := time.Now()
	tr = NewTrie()
	for _, word := range entries {
		tr.Add(word)
	}
	fmt.Printf("adding words to index took: %v\n", time.Since(startTime))

	return
}

type Branch struct {
	sync.RWMutex
	Branches  map[byte]*Branch
	LeafValue []byte // tail end
	End       bool
}

func NewBranch() *Branch {
	return &Branch{
		Branches: make(map[byte]*Branch),
	}
}

func (b *Branch) Add(entry []byte) {
	if b.LeafValue == nil && len(b.Branches) == 0 {
		b.Lock()
		b.LeafValue = entry
		b.End = true
		b.Unlock()
		return
	}

	// something came in but we already have branches for it
	// so the tail was the current branches index but no value
	// to push. just mark the current idx position as End
	if len(b.LeafValue) == 0 && len(entry) == 0 {
		b.Lock()
		b.End = true
		b.Unlock()
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
		// fmt.Printf("\n ++ ++ newLeafLen < len(b.LeafValue) |%s|  |%s|\n", string(newLeaf), string(entry))
		// fmt.Println("vv", string(b.LeafValue), b.End)
		tail := b.LeafValue[newLeafLen:]
		idx := tail[0]
		newBranch := NewBranch()
		newBranch.LeafValue = tail[1:]
		b.Lock()
		b.LeafValue = newLeaf
		newBranch.Branches, b.Branches = b.Branches, newBranch.Branches
		newBranch.End, b.End = b.End, newBranch.End
		b.Branches[idx] = newBranch
		b.Unlock()
	}

	// new leaf is smaller than the entry, which means there will be more stuff
	// that we need to push down
	if newLeafLen < len(entry) {
		tail := entry[newLeafLen:]
		idx := tail[0]
		// fmt.Printf("\nnewLeafLen < len(entry) |%s| |%s|\n", string(newLeaf), string(entry))
		// fmt.Println(">>>", string(b.LeafValue), b.End, "at idx", idx)

		// create new branch at idx if it does not exists yet
		b.Lock()
		if _, notPresent := b.Branches[idx]; !notPresent {
			b.Branches[idx] = NewBranch()
			// fmt.Printf("NewBranch at idx: %v for newleaf %s, entry %s \n", string(idx), string(newLeaf), string(entry))
		}
		// check whether the idx itself marks an End $. if so add a new idx
		// fmt.Println(">+>> send down", string(tail[1:]), "at idx", string(idx), "which currently has", len(b.Branches[idx].Branches), "branches and LeafVal:", b.Branches[idx].LeafValue)
		b.Branches[idx].Add(tail[1:])
		b.Unlock()
	} else {
		// if there is nothing else to be pushed down we just have to mark the
		// current branch as a end. this happens when you add a value that already
		// us covered by the index but this particular end had not been marked.
		// eg. you already have 'foo' in your index and now add 'f'.
		b.Lock()
		b.End = true
		b.Unlock()
	}
}

func (b *Branch) Members(branchPrefix []byte) (members []string) {
	if b.End {
		members = append(members, string(append(branchPrefix, b.LeafValue...)))
	}
	for idx, br := range b.Branches {
		newPrefix := append(append(branchPrefix, b.LeafValue...), idx)
		members = append(members, br.Members(newPrefix)...)
	}
	return
}

func (b *Branch) PrefixMembers(branchPrefix []byte, searchPrefix []byte) (members []string) {
	leafLen := len(b.LeafValue)
	searchPrefixLen := len(searchPrefix)

	// if the searchPrefix is empty we want all members
	if searchPrefixLen == 0 {
		members = append(members, b.Members(branchPrefix)...)
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
				members = append(members, br.PrefixMembers(append(append(branchPrefix, b.LeafValue...), idx), newSearchPrefix)...)
			}
		}
	} else if searchPrefixLen == leafLen {
		for i, sb := range searchPrefix {
			if sb != b.LeafValue[i] {
				return
			}
		}
		members = append(members, b.Members(branchPrefix)...)
	} else {
		if b.End {
			for i, sb := range searchPrefix {
				if sb != b.LeafValue[i] {
					return
				}
			}
			members = append(members, string(append(branchPrefix, b.LeafValue...)))
		}
	}
	return
}

func (b *Branch) Delete(entry []byte) (deleted bool) {
	leafLen := len(b.LeafValue)
	entryLen := len(entry)

	// we are at a leaf end.
	if b.End && entryLen == leafLen {
		for i, pb := range entry {
			if pb != b.LeafValue[i] {
				return false
			}
		}
		if len(b.Branches) > 0 {
			b.End = false
		}
		return true
	}

	// does the leaf match
	if leafLen > 0 {
		if entryLen <= leafLen {
			for i, pb := range entry {
				if pb != b.LeafValue[i] {
					return false
				}
			}
		} else {
			for i, lb := range b.LeafValue {
				if entry[i] != lb {
					return false
				}
			}
		}
	}

	if entryLen > leafLen {
		if br, present := b.Branches[entry[leafLen]]; present {
			return br.Has(entry[leafLen+1:])
		} else {
			return false
		}
	}

	return
}

func (b *Branch) Has(prefix []byte) bool {
	leafLen := len(b.LeafValue)
	prefixLen := len(prefix)

	if leafLen > 0 {
		if prefixLen <= leafLen {
			for i, pb := range prefix {
				if pb != b.LeafValue[i] {
					return false
				}
			}
		} else {
			for i, lb := range b.LeafValue {
				if prefix[i] != lb {
					return false
				}
			}
		}
	}

	if prefixLen > leafLen {
		// if len(b.Branches) == 0 {
		// 	return false
		// }
		if br, present := b.Branches[prefix[leafLen]]; present {
			return br.Has(prefix[leafLen+1:])
		} else {
			return false
		}
	}

	return true
}

const PADDING_CHAR = "-"

func (b *Branch) Dump(depth int) (out string) {
	if len(b.LeafValue) > 0 {
		out += fmt.Sprintf("%s V:%v\n", strings.Repeat(PADDING_CHAR, depth), string(b.LeafValue))
		// out += fmt.Sprintf("%s V:%v\n", strings.Repeat(PADDING_CHAR, depth), b.LeafValue)
	}

	if b.End {
		out += fmt.Sprintf("%s $\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)))
	}

	for idx, branch := range b.Branches {
		out += fmt.Sprintf("%s I:%v\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)), string(idx))
		// out += fmt.Sprintf("%s I:%v\n", strings.Repeat(PADDING_CHAR, depth+len(b.LeafValue)), idx)
		out += branch.Dump(depth + len(b.LeafValue) + 1)
	}

	return
}

func (b *Branch) PrintDump() {
	fmt.Printf("\n\n%s\n\n", b.Dump(0))
}

// persistence

// type EnDecodeTrie struct {
// 	Root *EnDecodeBranch
// }

// type EnDecodeBranch struct {
// 	Branches  map[byte]*Branch
// 	LeafValue []byte // tail end
// 	End       bool
// }

// func (t *Trie) FileDump(fname string) {

// 	edTrie := &EnDecodeTrie{
// 	Root: &EnDecodeBranch{
// 	Branches
// 	}
// 	}

// 	buf := new(bytes.Buffer)
// 	enc := gob.NewEncoder(buf)
// 	if err := enc.Encode(tr); err != nil {
// 		log.Println(err)
// 	}

// 	f, err := os.Create(fname)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer f.Close()

// 	w := bufio.NewWriter(f)
// 	blength, err := w.Write(buf.Bytes())
// 	log.Printf("wrote %d bytes\n", blength)
// 	w.Flush()
// }

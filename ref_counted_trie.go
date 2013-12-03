package trie

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// TRIE

type RefCountTrie struct {
	Root                     *RefCountBranch
	OpsCount                 int
	DumpOpsCount             int
	PersistThresholdOpsCount int
	// PersistThresholdTime time.Duration
}

/*
NewTrie returns the pointer to a new Trie with an initiallized root Branch
*/
func NewRefCountTrie() *RefCountTrie {
	t := &RefCountTrie{
		OpsCount:                 0,
		DumpOpsCount:             0,
		PersistThresholdOpsCount: 0,
	}
	t.Root = &RefCountBranch{
		Branches: make(map[byte]*RefCountBranch),
	}
	return t
}

/*
Add adds an entry to the trie and returns the branch node that the insertion
was made at - or rather where the end of the entry was marked.
*/
func (t *RefCountTrie) Add(entry string) *RefCountBranch {
	t.Root.Lock()
	b := t.Root.add([]byte(entry))
	t.OpsCount += 1
	t.Root.Unlock()
	return b
}

/*
Delete decrements the count of an existing entry by one. If the count equals
zero it removes an the entry from the trie. Returns true if the entry existed,
false otherwise. Note that the return value says something about the previous
existence of the entry - not whether it has been completely removed or just
its count decremented.
*/
func (t *RefCountTrie) Delete(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	t.Root.Lock()
	deleted := t.Root.delete([]byte(entry))
	t.OpsCount += 1
	t.Root.Unlock()
	return deleted
}

/*
Has returns true if the `entry` exists in the `Trie`
*/
func (t *RefCountTrie) Has(entry string) bool {
	return t.Root.has([]byte(entry))
}

/*
HasCount returns true  if the `entry` exists in the `Trie`. The second returned
value is the count how often the entry has been set.
*/
func (t *RefCountTrie) HasCount(entry string) (exists bool, count int64) {
	return t.Root.hasCount([]byte(entry))
}

/*
HasPrefix returns true if the the `Trie` contains entries with the given prefix
*/
func (t *RefCountTrie) HasPrefix(prefix string) bool {
	return t.Root.hasPrefix([]byte(prefix))
}

/*
HasPrefixCount returns true if the the `Trie` contains entries with the given
prefix. The second returned value is the count how often the entry has been set.
*/
func (t *RefCountTrie) HasPrefixCount(prefix string) (exists bool, count int64) {
	return t.Root.hasPrefixCount([]byte(prefix))
}

/*
Members returns all entries of the Trie with their counts as MemberInfo
*/
func (t *RefCountTrie) Members() []*MemberInfo {
	return t.Root.members([]byte{})
}

/*
Members returns a Slice of all entries of the Trie
*/
func (t *RefCountTrie) MembersList() (members []string) {
	for _, mi := range t.Root.members([]byte{}) {
		members = append(members, mi.Value)
	}
	return
}

/*
PrefixMembers returns all entries of the Trie that have the given prefix
with their counts as MemberInfo
*/
func (t *RefCountTrie) PrefixMembers(prefix string) []*MemberInfo {
	return t.Root.prefixMembers([]byte{}, []byte(prefix))
}

/*
PrefixMembers returns a List of all entries of the Trie that have the
given prefix
*/
func (t *RefCountTrie) PrefixMembersList(prefix string) (members []string) {
	for _, mi := range t.Root.prefixMembers([]byte{}, []byte(prefix)) {
		members = append(members, mi.Value)
	}
	return
}

func (t *RefCountTrie) DumpToFileWithMinOps(fname string) (err error) {
	if t.OpsCount >= t.PersistThresholdOpsCount {
		err = t.DumpToFile(fname)
	} else {
		// log.Println(t.OpsCount)
	}
	return
}

/*
DumpToFile dumps all values into a slice of strings and writes that to a file
using encoding/gob.

The Trie itself can currently not be encoded directly because gob does not
directly support structs with a sync.Mutex on them.
*/
func (t *RefCountTrie) DumpToFile(fname string) (err error) {
	t.Root.Lock()
	t.DumpOpsCount = t.OpsCount
	entries := t.Members()
	t.Root.Unlock()

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(entries); err != nil {
		err = errors.New(fmt.Sprintf("Could encode Trie entries for dump file: %v", err))
		return
	}

	f, err := os.Create(fname)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not save dump file: %v", err))
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	bl, err := w.Write(buf.Bytes())
	if err != nil {
		err = errors.New(fmt.Sprintf("Error writing to dump file: %v", err))
		return
	}
	log.Printf("wrote %d bytes to dumpfile %s\n", bl, fname)
	w.Flush()
	t.Root.Lock()
	t.OpsCount -= t.DumpOpsCount
	t.DumpOpsCount = 0
	t.Root.Unlock()
	return
}

/*
LoadFromFile loads a gib encoded wordlist from a file and creates a new Trie
by Add()ing all of them.
*/
func RCTLoadFromFile(fname string) (tr *RefCountTrie, err error) {
	log.Println("Load trie from", fname)
	f, err := os.Open(fname)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open Trie file: %v", err))
		tr = NewRefCountTrie()
	} else {
		defer f.Close()

		buf := bufio.NewReader(f)
		var entries []*MemberInfo
		dec := gob.NewDecoder(buf)
		if err = dec.Decode(&entries); err != nil {
			if err == io.EOF && entries == nil {
				log.Println("Nothing to decode. Seems the file is empty.")
				err = nil
			} else {
				err = errors.New(fmt.Sprintf("Decoding error: %v", err))
				return
			}
		}

		tr = NewRefCountTrie()
		startTime := time.Now()
		for _, mi := range entries {
			b := tr.Add(mi.Value)
			b.Count = mi.Count
		}
		tr.DumpOpsCount = 0
		tr.OpsCount = 0
		log.Printf("adding words to index took: %v\n", time.Since(startTime))
	}

	return
}

/*
Dump returns a string representation of the `Trie`
*/
func (t *RefCountTrie) Dump() string {
	return t.Root.Dump(0)
}

/*
 */
func (t *RefCountTrie) PrintDump() {
	t.Root.PrintDump()
}

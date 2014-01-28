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

type Trie struct {
	Root *Branch
}

/*
NewTrie returns the pointer to a new Trie with an initialized root Branch
*/
func NewTrie() *Trie {
	t := &Trie{
		Root: &Branch{
			Branches: make(map[byte]*Branch),
		},
	}
	return t
}

/*
Add adds an entry to the trie and returns the branch node that the insertion
was made at - or rather where the end of the entry was marked.
*/
func (t *Trie) Add(entry string) *Branch {
	t.Root.Lock()
	b := t.Root.add([]byte(entry))
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
func (t *Trie) Delete(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	t.Root.Lock()
	deleted := t.Root.delete([]byte(entry))
	t.Root.Unlock()
	return deleted
}

/*
GetBranch returns the branch end if the `entry` exists in the `Trie`
*/
func (t *Trie) GetBranch(entry string) *Branch {
	return t.Root.getBranch([]byte(entry))
}

/*
Has returns true if the `entry` exists in the `Trie`
*/
func (t *Trie) Has(entry string) bool {
	return t.Root.has([]byte(entry))
}

/*
HasCount returns true  if the `entry` exists in the `Trie`. The second returned
value is the count how often the entry has been set.
*/
func (t *Trie) HasCount(entry string) (exists bool, count int64) {
	return t.Root.hasCount([]byte(entry))
}

/*
HasPrefix returns true if the the `Trie` contains entries with the given prefix
*/
func (t *Trie) HasPrefix(prefix string) bool {
	return t.Root.hasPrefix([]byte(prefix))
}

/*
HasPrefixCount returns true if the the `Trie` contains entries with the given
prefix. The second returned value is the count how often the entry has been set.
*/
func (t *Trie) HasPrefixCount(prefix string) (exists bool, count int64) {
	return t.Root.hasPrefixCount([]byte(prefix))
}

/*
Members returns all entries of the Trie with their counts as MemberInfo
*/
func (t *Trie) Members() []*MemberInfo {
	return t.Root.members([]byte{})
}

/*
Members returns a Slice of all entries of the Trie
*/
func (t *Trie) MembersList() (members []string) {
	for _, mi := range t.Root.members([]byte{}) {
		members = append(members, mi.Value)
	}
	return
}

/*
PrefixMembers returns all entries of the Trie that have the given prefix
with their counts as MemberInfo
*/
func (t *Trie) PrefixMembers(prefix string) []*MemberInfo {
	return t.Root.prefixMembers([]byte{}, []byte(prefix))
}

/*
PrefixMembers returns a List of all entries of the Trie that have the
given prefix
*/
func (t *Trie) PrefixMembersList(prefix string) (members []string) {
	for _, mi := range t.Root.prefixMembers([]byte{}, []byte(prefix)) {
		members = append(members, mi.Value)
	}
	return
}

/*
Dump returns a string representation of the `Trie`
*/
func (t *Trie) Dump() string {
	return t.Root.Dump(0)
}

/*
 */
func (t *Trie) PrintDump() {
	t.Root.PrintDump()
}

/*
DumpToFile dumps all values into a slice of strings and writes that to a file
using encoding/gob.

The Trie itself can currently not be encoded directly because gob does not
directly support structs with a sync.Mutex on them.
*/
func (t *Trie) DumpToFile(fname string) (err error) {
	t.Root.Lock()
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
	_, err = w.Write(buf.Bytes())
	if err != nil {
		err = errors.New(fmt.Sprintf("Error writing to dump file: %v", err))
		return
	}
	// log.Printf("wrote %d bytes to dumpfile %s\n", bl, fname)
	w.Flush()
	return
}

/*
MergeFromFile loads a gib encoded wordlist from a file and Add() them to the `Trie`.
*/
// TODO: write tests for merge
func (t *Trie) MergeFromFile(fname string) (err error) {
	entries, err := loadTrieFile(fname)
	if err != nil {
		return
	}
	log.Printf("Got %v entries\n", len(entries))
	startTime := time.Now()
	for _, mi := range entries {
		b := t.GetBranch(mi.Value)
		if b != nil {
			b.Lock()
			b.Count += mi.Count
			b.Unlock()
		} else {
			b := t.Add(mi.Value)
			b.Lock()
			b.Count = mi.Count
			b.Unlock()
		}
	}
	log.Printf("merging words to index took: %v\n", time.Since(startTime))
	return
}

/*
LoadFromFile loads a gib encoded wordlist from a file and creates a new Trie
by Add()ing all of them.
*/
func LoadFromFile(fname string) (tr *Trie, err error) {
	tr = NewTrie()
	entries, err := loadTrieFile(fname)
	if err != nil {
		return
	}
	log.Printf("Got %v entries\n", len(entries))
	startTime := time.Now()
	for _, mi := range entries {
		b := tr.Add(mi.Value)
		b.Count = mi.Count
	}
	log.Printf("adding words to index took: %v\n", time.Since(startTime))

	return
}

func loadTrieFile(fname string) (entries []*MemberInfo, err error) {
	log.Println("Load trie from", fname)
	f, err := os.Open(fname)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open Trie file: %v", err))
	} else {
		defer f.Close()

		buf := bufio.NewReader(f)
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
	}

	return
}

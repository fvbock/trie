package trie

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func init() {
	// gob.Register(&trie.EnDecodeTrie{})
}

type Trie struct {
	Root *Branch
	// Cardinaliy int
}

/*
NewTrie returns the pointer to a new Trie with an initiallized root Branch
*/
func NewTrie() *Trie {
	t := &Trie{}
	t.Root = NewBranch()
	return t
}

/*
Add adds an non existing entry to the trie
*/
func (t *Trie) Add(entry string) {
	t.Root.Add([]byte(entry))
}

/*
Delete removes an existing entry from the trie
*/
func (t *Trie) Delete(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	return t.Root.delete([]byte(entry))
}

/*
Has returns true if the `entry` exists in the `Trie`
*/
func (t *Trie) Has(entry string) bool {
	return t.Root.Has([]byte(entry))
}

/*
Has returns true if the the `Trie` contains entries with the given prefix
*/
func (t *Trie) HasPrefix(prefix string) bool {
	return t.Root.HasPrefix([]byte(prefix))
}

/*
Members returns all entries of the Trie
*/
func (t *Trie) Members() []string {
	return t.Root.Members([]byte{})
}

/*
PrefixMembers returns all entries of the Trie that have the given prefix
*/
func (t *Trie) PrefixMembers(prefix string) []string {
	return t.Root.prefixMembers([]byte{}, []byte(prefix))
}

/*
DumpToFile dumps all values into a slice of strings and writes that to a file
using encoding/gob.

The Trie itself can currently not be encoded directly because gob does not
directly support structs with a sync.Mutex on them.
*/
func (t *Trie) DumpToFile(fname string) (err error) {
	entries := t.Members()
	// sort.Sort(sort.Reverse(sort.StringSlice(entries)))

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
	fmt.Printf("wrote %d bytes\n", bl)
	w.Flush()
	return
}

/*
LoadFromFile loads a gib encoded wordlist from a file and creates a new Trie
by Add()ing all of them.
*/
func LoadFromFile(fname string) (tr *Trie, err error) {
	fmt.Println("Load trie from", fname)
	f, err := os.Open(fname)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open Trie file: %v", err))
		tr = NewTrie()
	} else {
		defer f.Close()

		buf := bufio.NewReader(f)
		var entries []string
		dec := gob.NewDecoder(buf)
		if err = dec.Decode(&entries); err != nil {
			if err == io.EOF && entries == nil {
				fmt.Println("Nothing to decode. Seems the file is empty.")
				err = nil
			} else {
				err = errors.New(fmt.Sprintf("Decoding error: %v", err))
				return
			}
		}

		tr = NewTrie()
		startTime := time.Now()
		for _, word := range entries {
			tr.Add(word)
		}
		fmt.Printf("adding words to index took: %v\n", time.Since(startTime))
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

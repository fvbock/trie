package trie

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
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
 */
func (t *Trie) Add(entry string) {
	t.Root.Add([]byte(entry))
}

/*
 */
func (t *Trie) Delete(entry string) bool {
	if len(entry) == 0 {
		return false
	}
	return t.Root.delete([]byte(entry))
}

/*
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
 */
func (t *Trie) Has(entry string) bool {
	return t.Root.Has([]byte(entry))
}

/*
 */
func (t *Trie) HasPrefix(prefix string) bool {
	return t.Root.HasPrefix([]byte(prefix))
}

/*
 */
func (t *Trie) Members() []string {
	return t.Root.Members([]byte{})
}

/*
 */
func (t *Trie) PrefixMembers(prefix string) []string {
	return t.Root.prefixMembers([]byte{}, []byte(prefix))
}

/*
 */
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

/*
 */
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

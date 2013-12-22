package trie

import (
	// "fmt"
	"github.com/fvbock/uds-go/set"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(1)
}

func TestTrieAddSingle(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
}

func TestTrieAddBigSmall(t *testing.T) {
	tr := NewTrie()
	tr.Add("testing")
	tr.Add("tests")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestTrieAddSmallBig(t *testing.T) {
	tr := NewTrie()
	tr.Add("tests")
	tr.Add("testing")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestTrieGetBranch(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	tr.Add("testing")
	t.Logf("\n%s", tr.Dump())

	b1 := tr.GetBranch("test")
	if b1 == nil {
		t.Error("Expected to find a branch for 'test'.")
	}

	b2 := tr.GetBranch("tests")
	if b2 != nil {
		t.Error("Expected not to find a branch for 'tests'.")
	}

	b3 := tr.GetBranch("testing")
	if b3 == nil {
		t.Error("Expected to find a branch for 'testing'.")
	}

	b4 := tr.GetBranch("testi")
	if b4 != nil {
		t.Error("Expected not to find a branch for 'testi'.")
	}

	b5 := tr.GetBranch("tessi")
	if b5 != nil {
		t.Error("Expected not to find a branch for 'tessi'.")
	}
}

// func TestTrieAddEmptyBranch(t *testing.T) {
// 	tr := NewTrie()
// 	tr.Add("foobar")
// 	tr.Add("fooc")
// 	if tr.Root.End {
// 		t.Error("Expected Root End to be false")
// 	}
// 	t.Logf("\n%s", tr.Dump())
// 	// tr.Add("foob")
// 	// tr.Add("fooba")
// 	// tr.Add("fooca")
// 	// t.Logf("\n%s", tr.Dump())
// 	// if !tr.Root.End {
// 	// 	t.Error("Expected Root End to be true")
// 	// }
// }

func TestTrieAddBigSmallMulti(t *testing.T) {
	tr := NewTrie()
	tr.Add("testing")
	tr.Add("testing")
	tr.Add("tests")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
	_, c1 := tr.HasCount("testing")
	if c1 != 2 {
		t.Errorf("Expected count for testing to be 2. got %v instead", c1)
	}
	_, c2 := tr.HasCount("tests")
	if c2 != 1 {
		t.Errorf("Expected count for tests to be 1. got %v instead.", c2)
	}
}

func TestTrieAddSmallBigMulti(t *testing.T) {
	tr := NewTrie()
	tr.Add("tests")
	tr.Add("tests")
	tr.Add("testing")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
	_, c1 := tr.HasCount("testing")
	if c1 != 1 {
		t.Errorf("Expected count for testing to be 1. got %v instead", c1)
	}
	_, c2 := tr.HasCount("tests")
	if c2 != 2 {
		t.Errorf("Expected count for tests to be 2. got %v instead.", c2)
	}
}

func TestTrieAddTestFirst(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	tr.Add("testing")
	tr.Add("tests")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
	if !tr.Root.End {
		t.Error("Expected trunk End to be true")
	}
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestTrieAddTestLast(t *testing.T) {
	tr := NewTrie()
	tr.Add("testing")
	tr.Add("tests")
	tr.Add("test")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestTrieDump(t *testing.T) {
	tr := NewTrie()
	tr.Add("teased")
	tr.Add("test")
	tr.Add("test")
	tr.Add("testing")
	tr.Add("tea")
	t.Logf("\n%s", tr.Dump())
}

func TestTrieMembersCount(t *testing.T) {
	tr := NewTrie()
	tr.Add("teased")
	tr.Add("test")
	tr.Add("test")
	tr.Add("testing")

	if len(tr.Members()) != 3 {
		t.Error("Expected 3 members")
	}
	for _, mi := range tr.Members() {
		if mi.Value == "teased" && mi.Count != 1 {
			t.Error("Expected teased to have Count 1")
			continue
		}
		if mi.Value == "test" && mi.Count != 2 {
			t.Error("Expected test to have Count 2")
			continue
		}
		if mi.Value == "testing" && mi.Count != 1 {
			t.Error("Expected testing to have Count 1")
			continue
		}
	}
	t.Logf("\n%v", tr.Members())
}

// // todo
// func TestTriePrefixMembersCount(t *testing.T) {
// 	tr := NewTrie()
// 	tr.Add("foo")
// 	tr.Add("foobar")
// 	tr.Add("bar")

// 	if tr.MembersCount("test") != 0 {
// 		t.Error("Expected HasCount for test to be 0")
// 	}
// }

func TestTriePrefixMembers(t *testing.T) {
	tr := NewTrie()
	tr.Add("teased")
	tr.Add("test")
	tr.Add("test")
	tr.Add("testing")
	if len(tr.PrefixMembers("test")) != 2 {
		t.Error("Expected PrefixMembers('test') to have length 2")
	}
	if len(tr.PrefixMembers("te")) != 3 {
		t.Error("Expected PrefixMembers('te') to have length 3")
	}
	if len(tr.PrefixMembers("a")) != 0 {
		t.Error("Expected PrefixMembers('a') to have length 0")
	}
	tr.PrintDump()
	if len(tr.PrefixMembers("ta")) != 0 {
		t.Error("Expected PrefixMembers('ta') to have length 0")
	}
	if len(tr.PrefixMembers("")) != 3 {
		t.Error("Expected PrefixMembers('') to have length 3")
	}
	if len(tr.PrefixMembersList("a")) != 0 {
		t.Error("Expected PrefixMembersList('a') to have length 0")
	}
	if len(tr.PrefixMembersList("")) != 3 {
		t.Error("Expected PrefixMembersList('') to have length 3")
	}

	// cover different code paths
	tr.Add("te")
	tr.PrintDump()
	if len(tr.PrefixMembers("a")) != 0 {
		t.Error("Expected PrefixMembers('a') to have length 0")
	}
	tl := len(tr.PrefixMembers("t"))
	if tl != 4 {
		t.Errorf("Expected PrefixMembers('t') to have length 4, got %v instead.", tl)
	}
}

func TestTrieHasPrefixEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.HasPrefix("test") {
		t.Error("Expected no prefix test")
	}
	_, c := tr.HasPrefixCount("test")
	if c != 0 {
		t.Error("Expected no prefix test")
	}
}

func TestTrieHasPrefixOne(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.HasPrefix("test") {
		t.Error("Expected prefix test")
	}
	_, c := tr.HasPrefixCount("test")
	if c != 1 {
		t.Error("Expected prefix test to have count 1")
	}
}

func TestTrieHasPrefixMany(t *testing.T) {
	tr := NewTrie()
	tr.Add("tease")
	tr.Add("teases")
	tr.Add("teased")
	tr.Add("teaser")
	tr.Add("tests")
	tr.Add("test")
	tr.Add("tested")
	tr.Add("testing")
	if tr.HasPrefix("ted") {
		t.Error("Expected no prefix ted")
	}
	if !tr.HasPrefix("tease") {
		t.Error("Expected prefix tease")
	}
	if !tr.HasPrefix("testing") {
		t.Error("Expected prefix testing")
	}

	// prefixCount
	_, ctest := tr.HasPrefixCount("test")
	if ctest != 4 {
		t.Errorf("Expected prefix test to have count 4, got %v instead.", ctest)
	}
	_, ctes := tr.HasPrefixCount("tes")
	if ctes != 4 {
		t.Errorf("Expected prefix tes to have count 4, got %v instead.", ctes)
	}

	_, ctea := tr.HasPrefixCount("tea")
	if ctea != 4 {
		t.Errorf("Expected prefix tea to have count 4, got %v instead.", ctea)
	}
	tr.Add("tea")
	_, ctea = tr.HasPrefixCount("tea")
	if ctea != 5 {
		t.Errorf("Expected prefix tea to have count 5, got %v instead.", ctea)
	}

	// test false cases with shorter and longer than leaf prefixes
	_, ca := tr.HasPrefixCount("a")
	if ca != 0 {
		t.Errorf("Expected prefix a to have count 0, got %v instead.", ca)
	}
	_, casdf := tr.HasPrefixCount("asdf")
	if casdf != 0 {
		t.Errorf("Expected prefix asdf to have count 0, got %v instead.", casdf)
	}
}

func TestTrieHasEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.Has("test") {
		t.Error("Expected no test")
	}
}

func TestTrieHasOne(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.Has("test") {
		t.Error("Expected test")
	}
}

func TestTrieHasMany(t *testing.T) {
	tr := NewTrie()
	tr.Add("tease")
	tr.Add("teases")
	tr.Add("teased")
	tr.Add("teaser")
	tr.Add("tests")
	tr.Add("test")
	tr.Add("tested")
	tr.Add("testing")
	if tr.Has("testi") {
		t.Error("Expected no testi")
	}
	if tr.Has("te") {
		t.Error("Expected no te")
	}
	if !tr.Has("tease") {
		t.Error("Expected tease")
	}
	if !tr.Has("testing") {
		t.Error("Expected testing")
	}
}

func TestTrieHasPrefixManyMultibyte(t *testing.T) {
	tr := NewTrie()
	tr.Add("日本人")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")
	if tr.HasPrefix("ä") {
		t.Error("Expected no prefix ä")
	}
	if tr.HasPrefix("無い") {
		t.Error("Expected no prefix 無い")
	}
	if !tr.HasPrefix("日本語") {
		t.Error("Expected prefix 日本語")
	}
	if !tr.HasPrefix("日") {
		t.Error("Expected prefix 日")
	}
}

func TestTrieHasManyMultibyte(t *testing.T) {
	tr := NewTrie()
	tr.Add("日本人")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")
	if tr.Has("ä") {
		t.Error("Expected no ä")
	}
	if tr.Has("無い") {
		t.Error("Expected no 無い")
	}
	if tr.Has("日") {
		t.Error("Expected no 日")
	}
	if !tr.Has("日本語") {
		t.Error("Expected 日本語")
	}
	if !tr.Has("学校") {
		t.Error("Expected 学校")
	}
}

func TestTrieDeleteEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.Delete("test") {
		t.Error("Expected false for tr.Delete('test')")
	}
}

func TestTrieDeleteNothing(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if tr.Delete("") {
		t.Error("Expected false for tr.Delete('')")
	}

	_, c1 := tr.HasCount("test")
	if c1 != 1 {
		t.Errorf("Expected count for test to be 1. got %v instead", c1)
	}

	if tr.Delete("tes") {
		t.Error("Expected false for tr.Delete('tes')")
	}
}

func TestTrieDeleteOne(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}

	// // delete a branch that has no further branches
	// tr.Add("teste")
	// if !tr.Delete("teste") {
	// 	t.Error("Expected true for tr.Delete('test')")
	// }
}

func TestTrieDeleteDouble(t *testing.T) {
	tr := NewTrie()
	tr.Add("foo")
	tr.Add("test")
	tr.Add("test")
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}
	tr.PrintDump()
	t.Log(tr.Members())
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}
	tr.PrintDump()
	t.Log(tr.Members())
}

func TestTrieDeletePrefixCount(t *testing.T) {
	tr := NewTrie()
	tr.Add("foo")
	tr.Add("foo")
	tr.Add("foobar")
	tr.PrintDump()
	if tr.Delete("test") {
		t.Error("Expected false for tr.Delete('test')")
	}
	if !tr.Delete("foo") {
		t.Error("Expected true for tr.Delete('foo')")
	}
	tr.PrintDump()
	_, cfoo := tr.HasCount("foo")
	if cfoo != 1 {
		t.Errorf("Expected count for foo to be 1. got %v instead.", cfoo)
	}
	_, cfoobar := tr.HasCount("foobar")
	if cfoobar != 1 {
		t.Errorf("Expected count for foobar to be 1. got %v instead.", cfoobar)
	}
	if !tr.Delete("foo") {
		t.Error("Expected true for tr.Delete('foo')")
	}
	tr.PrintDump()
	_, cfoo = tr.HasCount("foo")
	if cfoo != 0 {
		t.Errorf("Expected count for foo to be 0. got %v instead.", cfoo)
	}
	_, cfoobar = tr.HasCount("foobar")
	if cfoobar != 1 {
		t.Errorf("Expected count for foobar to be 1. got %v instead.", cfoobar)
	}
}

func TestTrieDeleteMany(t *testing.T) {
	tr := NewTrie()
	tr.Add("tease")
	tr.Add("teases")
	tr.Add("teased")
	tr.Add("test")
	tr.Add("test")

	// if tr.Delete("te") {
	// 	t.Error("Expected false for tr.Delete('te')")
	// }
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}

	expectedMembers := make(map[string]bool)
	expectedMembers["tease"] = true
	expectedMembers["teases"] = true
	expectedMembers["teased"] = true
	expectedMembers["test"] = true
	// expectedMembers["test"] = true
	for _, m := range tr.Members() {
		if m.Count != 1 {
			t.Errorf("Expected Count for %s to be 1 - not %v.", m.Value, m.Count)
		} else {
			ec := len(expectedMembers)
			delete(expectedMembers, m.Value)
			if len(expectedMembers) == ec {
				t.Errorf("Not expected member %s.", m.Value)
			}
		}
	}

	if len(expectedMembers) != 0 {
		t.Log(tr.Members())
		t.Error("Deletion seems to have deleted more than just 'test' (once).", expectedMembers)
	}

	if !tr.Delete("tease") {
		t.Error("Expected true for tr.Delete('tease')")
	}
	if !tr.Delete("teases") {
		t.Error("Expected true for tr.Delete('tease')")
	}
	if !tr.Delete("teased") {
		t.Error("Expected true for tr.Delete('tease')")
	}

	tr.PrintDump()
	t.Log(tr.Members())
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}

	tr.PrintDump()
	t.Log(tr.Members())

	if len(tr.Root.Branches) != 0 {
		t.Error("Expected 0 Branches on Root")
	}
	if len(tr.Root.LeafValue) != 0 {
		t.Error("Expected no LeafValue on Root")
	}
	if tr.Root.End {
		t.Error("Expected End to be false on Root")
	}
}

func TestTrieDeleteManyRandom_az(t *testing.T) {
	tr := NewTrie()
	var prefix = "prefix"
	var words []string
	var str []byte
	var n = 0
	for n < 100 {
		i := 0
		str = []byte{}
		for i < 10 {
			rn := 0
			for rn < 97 {
				rn = rand.Intn(123)
			}
			str = append(str, byte(rn))
			i++
		}
		if rand.Intn(2) == 1 {
			words = append(words, prefix+string(str))
			tr.Add(prefix + string(str))
		} else {
			words = append(words, string(str))
			tr.Add(string(str))
		}
		n++
	}
	// t.Log(words)
	// tr.PrintDump()
	for wi, w := range words {
		if !tr.Delete(w) {
			t.Errorf("Expected true for tr.Delete('%s')", w)
		}
		// expect to still find the rest
		if wi+1 < len(words) {
			for _, ow := range words[wi+1:] {
				// t.Logf("Checking for %s", ow)
				if !tr.Has(ow) {
					t.Errorf("Expected to still find %s", ow)
				}
			}
		}
	}
	tr.PrintDump()
	if len(tr.Root.Branches) != 0 {
		t.Error("Expected 0 Branches on Root")
	}
	if len(tr.Root.LeafValue) != 0 {
		t.Error("Expected no LeafValue on Root")
	}
	if tr.Root.End {
		t.Error("Expected End to be false on Root")
	}
}

func TestTrieMultiAdd(t *testing.T) {
	tr := NewTrie()
	words := []string{"foodie", "foods", "foodchain", "foodcrave", "food", "人", "日本", "日本語学校", "学校", "日本語"}
	// words := []string{"日本語", "日本語学校"}
	// words := []string{"日本語学校", "日本"}
	wg := sync.WaitGroup{}
	for _, w := range words {
		// wg.Add(1)
		// go func(word string) {
		// 	tr.Add(word)
		// 	wg.Done()
		// }(w)

		// tr.Add(w)
		// tr.Add(w)
		// if w == "日本" {
		// 	tr.PrintDump()
		// 	tr.Delete(w)
		// 	tr.PrintDump()
		// }

		// wg.Add(2)
		// go func(word string) {
		// 	tr.Add(word)
		// 	wg.Done()
		// }(w)
		// go func(word string) {
		// 	tr.Add(word)
		// 	wg.Done()
		// }(w)
		// go func(word string) {
		// 	wg.Add(1)
		// 	if word == "日本" {
		// 		tr.PrintDump()
		// 		tr.Delete(word)
		// 		tr.PrintDump()
		// 	}
		// 	wg.Done()
		// }(w)

		// wg.Add(3)
		// go func(word string) {
		// 	tr.Add(word)
		// 	wg.Done()
		// }(w)
		// go func(word string) {
		// 	tr.Delete(word)
		// 	wg.Done()
		// }(w)
		// go func(word string) {
		// 	tr.Add(word)
		// 	wg.Done()
		// }(w)

		wg.Add(5)
		go func(word string) {
			tr.Add(word)
			wg.Done()
		}(w)
		go func(word string) {
			tr.Delete(word)
			wg.Done()
		}(w)
		go func(word string) {
			tr.Add(word)
			wg.Done()
		}(w)
		go func(word string) {
			tr.Delete(word)
			wg.Done()
		}(w)
		go func(word string) {
			tr.Add(word)
			wg.Done()
		}(w)

	}
	wg.Wait()
	tr.PrintDump()
	t.Log(tr.Members())
}

func TestTrieDumpToFileLoadFromFile(t *testing.T) {
	tr := NewTrie()
	var prefix = "prefix"
	var words []string
	var str []byte
	var insert string
	var n = 0
	for n < 100 {
		i := 0
		str = []byte{}
		for i < 10 {
			rn := 0
			for rn < 97 {
				rn = rand.Intn(123)
			}
			str = append(str, byte(rn))
			i++
		}
		if rand.Intn(2) == 1 {
			insert = prefix + string(str)
		} else {
			insert = string(str)
		}
		words = append(words, insert)
		tr.Add(insert)
		if rand.Intn(2) == 1 {
			tr.Add(insert)
		}
		n++
	}
	err := tr.DumpToFile("testfiles/TestDumpToFileLoadFromFile")

	loadedTrie, err := LoadFromFile("testfiles/TestDumpToFileLoadFromFile")
	if err != nil {
		t.Errorf("Failed to load Trie from file: %v", err)
	}
	for _, w := range words {
		// t.Logf("Checking for %s", w)
		if !loadedTrie.Has(w) {
			t.Errorf("Expected to find %s", w)
		}
	}

	trMembers := set.NewStringSet(tr.MembersList()...)
	loadedTrieMembers := set.NewStringSet(loadedTrie.MembersList()...)

	t.Log("trMembers.IsEqual(loadedTrieMembers):", trMembers.IsEqual(loadedTrieMembers))

	diff := trMembers.Difference(loadedTrieMembers)
	if diff.Len() > 0 {
		t.Error("Dump() of the original and the LoadFromFile() version of the Trie are different.")
	}

	// check counts
	for _, mi := range tr.Members() {
		_, count := loadedTrie.HasCount(mi.Value)
		if count != mi.Count {
			t.Errorf("Count for member %s differs: orig was %v, restored trie has %v", mi.Value, mi.Count, count)
		}
	}

	// test expected failures
	if tr.DumpToFile("dirdoesnotexist/TestDumpToFileLoadFromFile") == nil {
		t.Error("expected DumpToFile() to fail with non existent directory.")
	}
}

func TestTrieLoadFromFileEmpty(t *testing.T) {
	loadedTrie, err := LoadFromFile("testfiles/empty")
	if err != nil {
		t.Errorf("Failed to load Trie from file: %v", err)
	}

	loadedTrieMembers := set.NewStringSet(loadedTrie.MembersList()...)
	t.Log(loadedTrieMembers)
	t.Log(loadedTrieMembers.Len())
	if loadedTrieMembers.Len() > 0 {
		t.Error("Expected 0 Members from LoadFromFile() with an empty file.")
	}
}

func TestTrieLoadFromFileExpectedFailures(t *testing.T) {
	_, err := LoadFromFile("testfiles/notatriedump")
	if err == nil {
		t.Error("Expected LoadFromFile to fail - file testfiles/notatriedump is not a valid Trie dump.")
	}
	_, err = LoadFromFile("doesnotexist/doesnotexist")
	if err == nil {
		t.Error("Expected LoadFromFile to fail - file testfiles/doesnotexist does notexist.")
	}
}

func TestTrieDumpToFileMergeFromFile(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	tr.Add("test")
	tr.Add("tested")
	tr.Add("tent")
	tr.DumpToFile("testfiles/TestDumpToFileMergeFromFile")

	tr2 := NewTrie()
	tr2.Add("tea")
	tr2.Add("tested")

	err := tr2.MergeFromFile("testfiles/TestDumpToFileMergeFromFile")
	if err != nil {
		t.Errorf("Failed to merge Trie from file: %v", err)
	}

	_, ctest := tr2.HasCount("test")
	if ctest != 2 {
		t.Errorf("Expected count for test to be 2. got %v instead.", ctest)
	}
	_, ctested := tr2.HasCount("tested")
	if ctested != 2 {
		t.Errorf("Expected count for tested to be 2. got %v instead.", ctested)
	}
	_, ctea := tr2.HasCount("tea")
	if ctea != 1 {
		t.Errorf("Expected count for tea to be 1. got %v instead.", ctea)
	}
	_, ctent := tr2.HasCount("tent")
	if ctent != 1 {
		t.Errorf("Expected count for tent to be 1. got %v instentd.", ctent)
	}

	// expected failure
	err = tr2.MergeFromFile("doesnotexist/doesnotexist")
	if err == nil {
		t.Error("Expected MergeFromFile to fail - file testfiles/doesnotexist does notexist.")
	}
}

// some simple benchmarks

func BenchmarkTrieBenchAdd(b *testing.B) {
	randstr := make([][]byte, 100)
	i := 0
	for i < 100 {
		randstr[i] = []byte{}
		n := 0
		for n < 100 {
			randstr[i] = append(randstr[i], byte(rand.Intn(255)))
			n++
		}
		i++
	}

	tr := NewTrie()
	for x := 0; x < b.N; x++ {
		tr.Add(string(randstr[x%100]))
	}
}

func BenchmarkTrieBenchHasPrefix(b *testing.B) {
	tr := NewTrie()
	randstr := make([][]byte, 10000)
	i := 0
	for i < 10000 {
		randstr[i] = []byte{}
		n := 0
		for n < 100 {
			randstr[i] = append(randstr[i], byte(rand.Intn(255)))
			n++
		}
		i++
	}

	for x := 0; x < 100000; x++ {
		tr.Add(string(randstr[x%10000]))
	}

	for x := 0; x < b.N; x++ {
		tr.HasPrefix(string(randstr[x%10000]))
	}
}

func BenchmarkTrieBenchHas(b *testing.B) {
	tr := NewTrie()
	randstr := make([][]byte, 10000)
	i := 0
	for i < 10000 {
		randstr[i] = []byte{}
		n := 0
		for n < 100 {
			randstr[i] = append(randstr[i], byte(rand.Intn(255)))
			n++
		}
		i++
	}

	for y := 0; y < 100000; y++ {
		tr.Add(string(randstr[y%10000]))
	}

	for x := 0; x < b.N; x++ {
		tr.Has(string(randstr[x%10000]))
	}
}

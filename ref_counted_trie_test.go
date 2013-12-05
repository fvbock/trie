package trie

import (
	"github.com/fvbock/uds-go/set"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(4)
}

func TestRefCountTrieAddSingle(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("test")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
}

func TestRefCountTrieAddBigSmall(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("testing")
	tr.Add("tests")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestRefCountTrieAddSmallBig(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("tests")
	tr.Add("testing")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestRefCountTrieGetBranch(t *testing.T) {
	tr := NewRefCountTrie()
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
}

func TestRefCountTrieAddBigSmallMulti(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieAddSmallBigMulti(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieAddTestFirst(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieAddTestLast(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieDump(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("teased")
	tr.Add("test")
	tr.Add("test")
	tr.Add("testing")
	tr.Add("tea")
	t.Logf("\n%s", tr.Dump())
}

func TestRefCountTrieMembersCount(t *testing.T) {
	tr := NewRefCountTrie()
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
		// t.Errorf("Unexpected member: %v", mi)
	}
	t.Logf("\n%v", tr.Members())
}

// // todo
// func TestRefCountTriePrefixMembersCount(t *testing.T) {
// 	tr := NewRefCountTrie()
// 	tr.Add("foo")
// 	tr.Add("foobar")
// 	tr.Add("bar")

// 	if tr.MembersCount("test") != 0 {
// 		t.Error("Expected HasCount for test to be 0")
// 	}
// }

func TestRefCountTriePrefixMembersCountFromFile(t *testing.T) {
	tr := NewRefCountTrie()
	tr, err := RCTLoadFromFile("testfiles/trie_idx_5018d345558fbe46c4000001")
	// tr, err := RCTLoadFromFile("/tmp/trie_idx_5018d345558fbe46c4000001")
	if err != nil {
		t.Errorf("Failed to load Trie from file: %v", err)
	}
	t.Logf("\n%v", len(tr.Members()))
	t.Logf("\n%v", tr.PrefixMembers("test"))
	t.Logf("\n%v", tr.PrefixMembersList("test"))
	// tr.PrintDump()
}

func TestRefCountTrieHasPrefixEmpty(t *testing.T) {
	tr := NewRefCountTrie()
	if tr.HasPrefix("test") {
		t.Error("Expected no prefix test")
	}
}

func TestRefCountTrieHasPrefixOne(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("test")
	if !tr.HasPrefix("test") {
		t.Error("Expected prefix test")
	}
}

func TestRefCountTrieHasPrefixMany(t *testing.T) {
	tr := NewRefCountTrie()
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
}

func TestRefCountTrieHasEmpty(t *testing.T) {
	tr := NewRefCountTrie()
	if tr.Has("test") {
		t.Error("Expected no test")
	}
}

func TestRefCountTrieHasOne(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("test")
	if !tr.Has("test") {
		t.Error("Expected test")
	}
}

func TestRefCountTrieHasMany(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieHasPrefixManyMultibyte(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieHasManyMultibyte(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieDeleteEmpty(t *testing.T) {
	tr := NewRefCountTrie()
	if tr.Delete("test") {
		t.Error("Expected false for tr.Delete('test')")
	}
}

func TestRefCountTrieDeleteOne(t *testing.T) {
	tr := NewRefCountTrie()
	tr.Add("test")
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}
}

func TestRefCountTrieDeleteDouble(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieDeletePrefixCount(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieDeleteMany(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieDeleteManyRandom_az(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieMultiAdd(t *testing.T) {
	tr := NewRefCountTrie()
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

func TestRefCountTrieDumpToFileRCTLoadFromFile(t *testing.T) {
	tr := NewRefCountTrie()
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
	tr.DumpToFile("testfiles/TestDumpToFileRCTLoadFromFile")

	loadedTrie, err := RCTLoadFromFile("testfiles/TestDumpToFileRCTLoadFromFile")
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
}

func TestRefCountTrieLoadFromFileEmpty(t *testing.T) {
	loadedTrie, err := RCTLoadFromFile("testfiles/empty")
	if err != nil {
		t.Errorf("Failed to load Trie from file: %v", err)
	}

	loadedTrieMembers := set.NewStringSet(loadedTrie.MembersList()...)
	t.Log(loadedTrieMembers)
	t.Log(loadedTrieMembers.Len())
	if loadedTrieMembers.Len() > 0 {
		t.Error("Expected 0 Members from RCTLoadFromFile() with an empty file.")
	}
}

// some simple benchmarks

func BenchmarkRefCountTrieBenchAdd(b *testing.B) {
	for x := 0; x < b.N; x++ {
		tr := NewRefCountTrie()
		tr.Add("foodie")
		tr.Add("foods")
		tr.Add("foodchain")
		tr.Add("foodcrave")
		tr.Add("food")
		tr.Add("人")
		tr.Add("日本")
		tr.Add("日本語学校")
		tr.Add("学校")
		tr.Add("日本語")
	}
}

func BenchmarkRefCountTrieBenchHasPrefix(b *testing.B) {
	tr := NewRefCountTrie()
	tr.Add("foodie")
	tr.Add("foods")
	tr.Add("foodchain")
	tr.Add("foodcrave")
	tr.Add("food")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")

	for x := 0; x < b.N; x++ {
		tr.HasPrefix("foodie")
		tr.HasPrefix("foodcrave")
		tr.HasPrefix("日本")
		tr.HasPrefix("学校")
		tr.HasPrefix("thisisnotinhere")
		tr.HasPrefix("学日本校")
	}
}

func BenchmarkRefCountTrieBenchHas(b *testing.B) {
	tr := NewRefCountTrie()
	tr.Add("foodie")
	tr.Add("foods")
	tr.Add("foodchain")
	tr.Add("foodcrave")
	tr.Add("food")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")

	for x := 0; x < b.N; x++ {
		tr.Has("foodie")
		tr.Has("foodcrave")
		tr.Has("日本")
		tr.Has("学校")
		tr.Has("thisisnotinhere")
		tr.Has("学日本校")
	}
}

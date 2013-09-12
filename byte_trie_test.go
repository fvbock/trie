package trie

import (
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	tr := NewTrie()
	t.Logf("Empty Has('foo')? %v", tr.Has("foo"))

	tr.Add("foo")
	tr.Add("foodie")
	tr.PrintDump()
	tr.Add("foods")
	tr.PrintDump()
	tr.Add("foodchain")
	tr.Add("foodcrave")
	tr.Add("food")
	// tr.Add("f")

	tr.PrintDump()

	t.Logf("Members(): %v", tr.Members())

	t.Logf("PrefixMembers('foo'): %v", tr.PrefixMembers("foo"))
	t.Logf("PrefixMembers('food'): %v", tr.PrefixMembers("food"))

	t.Logf("Has('foo')? %v", tr.Has("foo"))
	t.Logf("Has('food')? %v", tr.Has("food"))
	t.Logf("Has('foodie')? %v", tr.Has("foodie"))
	t.Logf("Has('foods')? %v", tr.Has("foods"))
	t.Logf("Has('foodstore')? %v", tr.Has("foodstore"))

	tr2 := NewTrie()
	tr2.Add("日本人")
	tr2.Add("人")
	tr2.Add("日本")
	tr2.Add("日本語学校")
	tr2.Add("学校")
	tr2.Add("日本語")
	tr2.Add("ä")

	t.Logf("Members(): %v", tr2.Members())
	t.Logf("PrefixMembers('日本語'): %v", tr2.PrefixMembers("日本語"))
	t.Logf("PrefixMembers('日本語'): %v", tr2.PrefixMembers("日本"))

	tr2.PrintDump()

	t.Logf("Has('日本')? %v", tr2.Has("日本"))
	t.Logf("Has('日')? %v", tr2.Has("日"))
	t.Logf("Has('日本語')? %v", tr2.Has("日本語"))
	t.Logf("Has('{')? %v", tr2.Has("{"))
	t.Logf("Has('æ')? %v", tr2.Has("æ"))
	t.Logf("Has('ä')? %v", tr2.Has("ä"))

	// t.Logf("order does not matter: %v\n", tr3.Root.Dump(0) == tr4.Root.Dump(0))
}

func TestMultiAdd(t *testing.T) {
	tr := NewTrie()
	words := []string{"foodie", "foods", "foodchain", "foodcrave", "food", "人", "日本", "日本語学校", "学校", "日本語"}
	wg := sync.WaitGroup{}
	for _, w := range words {
		wg.Add(1)
		go func(word string) {
			tr.Add(word)
			wg.Done()
		}(w)
	}
	wg.Wait()
	tr.PrintDump()
}

func BenchmarkBenchAdd(b *testing.B) {
	for x := 0; x < b.N; x++ {
		tr := NewTrie()
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

func BenchmarkBenchHas(b *testing.B) {
	tr := NewTrie()
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

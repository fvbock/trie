trie
====

A Trie (Prefix Index) implementation in golang. It works fine with all unicode characters.

Documentation can be found [at godoc.org](http://godoc.org/github.com/fvbock/trie).

Entries are reference counted: If you `Add("foo")` twice and `Del("foo")` it once it will still be found.

[![Build Status](https://travis-ci.org/fvbock/trie.png)](https://travis-ci.org/fvbock/trie)

Example
=======

	t := trie.NewTrie()
	t.Add("foo")
	t.Add("bar")
	t.PrintDump()

	// output:
	//  I:f (-)
	// - V:oo (1)
	// --- $
	//  I:b (-)
	// - V:ar (1)
	// --- $

	t.Add("foo")
	t.PrintDump()

	// output:
	//  I:f (-)
	// - V:oo (2)
	// --- $
	//  I:b (-)
	// - V:ar (1)
	// --- $

	fmt.Println(t.Has("foo"))
	// output: true

	fmt.Println(t.HasCount("foo"))
	// output: true 2

	fmt.Println(t.Has("foobar"))
	// output: false

	fmt.Println(t.Members())
	// output: [foo(2) bar(1)]

	t.Add("food")
	t.Add("foobar")
	t.Add("foot")
	fmt.Println(t.HasPrefix("foo"))
	// output: true

	fmt.Println(t.PrefixMembers("foo"))
	// output: [foo(2) food(1) foobar(1) foot(1)]


A `Trie` can be dumped into a file with

	t.DumpToFile("/tmp/trie_foo")

And loaded with

	t2, _ := trie.LoadFromFile("/tmp/trie_foo")
	fmt.Println(t2.Members())
	// output: [foo(2) food(1) foobar(1) foot(1) bar(1)]

An existing `Trie` can be merged with a stored one with

	t3 := trie.NewTrie()
	t3.Add("フー")
	t3.Add("バー")
	t3.Add("日本語")
	fmt.Println(t3.Members())
	// output: [フー(1) バー(1) 日本語(1)]

	t3.MergeFromFile("/tmp/trie_foo")
	fmt.Println(t3.Members())
	// output: [フー(1) バー(1) 日本語(1) foo(2) food(1) foobar(1) foot(1) bar(1)]

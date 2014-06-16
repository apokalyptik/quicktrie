package trie

import "testing"

func TestTrie(t *testing.T) {
	Debug = true
	trie := NewTrie()
	if trie == nil {
		t.Error("Expected a new Trie")
		return
	}
	if trie.GetString("asdf") {
		t.Errorf("Expected get on empty trie to fail")
	}
	trie.AddString("asdf")
	if trie.GetString("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}
	if trie.GetString("AAA") {
		t.Errorf("Expected trie to be case insensitive")
	}
	if trie.GetString("aaaa") {
		t.Errorf("Expected get on non-added key to still be false")
	}
	if trie.GetString("aa") {
		t.Errorf("Expected get on non-added prefix of a key to be false")
	}
	trie.AddString("aaab")
	if trie.GetString("aaab") == false {
		t.Errorf("Expected get on added key to be true")
	}
	if trie.GetString("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}
	trie.DelString("aaab")
	if trie.GetString("aaab") {
		t.Errorf("Expected get on deleted key to be false")
	}
	if trie.GetString("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}

	trie.AddString("y")
	trie.AddString("yy")
	trie.AddString("yyy")
	trie.AddString("yyyy")
	trie.AddString("yyyyy")
	trie.DelString("y")
	trie.DelString("yy")
	trie.DelString("yyy")
	trie.DelString("yyyy")
	trie.DelString("yyyyy")

	trie.AddString("w")
	trie.AddString("ww")
	trie.AddString("www")
	trie.AddString("wwww")
	trie.AddString("wwwww")
	trie.DelString("wwwww")
	trie.DelString("wwww")
	trie.DelString("www")
	trie.DelString("ww")
	trie.DelString("w")

	trie.AddString("z")
	trie.AddString("zz")
	trie.AddString("zzz")
	trie.AddString("zzzz")
	trie.AddString("zzzzz")
	trie.DelString("zz")
	trie.DelString("z")
	trie.DelString("zzzz")
	trie.DelString("zzz")
	trie.DelString("zzzzz")

	trie.AddString("bbb")
	trie.AddString("xx")
	trie.AddString("qqq")
	trie.AddString("qqqq")
	trie.AddString("qqqqq")
	trie.AddString("xxx")
	trie.AddString("xxx")
	trie.AddString("xxx")
	trie.DelString("qqq")

	var found = 0
	trie.IterateString(func(key string) { found++ })
	if found != 6 {
		t.Errorf("Expected 6 values when iterating, got %d", found)
	}

	trie.AddString("nnn")
	trie.AddString("nnnnn")
	trie.AddString("nnnnnn")
	trie.AddString("nnnnnnn")
	trie.AddString("nnnnnnnn")
	trie.AddString("nnnnnnnnn")

	found = 0
	trie.IterateString(func(key string) { found++ })
	if found != 12 {
		t.Errorf("Expected 12 values when iterating, got %d", found)
	}

	trie.DropString("nnnnnnn")

	found = 0
	trie.IterateString(func(key string) { found++ })
	if found != 9 {
		t.Errorf("Expected 9 values when iterating, got %d", found)
	}

	trie.DropString("nnnn")

	found = 0
	trie.IterateString(func(key string) { found++ })
	if found != 7 {
		t.Errorf("Expected 7 values when iterating, got %d", found)
	}

	trie.DropString("n")

	found = 0
	trie.IterateString(func(key string) { found++ })
	if found != 6 {
		t.Errorf("Expected 6 values when iterating, got %d", found)
	}
	if Debug {
		trie.Log()
	}
}
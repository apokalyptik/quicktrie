package trie

import "testing"

func TestBWTrie(t *testing.T) {
	trie := NewBWTrie()
	if trie == nil {
		t.Error("Expected a new Trie")
		return
	}
	if trie.Exists("asdf") {
		t.Errorf("Expected get on empty trie to fail")
	}
	trie.Add("asdf")
	if trie.Exists("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}
	if trie.Exists("AAA") {
		t.Errorf("Expected trie to be case insensitive")
	}
	if trie.Exists("aaaa") {
		t.Errorf("Expected get on non-added key to still be false")
	}
	if trie.Exists("aa") {
		t.Errorf("Expected get on non-added prefix of a key to be false")
	}
	trie.Add("aaab")
	if trie.Exists("aaab") == false {
		t.Errorf("Expected get on added key to be true")
	}
	if trie.Exists("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}
	trie.Del("aaab")
	if trie.Exists("aaab") {
		t.Errorf("Expected get on deleted key to be false")
	}
	if trie.Exists("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}

	trie.Add("y")
	trie.Add("yy")
	trie.Add("yyy")
	trie.Add("yyyy")
	trie.Add("yyyyy")
	trie.Del("y")
	trie.Del("yy")
	trie.Del("yyy")
	trie.Del("yyyy")
	trie.Del("yyyyy")

	trie.Add("w")
	trie.Add("ww")
	trie.Add("www")
	trie.Add("wwww")
	trie.Add("wwwww")
	trie.Del("wwwww")
	trie.Del("wwww")
	trie.Del("www")
	trie.Del("ww")
	trie.Del("w")

	trie.Add("z")
	trie.Add("zz")
	trie.Add("zzz")
	trie.Add("zzzz")
	trie.Add("zzzzz")
	trie.Del("zz")
	trie.Del("z")
	trie.Del("zzzz")
	trie.Del("zzz")
	trie.Del("zzzzz")

	trie.Add("bbb")
	trie.Add("xx")
	trie.Add("qqq")
	trie.Add("qqqq")
	trie.Add("qqqqq")
	trie.Add("xxx")
	trie.Add("xxx")
	trie.Add("xxx")
	trie.Del("qqq")

	var found = trie.Count()

	trie.Add("nnn")
	trie.Add("nnnnn")
	trie.Add("nnnnnn")
	trie.Add("nnnnnnn")
	trie.Add("nnnnnnnn")
	trie.Add("nnnnnnnnn")

	found = trie.Count()
	if found != 12 {
		t.Errorf("Expected 12 values when iterating, got %d", found)
	}

	trie.Drop("nnnnnnn")

	found = trie.Count()
	if found != 9 {
		t.Errorf("Expected 9 values when iterating, got %d", found)
	}

	trie.Drop("nnnn")

	found = trie.Count()
	if found != 7 {
		t.Errorf("Expected 7 values when iterating, got %d", found)
	}

	trie.Drop("n")

	found = trie.Count()
	if found != 6 {
		t.Errorf("Expected 6 values when iterating, got %d", found)
	}
}

func TestKVIterateFrom(t *testing.T) {
	var trie = NewKVTrie()
	trie.Add("b")
	trie.Add("a")
	trie.Add("aa")
	trie.Add("aaa")
	trie.Add("aab")
	trie.Add("aabb")
	trie.Add("aaaa")
	if c := len(trie.GetBranch("")); c != 7 {
		t.Errorf("Expected prefix iteration to return 6 results from '', found %d", c)
	}
	if c := len(trie.GetBranch("a")); c != 6 {
		t.Errorf("Expected prefix iteration to return 6 results from 'a', found %d", c)
	}
	if c := len(trie.GetBranch("aa")); c != 5 {
		t.Errorf("Expected prefix iteration to return 5 results from 'aa', found %d", c)
	}
	if c := len(trie.GetBranch("aaa")); c != 2 {
		t.Errorf("Expected prefix iteration to return 2 results from 'aaa', found %d", c)
	}
}

func TestBWIterateFrom(t *testing.T) {
	var trie = NewBWTrie()
	trie.Add("b")
	trie.Add("a")
	trie.Add("aa")
	trie.Add("aaa")
	trie.Add("aab")
	trie.Add("aabb")
	trie.Add("aaaa")
	if c := len(trie.GetBranch("")); c != 7 {
		t.Errorf("Expected prefix iteration to return 6 results from '', found %d", c)
	}
	if c := len(trie.GetBranch("a")); c != 6 {
		t.Errorf("Expected prefix iteration to return 6 results from 'a', found %d", c)
	}
	if c := len(trie.GetBranch("aa")); c != 5 {
		t.Errorf("Expected prefix iteration to return 5 results from 'aa', found %d", c)
	}
	if c := len(trie.GetBranch("aaa")); c != 2 {
		t.Errorf("Expected prefix iteration to return 2 results from 'aaa', found %d", c)
	}
}

func TestKVTrie(t *testing.T) {
	trie := NewKVTrie()
	trie.Add("to", "data: to")
	trie.Add("tea", "data: tea")
	trie.Add("ten", "data: ten")

	if e, v := trie.Get("to"); !e || v == nil {
		t.Errorf("Expected key 'to' to exist, and return non nil data")
	} else if v.(string) != "data: to" {
		t.Errorf("Expected 'to' to conaint 'data: to', got: '%s'", v.(string))
	}

	if e, v := trie.Get("tea"); !e || v == nil {
		t.Errorf("Expected key 'tea' to exist, and return non nil data")
	} else if v.(string) != "data: tea" {
		t.Errorf("Expected 'tea' to contain 'data: tea', got: '%s'", v.(string))
	}

	trie.Set("to", "newdata")
	if e, v := trie.Get("to"); !e || v == nil {
		t.Errorf("Expected key 'to' to exist, and return non nil data")
	} else if v.(string) != "newdata" {
		t.Errorf("Expected 'to' to conaint 'newdata', got: '%s'", v.(string))
	}

	trie.Del("to")

	if e, v := trie.Get("to"); e || v != nil {
		t.Errorf("Expected key 'tea' to be nil, and return nil data")
	}

	if e, v := trie.Get("tea"); !e || v == nil {
		t.Errorf("Expected key 'tea' to exist, and return non nil data")
	} else if v.(string) != "data: tea" {
		t.Errorf("Expected 'tea' to contain 'data: tea', got: '%s'", v.(string))
	}

	trie.Drop("te")

	if e, v := trie.Get("ten"); e || v != nil {
		t.Errorf("Expected key 'ten' to be nil, and return nil data")
	}

	if e, v := trie.Get("tea"); e || v != nil {
		t.Errorf("Expected key 'tea' to be nil, and return nil data")
	}
}

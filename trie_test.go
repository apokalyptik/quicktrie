package trie

import (
	"log"
	"strconv"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := NewTrie()
	if trie == nil {
		t.Error("Expected a new Trie")
		return
	}
	if trie.Get("asdf") {
		t.Errorf("Expected get on empty trie to fail")
	}
	trie.Add("asdf")
	if trie.Get("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}
	if trie.Get("AAA") {
		t.Errorf("Expected trie to be case insensitive")
	}
	if trie.Get("aaaa") {
		t.Errorf("Expected get on non-added key to still be false")
	}
	if trie.Get("aa") {
		t.Errorf("Expected get on non-added prefix of a key to be false")
	}
	trie.Add("aaab")
	if trie.Get("aaab") == false {
		t.Errorf("Expected get on added key to be true")
	}
	if trie.Get("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}
	trie.Del("aaab")
	if trie.Get("aaab") {
		t.Errorf("Expected get on deleted key to be false")
	}
	if trie.Get("asdf") == false {
		t.Errorf("Expected get on added key to be true")
	}

	trie.Add("bbb")
	trie.Add("xx")
	trie.Add("qqq")
	trie.Add("qqqq")
	trie.Add("qqqqq")
	trie.Add("xxx")
	trie.Add("xxx")
	trie.Add("xxx")
	trie.Del("qqq")

	var found = 0
	trie.Iterate(func(key string) {
		if key == "" {
			return
		}
		found++
	})
	if found != 6 {
		log.Printf("Expected 4 values when iterating")
	}
}

func BenchmarkDuplicateSmallEntries(b *testing.B) {
	trie := NewTrie()
	for n := 0; n < b.N; n++ {
		trie.Add("aaa")
	}
}

func BenchmarkDuplicateLargeEntries(b *testing.B) {
	trie := NewTrie()
	for n := 0; n < b.N; n++ {
		trie.Add("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}
}

func BenchmarkNonDupSmallEntries(b *testing.B) {
	trie := NewTrie()
	for n := 0; n < b.N; n++ {
		trie.Add("aaa" + strconv.Itoa(n))
	}
}

func BenchmarkNonDupLargeEntries(b *testing.B) {
	trie := NewTrie()
	for n := 0; n < b.N; n++ {
		trie.Add("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + strconv.Itoa(n))
	}
}

func BenchmarkDeleteEmptyTrie(b *testing.B) {
	trie := NewTrie()
	for n := 0; n < b.N; n++ {
		trie.Del("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	}
}

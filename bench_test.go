package trie

import (
	"fmt"
	"testing"
)

func BenchmarkBwInsert(b *testing.B) {
	t := NewBWTrie()
	for n := 0; n < b.N; n++ {
		s := fmt.Sprintf("key %d", n)
		t.Add(s)
	}
}

func BenchmarkKvInsert(b *testing.B) {
	t := NewKVTrie()
	for n := 0; n < b.N; n++ {
		s := fmt.Sprintf("key %d", n)
		t.Add(s, n)
	}
}

package trie

import (
	"fmt"
	"log"
)

type IterFunc func([]byte)

type node interface {
	get([]byte) bool
	add([]byte)
	del([]byte)
	drop([]byte)
	iterate([]byte, IterFunc)
	log(...int)
}

var Debug bool

type Trie struct {
	root node
}

func NewTrie() *Trie {
	return &Trie{
		root: &BWTrie{
			children: []*BWTrie{},
		},
	}
}

func (t *Trie) AddString(key string) {
	t.Add([]byte(key))
}

func (t *Trie) Add(key []byte) {
	if Debug {
		log.Printf("add: %s\n", string(key))
	}
	t.root.add([]byte(key))
}

func (t *Trie) DropString(key string) {
	t.Drop([]byte(key))
}

func (t *Trie) Drop(key []byte) {
	t.root.drop(key)
}

func (t *Trie) DelString(key string) {
	t.Del([]byte(key))
}

func (t *Trie) Del(key []byte) {
	if Debug {
		log.Printf("del: %s\n", string(key))
	}
	t.root.del(key)
}

func (t *Trie) GetString(key string) bool {
	return t.Get([]byte(key))
}

func (t *Trie) Get(key []byte) bool {
	return t.root.get(key)
}

func (t *Trie) IterateString(callback func(string)) {
	t.Iterate(func(b []byte) { callback(string(b)) })
}

func (t *Trie) Iterate(callback func([]byte)) {
	t.root.iterate([]byte{}, callback)
}

func (t *Trie) Print() {
	t.root.iterate([]byte{}, func(k []byte) { fmt.Println(string(k)) })
}

func (t *Trie) Log() {
	t.root.log(0)
}

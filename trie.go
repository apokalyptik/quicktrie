package trie

import (
	"fmt"
	"log"
)

type IterFunc func([]byte, interface{})
type IterStringFunc func(string, interface{})

type node interface {
	get([]byte) (bool, interface{})
	add([]byte, ...interface{})
	del([]byte)
	drop([]byte)
	iterate([]byte, IterFunc)
	log(...int)
}

var Debug bool

type Trie struct {
	kind uint8
	root node
}

func NewTrie() *Trie {
	return &Trie{
		root: &BWTrie{
			children: []*BWTrie{},
		},
	}
}

func NewDataTrie() *Trie {
	return &Trie{
		kind: 1,
		root: &BWTrie{
			children: []*BWTrie{},
		},
	}
}

func (t *Trie) AddString(key string, data ...interface{}) {
	t.Add([]byte(key))
}

func (t *Trie) Add(key []byte, data ...interface{}) {
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

func (t *Trie) Exists(key interface{}) bool {
	switch key := key.(type) {
	case []byte:
		k, _ := t.root.get(key)
		return k
	case string:
		k, _ := t.root.get([]byte(key))
		return k
	default:
		return false
	}
}

func (t *Trie) GetString(key string) (bool, interface{}) {
	return t.Get([]byte(key))
}

func (t *Trie) Get(key []byte) (bool, interface{}) {
	return t.root.get(key)
}

func (t *Trie) IterateString(callback IterStringFunc) {
	t.Iterate(func(b []byte, i interface{}) { callback(string(b), i) })
}

func (t *Trie) Iterate(callback IterFunc) {
	t.root.iterate([]byte{}, callback)
}

func (t *Trie) Print() {
	t.root.iterate([]byte{}, func(k []byte, _ interface{}) { fmt.Println(string(k)) })
}

func (t *Trie) Log() {
	t.root.log(0)
}

func (t *Trie) Count() int {
	n := 0
	t.Iterate(func(_ []byte, _ interface{}) { n++ })
	return n
}

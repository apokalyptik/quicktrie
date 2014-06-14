package trie

import (
	"fmt"
	"log"
	"strings"
)

var Debug bool

type Trie struct {
	key      []byte
	children [256]*Trie
	endpoint uint8
}

func NewTrie() *Trie {
	return &Trie{
		children: [256]*Trie{},
	}
}

func (t *Trie) longestCommonPrefix(a, b []byte) int {
	var i = 0
	for ; i < len(b) && i < len(a) && b[i] == a[i]; i++ {
	}
	return i
}

func (t *Trie) AddString(key string) {
	if t.key == nil {
		t.Add([]byte(key))
	}
}

func (t *Trie) Add(key []byte) {
	if t.key == nil {
		if Debug {
			log.Printf("add: %s\n", string(key))
		}
		t.add([]byte(key))
	}
}

func (t *Trie) add(key []byte) {
	index := uint8(key[0])
	klen := len(key)
	if t.children[index] == nil {
		// We can just add our new trie node here
		t.children[index] = &Trie{
			key:      key,
			endpoint: 1,
		}
	} else {
		// The node we want to act on exists...
		child := t.children[index]
		clen := len(child.key)
		lcp := t.longestCommonPrefix(key, child.key)
		if lcp == klen && lcp == clen {
			// This key exists exactly. eg: have "aa", adding "aa"
			child.endpoint = 1
		} else if lcp == clen {
			// the entire new key is a prefix of the existing child key eg: have "aaa", adding "aa"
			child.add(key[lcp:])
		} else if lcp == klen {
			// the entire child key is a prefix for the new key. eg: have "aa", adding "aa.+"
			//		1. this becomes the new child node
			//		2. the old child node is rekeyed and becomes a child of the new one
			oldChild := t.children[index]
			oldChild.key = oldChild.key[lcp:]
			newChild := &Trie{
				key:      key[lcp:],
				endpoint: 1,
			}
			newChild.children[uint8(oldChild.key[0])] = oldChild
			t.children[index] = newChild
		} else {
			// the child key shares a prefix with the new key, both are longer than the prefix
			// eg: have "abc", adding "abd"
			//		1. create a new child node with the common prefix as the key
			//		2. the new child node gets a child node for the unique part of the new key
			//		3. the new child node gets a child node for the unique part of the old child node key
			newChild := &Trie{
				key: key[:lcp],
			}
			child.key = child.key[lcp:]
			newChild.children[uint8(child.key[0])] = child
			newChild.children[uint8(key[lcp])] = &Trie{
				key:      key[lcp:],
				endpoint: 1,
			}
			t.children[index] = newChild
		}
	}
}

func (t *Trie) DropString(key string) {
	t.Drop([]byte(key))
}

func (t *Trie) Drop(key []byte) {
	if t.key == nil {
		if Debug {
			log.Printf("drop: " + string(key))
			t.Log()
		}
		t.drop(key)
		if Debug {
			t.Log()
		}
	}
}

func (t *Trie) drop(key []byte) {
	index := uint8(key[0])
	if t.children[index] == nil {
		// nothing to drop
		return
	}

	child := t.children[index]
	klen := len(key)
	//clen := len(child.key)
	lcp := t.longestCommonPrefix(key, child.key)

	if lcp == klen {
		// The key we're dropping is exactly the child key or longer. Since drop is
		// recursive we want to drop the child eg: have "abcd", dropping "abcd"
		t.children[index] = nil
		return
	}

	// The key we're dropping shares a common prefix with the entire child key, recurse
	// eg: have "abc", dropping "abcdef"
	child.drop(key[lcp:])
}

func (t *Trie) DelString(key string) {
	t.Del([]byte(key))
}

func (t *Trie) Del(key []byte) {
	if t.key == nil {
		if Debug {
			log.Printf("del: %s\n", string(key))
		}
		t.del(key)
	}
}

func (t *Trie) del(key []byte) {
	index := uint8(key[0])
	if t.children[index] == nil {
		// We have reached the end of this branch of the trie. Nothing to delete
		return
	}
	child := t.children[index]
	clen := len(child.key)
	lcp := t.longestCommonPrefix(key, child.key)

	if lcp < clen {
		// If the key here is longer than the key we're deleting then the key we're
		// deleting cannot exist at this node. Nothing to delete. eg: have "abcd" del "ab"
		return
	}

	klen := len(key)

	if lcp == klen {
		// The key we're deleting is exactly the child key (we made sure lcp == clen above)
		// eg: have "abcd", del "abcd"
		child.endpoint = 0

		// TODO: this could be optimized away with a uint8 counter
		for _, v := range child.children {
			if v != nil {
				// The child has at least one non nil child of its own. We don't want to
				// prune the tree here...
				return
			}
		}

		// The child is not an endpoint and has no non-nil children of its own. Prune
		t.children[index] = nil
		return
	}

	// The key we're deleting shares a common prefix with the entire child key, recurse
	// eg: have "abc", del "abcdef"
	child.del(key[lcp:])

	if child.endpoint == 1 {
		// This child is an endpoint, we don't want to prune the tree here
		return
	}

	// TODO: this could be optimized away with a uint8 counter
	for _, v := range child.children {
		if v != nil {
			// The child has at least one non nil child of its own. We don't want to
			// prune the tree here...
			return
		}
	}

	// The child is not an endpoint and has no non-nil children of its own. Prune
	t.children[index] = nil
}

func (t *Trie) GetString(key string) bool {
	return t.Get([]byte(key))
}

func (t *Trie) Get(key []byte) bool {
	if t.key == nil {
		if Debug {
			log.Printf("get: %s\n", string(key))
		}
		return t.get(key)
	}
	return false
}

func (t *Trie) get(key []byte) bool {
	index := uint8(key[0])

	if t.children[index] == nil {
		// branch ends here with nothing to show for our work
		return false
	}

	child := t.children[index]
	klen := len(key)
	clen := len(child.key)
	lcp := t.longestCommonPrefix(key, child.key)

	if lcp < clen {
		// If the key here is longer than the key we're wanting then the key we're
		// wanting cannot exist at this node. Nothing to find. eg: have "abcd" want "ab"
		return false
	}

	if lcp == klen {
		// The key we're wanting is exactly the child key (we made sure lcp == clen above)
		// eg: have "abcd", want "abcd"
		return child.endpoint == 1
	}

	// The key we're wanting shares a common prefix with the entire child key, recurse
	// eg: have "abc", want "abcdef"
	return child.get(key[lcp:])
}

func (t *Trie) IterateString(callback func(string)) {
	if t.key == nil {
		t.Iterate(func(b []byte) { callback(string(b)) })
	}
}

func (t *Trie) Iterate(callback func([]byte)) {
	if t.key == nil {
		t.iterate([]byte{}, callback)
	}
}

func (t *Trie) iterate(key []byte, callback func([]byte)) {
	if t.endpoint != 0 {
		callback(key)
	}
	for _, v := range t.children {
		if v != nil {
			v.iterate(append(key, v.key...), callback)
		}
	}
}

func (t *Trie) Print() {
	for _, v := range t.children {
		if v != nil {
			fmt.Println(string(v.key))
			t.iterate([]byte{}, func(k []byte) { fmt.Println(string(k)) })
		}
	}
}

func (t *Trie) Log(indent ...int) {
	var indentLevel = len(indent)
	if indentLevel > 0 {
		indentLevel = indent[0]
	}
	var prefix = strings.Repeat(" ", indentLevel) + "↪-"
	if indentLevel > 0 {
		prefix = "|" + prefix
	}
	if t.key != nil {
		if t.endpoint == 0 {
			log.Printf("  %s%s\n", prefix, string(t.key))
		} else {
			log.Printf("* %s%s\n", prefix, string(t.key))
		}
		indentLevel++
	} else {
		log.Printf("  ++\n")
	}
	for _, v := range t.children {
		if v != nil {
			v.Log(indentLevel)
		}
	}
}

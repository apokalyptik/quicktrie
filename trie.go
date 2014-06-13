package trie

import (
	"fmt"
	"log"
	"strings"
)

var Debug bool

type Trie struct {
	key      []byte
	children []*Trie
	endpoint *struct{}
}

func NewTrie() *Trie {
	return &Trie{
		children: []*Trie{},
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
	for k, v := range t.children {
		if lcp := t.longestCommonPrefix(t.children[k].key, key); lcp > 0 {
			if lcp == len(key) && lcp == len(v.key) {
				// This key exists exactly
				// eg: have "aa", adding "aa"
				if v.endpoint == nil {
					v.endpoint = &struct{}{}
				}
			} else if lcp == len(key) {
				// the entire key is a sub-key of the child key
				// eg: have "aa", adding "aaa"
				oldChild := v
				oldChild.key = v.key[lcp:]
				newChild := &Trie{
					endpoint: &struct{}{},
					key:      key[:lcp],
					children: []*Trie{oldChild},
				}
				t.children[k] = newChild
			} else if lcp == len(v.key) {
				// the entire child key is a prefix for the key
				// eg: have "aa", adding "aaa"
				v.add(key[lcp:])
			} else {
				// the key and child key share a common prefix but are both going to
				// end up as their own children of the common prefix on account of
				// being larger than said prefix
				// eg: have "abc", adding "ayz"
				oldChild := v
				oldChild.key = oldChild.key[lcp:]
				newChild := &Trie{
					key: key[:lcp],
					children: []*Trie{
						oldChild,
						&Trie{
							endpoint: &struct{}{},
							key:      key[lcp:],
						},
					},
				}
				t.children[k] = newChild
			}
			return
		}
	}
	t.children = append(t.children, &Trie{key: key, endpoint: &struct{}{}})
}

func (t *Trie) DropString(key string) {
	t.Drop([]byte(key))
}

func (t *Trie) Drop(key []byte) {
	if t.key == nil {
		t.drop(key)
	}
}

func (t *Trie) drop(key []byte) {
	if t.key == nil && key == nil {
		t.children = []*Trie{}
		return
	}
	for k, v := range t.children {
		if lcp := t.longestCommonPrefix(key, v.key); lcp > 0 {
			if lcp == len(key) {
				if k == 0 {
					t.children = append(t.children[1:])
				} else {
					t.children = append(t.children[:k], t.children[k+1:]...)
				}
			} else if lcp < len(v.key) {
				// The delete key is part of, but less than the entire child key, Cannot be an exact match
				return
			} else {
				// The child key is less than the delete key but it is a prefix of the delete key, recurse
				v.drop(key[lcp:])
				if v.endpoint == nil && len(v.children) == 0 {
					if k == 0 {
						t.children = append(t.children[1:])
					} else {
						t.children = append(t.children[:k], t.children[k+1:]...)
					}
				}
			}
		}
	}
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
	for k, v := range t.children {
		if lcp := t.longestCommonPrefix(key, v.key); lcp > 0 {
			if lcp < len(v.key) {
				// Delete key is a prefix of child, but is not child or further, nothing to do
				return
			}
			if lcp == len(key) && lcp == len(v.key) {
				// This is the key we came for
				v.endpoint = nil
				if len(v.children) == 0 {
					if k == 0 {
						t.children = append(t.children[1:])
					} else {
						t.children = append(t.children[:k], t.children[k+1:]...)
					}
				}
				return
			}
			if lcp == len(v.key) {
				v.del(key[lcp:])
				if v.endpoint == nil && len(v.children) == 0 {
					if k == 0 {
						t.children = append(t.children[1:])
					} else {
						t.children = append(t.children[:k], t.children[k+1:]...)
					}
				}
				return
			}
		}
	}
	// No such key found in the tree
}

func (t *Trie) GetString(key string) bool {
	return t.Get([]byte(key))
}

func (t *Trie) Get(key []byte) bool {
	if t.key == nil {
		if Debug {
			log.Printf("get: %s\n", string(key))
		}
		rc := make(chan bool, len(t.children))
		dc := make(chan struct{})
		gr := 0
		for _, v := range t.children {
			gr++
			go func(v *Trie) {
				var result bool
				if key[0] == v.key[0] {
					lcp := t.longestCommonPrefix(key, v.key)
					if lcp == len(key) && lcp == len(v.key) {
						result = true
					} else {
						result = v.get(key[lcp:])
					}
				} else {
				}
				select {
				case <-dc:
					return
				default:
					rc <- result
				}
			}(v)
		}
		for i := 0; i < gr; i++ {
			if <-rc {
				close(dc)
				return true
			}
		}
	}
	return false
}

func (t *Trie) get(key []byte) bool {
	for _, v := range t.children {
		if key[0] != v.key[0] {
			continue
		}
		if len(key) < len(v.key) {
			continue
		}
		// The following can only be true if the preceeding was true
		if lcp := t.longestCommonPrefix(key, v.key); lcp > 0 {
			if lcp < len(v.key) {
				// if the common prefix less than the entirety of the child
				// key, it cannot possibly match the child key
				return false
			}
			if lcp == len(key) {
				// the child is exactly the key we're looking for
				if v.endpoint != nil {
					return true
				} else {
					return false
				}
			}
			return v.get(key[lcp:])
		}
	}
	return false
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
	if t.endpoint != nil {
		callback(key)
	}
	for _, v := range t.children {
		v.iterate(append(key, v.key...), callback)
	}
}

func (t *Trie) Print() {
	for _, v := range t.children {
		fmt.Println(string(v.key))
		t.iterate([]byte{}, func(k []byte) { fmt.Println(string(k)) })
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
		if t.endpoint == nil {
			log.Printf("  %s%s\n", prefix, string(t.key))
		} else {
			log.Printf("* %s%s\n", prefix, string(t.key))
		}
		indentLevel++
	} else {
		log.Printf("  ++\n")
	}
	for _, v := range t.children {
		v.Log(indentLevel)
	}
}

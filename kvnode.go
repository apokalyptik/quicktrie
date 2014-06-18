package trie

import (
	"log"
	"strings"
)

type kvTrie struct {
	key      []byte
	value    interface{}
	children []*kvTrie
	endpoint uint8
}

func (t *kvTrie) add(key []byte, vals ...interface{}) {
	if len(vals) < 1 {
		vals = []interface{}{nil}
	}
	for k, v := range t.children {
		if lcp := longestCommonPrefix(t.children[k].key, key); lcp > 0 {
			if lcp == len(key) && lcp == len(v.key) {
				// This key exists exactly
				// eg: have "aa", adding "aa"
				if v.endpoint == 0 {
					v.endpoint = 1
				}
			} else if lcp == len(key) {
				// the entire key is a sub-key of the child key
				// eg: have "aa", adding "aaa"
				oldChild := v
				oldChild.key = v.key[lcp:]
				newChild := &kvTrie{
					endpoint: 1,
					key:      key[:lcp],
					value:    vals[0],
					children: []*kvTrie{oldChild},
				}
				t.children[k] = newChild
			} else if lcp == len(v.key) {
				// the entire child key is a prefix for the key
				// eg: have "aa", adding "aaa"
				v.add(key[lcp:], vals...)
			} else {
				// the key and child key share a common prefix but are both going to
				// end up as their own children of the common prefix on account of
				// being larger than said prefix
				// eg: have "abc", adding "ayz"
				oldChild := v
				oldChild.key = oldChild.key[lcp:]
				newChild := &kvTrie{
					key: key[:lcp],
					children: []*kvTrie{
						oldChild,
						&kvTrie{
							endpoint: 1,
							value:    vals[0],
							key:      key[lcp:],
						},
					},
				}
				t.children[k] = newChild
			}
			return
		}
	}
	t.children = append(t.children, &kvTrie{key: key, endpoint: 1, value: vals[0]})
}

func (t *kvTrie) drop(key []byte) {
	if t.key == nil && key == nil {
		t.children = []*kvTrie{}
		return
	}
	for k, v := range t.children {
		if lcp := longestCommonPrefix(key, v.key); lcp > 0 {
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
				if v.endpoint == 0 && len(v.children) == 0 {
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

func (t *kvTrie) del(key []byte) {
	for k, v := range t.children {
		if lcp := longestCommonPrefix(key, v.key); lcp > 0 {
			if lcp < len(v.key) {
				// Delete key is a prefix of child, but is not child or further, nothing to do
				return
			}
			if lcp == len(key) && lcp == len(v.key) {
				// This is the key we came for
				v.endpoint = 0
				v.value = nil
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
				if v.endpoint == 0 && len(v.children) == 0 {
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

func (t *kvTrie) get(key []byte) (bool, interface{}) {
	for _, v := range t.children {
		if key[0] != v.key[0] {
			continue
		}
		if len(key) < len(v.key) {
			continue
		}
		// The following can only be true if the preceeding was true
		if lcp := longestCommonPrefix(key, v.key); lcp > 0 {
			if lcp < len(v.key) {
				// if the common prefix less than the entirety of the child
				// key, it cannot possibly match the child key
				return false, nil
			}
			if lcp == len(key) {
				// the child is exactly the key we're looking for
				if v.endpoint != 0 {
					return true, v.value
				}
				return false, nil
			}
			return v.get(key[lcp:])
		}
	}
	return false, nil
}

func (t *kvTrie) iterate(key []byte, callback IterFunc) {
	if t.endpoint != 0 {
		callback(key, t.value)
	}
	for _, v := range t.children {
		v.iterate(append(key, v.key...), callback)
	}
}

func (t *kvTrie) iterateFrom(prefix []byte, callback IterFunc) {
	for _, v := range t.children {
		if len(prefix) == 0 {
			v.iterate(v.key, callback)
			continue
		}

		if prefix[0] != v.key[0] {
			// This child key cannot be prefixed by the prefix argument
			continue
		}
		// This child key must be prefixed by the prefix argument

		lcp := longestCommonPrefix(prefix, v.key)

		if lcp == len(prefix) {
			// the child key is entirely prefixed by the prefix argument
			v.iterate(v.key, callback)
		} else if lcp == len(v.key) {
			// the entire child key is a shared sub prefix of the prefix argument
			// time to recurse
			v.iterateFrom(prefix[lcp:], callback)
		}

		return
	}
}

func (t *kvTrie) log(indent ...int) {
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
		v.log(indentLevel)
	}
}

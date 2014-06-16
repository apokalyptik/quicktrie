package trie

import (
	"log"
	"strings"
)

type bwTrie struct {
	key      []byte
	children []*bwTrie
	endpoint uint8
}

func (t *bwTrie) add(key []byte, _ ...interface{}) {
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
				newChild := &bwTrie{
					endpoint: 1,
					key:      key[:lcp],
					children: []*bwTrie{oldChild},
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
				newChild := &bwTrie{
					key: key[:lcp],
					children: []*bwTrie{
						oldChild,
						&bwTrie{
							endpoint: 1,
							key:      key[lcp:],
						},
					},
				}
				t.children[k] = newChild
			}
			return
		}
	}
	t.children = append(t.children, &bwTrie{key: key, endpoint: 1})
}

func (t *bwTrie) drop(key []byte) {
	if t.key == nil && key == nil {
		t.children = []*bwTrie{}
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

func (t *bwTrie) del(key []byte) {
	for k, v := range t.children {
		if lcp := longestCommonPrefix(key, v.key); lcp > 0 {
			if lcp < len(v.key) {
				// Delete key is a prefix of child, but is not child or further, nothing to do
				return
			}
			if lcp == len(key) && lcp == len(v.key) {
				// This is the key we came for
				v.endpoint = 0
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

func (t *bwTrie) get(key []byte) (bool, interface{}) {
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
					return true, nil
				} else {
					return false, nil
				}
			}
			return v.get(key[lcp:])
		}
	}
	return false, nil
}

func (t *bwTrie) iterate(key []byte, callback IterFunc) {
	if t.endpoint != 0 {
		callback(key, nil)
	}
	for _, v := range t.children {
		v.iterate(append(key, v.key...), callback)
	}
}

func (t *bwTrie) log(indent ...int) {
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

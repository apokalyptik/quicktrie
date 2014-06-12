package trie

import "sync"

type Trie struct {
	source     *Trie
	leaf       bool
	entries    int
	childCount int
	children   map[uint8]*Trie
	lock       sync.RWMutex
}

func NewTrie() *Trie {
	return &Trie{
		source:   nil,
		children: map[uint8]*Trie{},
	}
}

func (t *Trie) Add(key string) bool {
	if t.source == nil {
		t.add(key, t)
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.source == nil
}

func (t *Trie) add(key string, source *Trie) {
	if len(key) == 0 {
		source.entries++
		t.leaf = true
		return
	} else {
		if v, ok := t.children[key[0]]; ok {
			v.add(key[1:], source)
		} else {
			v := &Trie{
				source:   source,
				children: map[uint8]*Trie{},
			}
			t.children[key[0]] = v
			t.childCount++
			v.add(key[1:], source)
		}
	}
}

func (t *Trie) Del(key string) {
	if t.source == nil {
		t.del(key, t)
	}
	t.lock.Lock()
	defer t.lock.Unlock()
}

func (t *Trie) del(key string, source *Trie) int {
	if len(key) == 0 {
		if t.leaf {
			source.entries--
			t.leaf = false
		}
	} else {
		if v, ok := t.children[key[0]]; ok {
			v.del(key[1:], source)
			if v.childCount == 0 && t.leaf == false {
				delete(t.children, key[0])
				t.childCount--
			}
		}
	}
	return t.childCount
}

func (t *Trie) Get(key string) bool {
	if t.source != nil {
		return false
	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.get(key)
}

func (t *Trie) get(key string) bool {
	if len(key) == 0 {
		return t.leaf
	}
	if t.childCount == 0 {
		return false
	}
	if v, ok := t.children[key[0]]; ok {
		return v.get(key[1:])
	} else {
		return false
	}
}

func (t *Trie) Iterate(callback func(string)) {
	if t.source != nil {
		return
	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.childCount > 0 {
		for k, v := range t.children {
			v.iterate(string(k), callback)
		}
	}
}

func (t *Trie) iterate(key string, callback func(string)) {
	if t.leaf {
		callback(key)
	}
	if t.childCount > 0 {
		for k, v := range t.children {
			v.iterate(key+string(k), callback)
		}
	}
}

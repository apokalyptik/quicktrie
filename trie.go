package trie

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

type Trie struct {
	root node
}

func NewTrie() *Trie {
	return &Trie{
		root: &bwTrie{
			children: []*bwTrie{},
		},
	}
}

func NewKVTrie() *Trie {
	return &Trie{
		root: &kvTrie{
			children: []*kvTrie{},
		},
	}
}

func (t *Trie) Add(key interface{}, data ...interface{}) {
	switch key := key.(type) {
	case []byte:
		t.root.add(key, data...)
	case string:
		t.root.add([]byte(key), data...)
	}
}

func (t *Trie) Drop(key interface{}) {
	switch key := key.(type) {
	case []byte:
		t.root.drop(key)
	case string:
		t.root.drop([]byte(key))
	}
}

func (t *Trie) Del(key interface{}) {
	switch key := key.(type) {
	case []byte:
		t.root.del(key)
	case string:
		t.root.del([]byte(key))
	}
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

func (t *Trie) Get(key interface{}) (bool, interface{}) {
	switch key := key.(type) {
	case []byte:
		return t.root.get(key)
	case string:
		return t.root.get([]byte(key))
	default:
		return false, nil
	}
}

func (t *Trie) Iterate(callback IterFunc) {
	t.root.iterate([]byte{}, callback)
}

func (t *Trie) Log() {
	t.root.log(0)
}

func (t *Trie) Count() int {
	n := 0
	t.Iterate(func(_ []byte, _ interface{}) { n++ })
	return n
}

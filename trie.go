/*
Trie presents a simple, clean, black and white radix trie interface. For more
information on what kind of data structure a radix trie is, please see
http://en.wikipedia.org/wiki/Radix_tree.

Two types of data structures are available. A black and white only tree
(default) which merely records the existence of keys added to the data
structure, and key-value black and white trie which stores an arbitrary
value in the trie attached to the trie.  Because the value is stored
seperately from the existence of the key you can store nil values in the trie
and be able to tell that apart from a key that does not exist.

The data structures inside this package are NOT synchronized, you'll want to
add a sync.Mutex or sync.RWMutex to your code if it needs to be thread safe.
At this point in time I do consider the structures to be *mostly* safe to use
concurrently without locking as long as you're willing to end up in a part of
the trie that existed when you started making your call but not by the time
the call has ended. For more read operations this is "OK" (please know the
gaurantees required by your application) but it is absolutely not OK for write
operations and you should absolutely synchronize those or risk fairly massive
corruption.

Example BW Trie usage
	trie := trie.NewBWTrie
	trie.Add("apple")
	trie.Add("apply")
	trie.Add("actually")
	trie.Add("actively")
	trie.Add("Alaska")
	trie.Get("actively") // true
	trie.Get("ancillary") // false
	trie.Log()
		// 2014/06/16 12:06:29   ++
		// 2014/06/16 12:06:29   ↪-a
		// 2014/06/16 12:06:29   | ↪-ppl
		// 2014/06/16 12:06:29 * |  ↪-e
		// 2014/06/16 12:06:29 * |  ↪-y
		// 2014/06/16 12:06:29   | ↪-ct
		// 2014/06/16 12:06:29 * |  ↪-ually
		// 2014/06/16 12:06:29 * |  ↪-ively
		// 2014/06/16 12:06:29 * ↪-Alaska
	trie.Drop("ap")
	trie.Log()
		// 2014/06/16 12:08:07   ++
		// 2014/06/16 12:08:07   ↪-a
		// 2014/06/16 12:08:07   | ↪-ct
		// 2014/06/16 12:08:07 * |  ↪-ually
		// 2014/06/16 12:08:07 * |  ↪-ively
		// 2014/06/16 12:08:07 * ↪-Alaska
*/
package trie

// IterFunc describes the function signature that it required for the callback
// portion of the Iterate function
type IterFunc func([]byte, interface{})

type node interface {
	get([]byte) (bool, interface{})
	add([]byte, ...interface{})
	del([]byte)
	drop([]byte)
	iterate([]byte, IterFunc)
	log(...int)
}

// Trie is the itnerface to your requested trie. This is the interface you'll
// use whether you requested a BW trie or a KV trie.
type Trie struct {
	root node
}

// NewTrie is a convenience function, it merely calls NewBWTrie. Please see the
// documentation for NewBWTrie
func NewTrie() *Trie {
	return NewBWTrie()
}

// NewBWTrie returns a new "black and white" radix trie.  A black and white radix
// trie records the existence of a string, but there is no data associated with
// it.  The existence is the data. This kind of trie is useful for simple
// "exists" stype lookups.
func NewBWTrie() *Trie {
	return &Trie{
		root: &bwTrie{
			children: []*bwTrie{},
		},
	}
}

// NewKVTrie returns a new Key/Value trie.  This is essentially a black and
// white radix trie where each node has arbitrary data associated with it.
func NewKVTrie() *Trie {
	return &Trie{
		root: &kvTrie{
			children: []*kvTrie{},
		},
	}
}

// Add allows you to add a key to your trie.  For BW tries the data argument
// is ignored and you may ommit it.  For KV tries you may pass any value that
// you wish, and it will be stored along with your key.  If you choose to
// ommit the data value for a KV trie then the stored data will be nil. Only
// the first data argument is recognized, so if you wish to store an array or
// slice with your key you should pass that, not depend on the variadic
func (t *Trie) Add(key interface{}, data ...interface{}) {
	switch key := key.(type) {
	case []byte:
		t.root.add(key, data...)
	case string:
		t.root.add([]byte(key), data...)
	}
}

// Drop allows you to cut an enitre branch off of your trie.  This means that
// every existing key of which the passed key is a prefix string (inclusive)
// will be removed from the trie.
func (t *Trie) Drop(key interface{}) {
	switch key := key.(type) {
	case []byte:
		t.root.drop(key)
	case string:
		t.root.drop([]byte(key))
	}
}

// Del allows you to remove a single key from the trie. Del will only delete
// the exactly matching key, unlike drop, and is therefor considerably safer
// unless you know why you would want to drop an entire prefix from your trie
func (t *Trie) Del(key interface{}) {
	switch key := key.(type) {
	case []byte:
		t.root.del(key)
	case string:
		t.root.del([]byte(key))
	}
}

// Exists allows you to check the existence of a key within the trie.  As it
// only returns asingle boolean value it's convenient to use inside of if
// statements
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

// Get allows you to fetch a key from the trie.  If the trie is a BW trie then
// get returns true, or false depending on whether the key has been added to
// the trie and nil for the data.  KV tries will return the data passed to add
// as the second return value
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

// Iterate allows you to run a function against every key inserted into the
// trie.
func (t *Trie) Iterate(callback IterFunc) {
	t.root.iterate([]byte{}, callback)
}

// Log prints a "pretty" representation of the trie. This is mainly useful for
// debugging, and it'll print arbitrary binary data if you've used something
// other than text strings for your keys
func (t *Trie) Log() {
	t.root.log(0)
}

// Count returns the number of keys in the trie.  Internally it uses the
// Iterate function to do this
func (t *Trie) Count() int {
	n := 0
	t.Iterate(func(_ []byte, _ interface{}) { n++ })
	return n
}

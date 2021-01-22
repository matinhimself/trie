package hashtable

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/matinhimself/trie/pkg/trie"
	"math"
	"os"
	"sync"
)

type node struct {
	Value HashAble
}

func (n node) String() string {
	return fmt.Sprintf("%s", n.Value)
}

type HashAble interface {
	Equals(other *HashAble) bool
	ToHash() uint64
	GetKey() string
}

// HashTable is a wrapper for a trie tree
// and a hashtable. it will store each
// student in a hashtable and its
// hash value in a trie tree as a pair of
// hash value and student id.
type HashTable struct {
	lock    sync.RWMutex
	size    int
	count   int
	buckets [][]node
	tree    *trie.Trie
}

func (hm *HashTable) Size() int {
	return hm.size
}

func NewHashTable(size int) (*HashTable, error) {
	hm := new(HashTable)
	if size <= 0 {
		return nil, errors.New("hashmap size should be > 1")
	}

	hm.buckets = make([][]node, size)
	hm.size = size
	hm.count = 0
	hm.tree = trie.NewTrie()
	for i := range hm.buckets {
		hm.buckets[i] = make([]node, 0, 20)
	}
	return hm, nil
}

// Safe multiplication for indexing
// in situation of overflow it will always
// return positive number and a false and
// false boolean.
func Multi64(u1, u2 uint64) (res uint64, ok bool) {
	if u2 >= (math.MaxUint64 / u1) {
		u3 := u1 * u2
		if u3 < 0 {
			u3 = -u3
		}
		return u3, false
	} else {
		return u1 * u2, true
	}
}

// Implements the Fowler–Noll–Vo hash function
// Tc: O(m) with m as length of a key, But for large number of entries,
// length of the keys is almost negligible. so hash computation can be considered
// to take place in constant time O(1).
func hash(key string) uint64 {
	var h uint64
	for _, b := range []byte(key) {
		h = h ^ uint64(b)
		tmp, _ := Multi64(h, 1099511628211)
		h = tmp
	}
	return h
}

// Implements the Jenkins hash function
//func Jenkins(key string) uint64 {
//	var h uint64
//	for _, c := range key{
//		h += uint64(c)
//		h += h << 10
//		h ^= h >> 6
//	}
//	h += h << 3
//	h ^= h >> 11
//	h += h << 15
//
//	return h
//}

// Implements mine algorithm
//func mineAlgo(key string) uint32 {
//	var h float64
//	for i := 0; i < len(key); i++ {
//		h += math.Pow(97,float64(i)) * float64(key[i])
//	}
//	h = math.Mod(h, float64(4999))
//
//	return uint32(h)
//}

// getIndex hashes the given key and returns its index in
// buckets array.
func (hm *HashTable) getIndex(key HashAble) uint64 {
	rn := key.ToHash() % uint64(hm.size)
	return rn
}

// Set sets the value for an associated key in the hashmap.
// given object should implements HashAble interface.
func (hm *HashTable) Set(obj HashAble) uint64 {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	index := hm.getIndex(obj)
	chain := hm.buckets[index]
	found := false

	// check if the key already exists
	for i := range chain {
		// if found, update the node
		node := &chain[i]
		if node.Value.Equals(&obj) {
			node.Value = obj
			found = true
		}
	}
	if found { // hashmap has been updated
		return index
	}

	// add a new node
	node := node{Value: obj}
	chain = append(chain, node)
	hm.buckets[index] = chain
	hm.count++
	hm.tree.Insert(obj.GetKey(), index)
	return index
}

// GetAllKeys returns all keys stored in the trie.
func (hm *HashTable) GetAllKeys() []string {
	return hm.tree.GetAllKeys()
}

// Prints all keys stored in the trie.
func (hm *HashTable) PrintAll() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Student ID"})
	hm.lock.RLock()
	defer hm.lock.RUnlock()

	res := hm.tree.GetAllKeys()
	for _, re := range res {
		t.AppendRow(table.Row{
			re,
		})
	}
	t.SetAutoIndex(true)
	t.SetStyle(table.StyleLight)
	t.Render()
}

// Get returns the value associated with a key in the hashTable,
// and an boolean indicating whether the value exists or not.
func (hm *HashTable) Get(studentId string) (*node, bool) {
	hm.lock.RLock()
	defer hm.lock.RUnlock()

	val, found := hm.tree.Search(studentId)
	if !found || val == nil {
		return nil, false
	}
	index := (*val).(uint64)
	chain := hm.buckets[index]
	for _, node := range chain {
		if node.Value.GetKey() == studentId {
			return &node, true
		}
	}
	return nil, false
}

// Delete deletes the Node associated with a key in the hashTable,
// and an boolean indicating whether the it was successfully deleted
// or not.
func (hm *HashTable) Delete(studentId string) (deleted bool) {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	ind, deleted := hm.tree.Delete(studentId)
	if !deleted || *ind == nil {
		return false
	}
	index := (*ind).(uint64)
	chain := hm.buckets[index]
	for i, node := range chain {
		if node.Value.GetKey() == studentId {
			hm.buckets[index] = append(chain[:i], chain[i+1:]...)
			return true
		}
	}
	return false
}

// GetKeysWithPrefix returns all keys exiting with a given prefix
func (hm *HashTable) GetKeysWithPrefix(studentId string) []string {
	hm.lock.RLock()
	defer hm.lock.RUnlock()

	keys := hm.tree.GetPrefixKeys(studentId)
	return keys
}

type pair struct {
	Key string
	Value HashAble
}


func (hm *HashTable) GetPairsWithPrefix(pref string) []pair {
	res := hm.GetKeysWithPrefix(pref)
	pairs := make([]pair, 0)
	for _, re := range res {
		elem, found := hm.Get(re)
		if found {
			pairs = append(
				pairs,
				pair{re, elem.Value},
			)
		}
	}
	return pairs
}

func (hm *HashTable) GetAllPairs() []pair {
	res := hm.GetAllKeys()
	pairs := make([]pair, 0)
	for _, re := range res {
		elem, found := hm.Get(re)
		if found {
			pairs = append(
				pairs,
				pair{re, elem.Value},
			)
		}
	}
	return pairs
}

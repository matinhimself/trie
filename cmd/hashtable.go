package main

import (
	"errors"
	"fmt"
	"github.com/matinhimself/trie/pkg/trie"
	"math"
	"sync"
)

type Node struct {
	Value *Student
}

func (n Node) String() string {
	return fmt.Sprintf("%s", n.Value)
}

// HashTable is a wrapper for a trie tree
// and a hashtable with will store each
// student in a hashtable and store its
// hash value in a trie tree as a pair of
// hash value and student id
type HashTable struct {
	size    int
	count   int
	buckets [][]Node
	tree    *trie.Trie
	lock    sync.RWMutex
}

func (hm *HashTable) Size() int {
	return hm.size
}

func NewHashTable(size int) (*HashTable, error) {
	hm := new(HashTable)
	if size <= 0 {
		return nil, errors.New("hashmap size should be > 1")
	}

	hm.buckets = make([][]Node, size)
	hm.size = size
	hm.count = 0
	hm.tree = trie.NewTrie()
	for i := range hm.buckets {
		hm.buckets[i] = make([]Node, 0, 20)
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

// Picking hash function:
// I was confused about picking hash function
// that is suitable for storing something like
// student id, so i wrote some tests for
// different hash functions,
// and here is the stats of best two hash functions:
// hashFunction       loadFactor       standardDeviation       emptyCells
//   Jenkins              4.0                1.77					2
//   Jenkins              6.0                2.27					0
//   Jenkins              8.0                2.60					0
// Fowler–Noll–Vo         4.0                1.47					0
// Fowler–Noll–Vo         6.0                1.870					0
// Fowler–Noll–Vo         8.0                2.20					0

// Implements the Fowler–Noll–Vo hash function
// Tc: O(m) with m as length of a key
// But for large number of entries,
// length of the keys is almost negligible.
// so hash computation can be considered
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

// returns the index at which the key needs to go
func (hm *HashTable) getIndex(key *Student) uint32 {
	rn := key.toHash() % uint32(hm.size)
	return rn
}

// Set the value for an associated key in the hashmap
func (hm *HashTable) Set(student *Student) uint32 {
	index := hm.getIndex(student)
	chain := hm.buckets[index]
	found := false

	// first see if the key already exists
	for i := range chain {
		// if found, update the student
		node := &chain[i]
		if node.Value.StudentID == student.StudentID {
			node.Value = student
			found = true
		}
	}
	if found { // hashmap has been updated
		return index
	}

	// add a new node
	node := Node{Value: student}
	chain = append(chain, node)
	hm.buckets[index] = chain
	hm.count++
	hm.tree.Insert([]byte(student.StudentID), index)
	return index
}

func (hm *HashTable) printAll() {
	res := hm.tree.GetAllKeys()
	for i, re := range res {
		fmt.Printf("%04d.", i)
		for _, b := range re {
			fmt.Print(string(b+byte('0')))
		}
		fmt.Print("\n")
	}
}

// Get returns the value associated with a key in the hashTable,
// and an error indicating whether the value exists or not.
func (hm *HashTable) Get(studentId StudentID) (*Node, bool) {
	val, found := hm.tree.Search(trie.Bytes(studentId))
	if !found || val == nil{
		return nil, false
	}
	index := (*val).(uint32)
	chain := hm.buckets[index]
	for _, node := range chain {
		if node.Value.StudentID == studentId {
			return &node, true
		}
	}
	return nil, false
}

func (hm *HashTable) Delete(studentId StudentID) (deleted bool) {
	ind, deleted := hm.tree.Delete(trie.Bytes(studentId))
	if !deleted || *ind == nil{
		return false
	}
	index := (*ind).(uint32)
	chain := hm.buckets[index]
	for i, node := range chain {
		if node.Value.StudentID == studentId {
			chain[i] = Node{}
			return true
		}
	}
	return false
}




func (hm *HashTable) GetKeysWithPrefix(studentId StudentID) []string {
	keys := hm.tree.GetPrefixKeys(trie.Bytes(studentId))
	var strs []string
	for _, key := range keys {
		if len(key) >= len(studentId) {
			tmp := ""
			for _, b := range key {
				tmp += string(b + byte('0'))
			}
			strs = append(strs, tmp)
		}
	}
	return strs
}


func (hm *HashTable) GetAllKeys(studentId StudentID) []string {
	keys := hm.tree.GetPrefixKeys(trie.Bytes(studentId))
	var strs []string
	for _, key := range keys {
		if len(key) >= len(studentId) {
			tmp := ""
			for _, b := range key {
				tmp += string(b + byte('0'))
			}
			strs = append(strs, tmp)
		}
	}
	return strs
}

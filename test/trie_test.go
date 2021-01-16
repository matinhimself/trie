package test

import (
	"bytes"
	"fmt"
	"github.com/matinhimself/trie/pkg/trie"
	"strconv"
	"testing"
)

func TestTrieAdd(t *testing.T) {
	tree := trie.NewTrie()
	for i := 0; i < 1000; i++ {
		tree.Insert(strconv.Itoa(i), i)
	}
	if tree.Size() != 1000 {
		t.Error("tree size is wrong")
	}
	for i := 0; i < 1000; i++ {
		val, found := tree.Search(trie.Bytes(strconv.Itoa(i)))
		if found == false {
			t.Error("value didn't find in tree.")
		} else if (*val).(int) != i {
			t.Error("value doesn't match.'")
		}
	}
}

func TestTrieEmptyAdd(t *testing.T) {
	tree := trie.NewTrie()
	tree.Insert("", 0)
	if tree.Size() != 0 {
		t.Error("Trie added empty expression")
	}
}

func TestTrieGetAllKeys(t *testing.T) {
	tree := trie.NewTrie()
	for i := 0; i < 1000; i++ {
		tree.Insert(strconv.Itoa(i), i)
	}
	res := tree.GetAllKeys()
	for i := 0; i < 1000; i++ {
		found := false
		for j := 0; j < 1000; j++ {
			if bytes.Equal(res[j], trie.Bytes(strconv.Itoa(i))) {
				found = true
			}
		}
		if !found {
			t.Errorf("key %d not found in GetAllKeys", i)
		}
	}
}

func TestTrieDelete(t *testing.T) {
	tree := trie.NewTrie()
	for i := 0; i < 1000; i++ {
		tree.Insert(strconv.Itoa(i), i)
	}
	for i := 0; i < 1000; i++ {
		tree.Delete(strconv.Itoa(i))
		_, found := tree.Search(trie.Bytes(strconv.Itoa(i)))
		if found == true {
			t.Error("deleted key found in trie.")
		}
	}
	if tree.Size() != 0 {
		t.Errorf("trie size isn't zero after deleting all nodes.")
	}
}

func TestTrieGetKeyPrefix(t *testing.T) {
	tree := trie.NewTrie()
	for i := 31; i > 0; i-- {
		tree.Insert(fmt.Sprintf("%b", 1<<i), i)
		res := tree.GetPrefixKeys(trie.Bytes(fmt.Sprintf("%b", 1<<i)))
		if len(res) != 32-i {
			t.Error("prefix size is wrong")
		}
	}
}
func TestTrieGetValuePrefix(t *testing.T) {
	tree := trie.NewTrie()
	for i := 31; i > 0; i-- {
		tree.Insert(fmt.Sprintf("%b", 1<<i), i)
		res := tree.GetPrefixValues(trie.Bytes(fmt.Sprintf("%b", 1<<i)))
		if len(res) != 32-i {
			t.Error("prefix size is wrong")
		}
	}
}

package test

import "testing"
import "github.com/matinhimself/trie/pkg/trie"


func TestTrieAdd(t *testing.T) {
	tree := trie.NewTrie()
	tree.Insert(trie.Bytes("1"), 10)

}
package trie

import (
	"sync"
)

// Bytes reflects a type alias for a byte slice
type Bytes []byte

// TrieNode implements a node that the Trie is composed of. Each node contains
// a symbol that a key can be composed of unless the node is the root. The node
// has a collection of children that is represented as a hashmap, although,
// traditionally an array is used to represent each symbol in the given
// alphabet. The node may also contain a Value that indicates a possible query
// result.
//
// TODO: Handle the case where the Value given is a dummy Value which can be
// nil. Perhaps it's best to not store values at all.
type TrieNode struct {
	children []*TrieNode
	symbol   byte
	Value    interface{}
	root     bool
}

// Trie implements a thread-safe search tree that stores byte key Value pairs
// and allows for efficient queries.
type Trie struct {
	rw   sync.RWMutex
	root *TrieNode
	size int
}

// NewTrie returns a new initialized empty Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{root: true, children: make([]*TrieNode, 10)},
		size: 1,
	}
}

func newNode(symbol byte) *TrieNode {
	return &TrieNode{children: make([]*TrieNode, 10), symbol: symbol}
}

// Size returns the total number of nodes in the trie. The size includes the
// root node.
func (t *Trie) Size() int {
	t.rw.RLock()
	defer t.rw.RUnlock()
	return t.size
}

// Insert inserts a key Value pair into the trie. If the key already exists,
// the Value is updated. Insertion is performed by starting at the root
// and traversing the nodes all the way down until the key is exhausted. Once
// exhausted, the currNode pointer should be a pointer to the last symbol in
// the key and reflect the terminating node for that key Value pair.
func (t *Trie) Insert(key Bytes, value interface{}) {
	t.rw.Lock()
	defer t.rw.Unlock()

	currNode := t.root

	for _, sym := range key {
		symbol := sym - byte('0')
		if currNode.children[symbol] == nil {
			currNode.children[symbol] = newNode(symbol)
		}

		currNode = currNode.children[symbol]
	}

	// Only increment size if the key Value pair is new, otherwise we consider
	// the operation as an update.
	if currNode.Value == nil {
		t.size++
	}

	currNode.Value = value
}

// Search attempts to search for a Value in the trie given a key. If such a key
// exists, it's Value is returned along with a boolean to reflect that the key
// exists. Otherwise, an empty Value and false is returned.
func (t *Trie) Search(key Bytes) (*TrieNode, bool) {
	t.rw.RLock()
	defer t.rw.RUnlock()

	currNode := t.root

	for _, symbol := range key {
		if currNode.children[symbol-byte('0')] == nil {
			return nil, false
		}

		currNode = currNode.children[symbol-byte('0')]
	}

	return currNode, true
}

// GetAllKeys returns all the keys that exist in the trie. Keys are retrieved
// by performing a DFS on the trie where at each node we keep track of the
// current path (key) traversed thusfar and if that node has a Value. If so,
// the full path (key) is appended to a list. After the trie search is
// exhausted, the final list is returned.
func (t *Trie) GetAllKeys() []Bytes {
	visited := make(map[*TrieNode]bool)
	var keys []Bytes

	var dfsGetKeys func(n *TrieNode, key Bytes)
	dfsGetKeys = func(n *TrieNode, key Bytes) {
		if n != nil {
			pathKey := append(key, n.symbol)
			visited[n] = true

			if n.Value != nil {
				fullKey := make(Bytes, len(pathKey))

				// Copy the contents of the current path (key) to a new key so
				// future recursive calls will contain the correct bytes.
				copy(fullKey, pathKey)

				// Append the path (key) to the key list ignoring the first
				// byte which is the root symbol.
				keys = append(keys, fullKey[1:])
			}

			for _, child := range n.children {
				if _, ok := visited[child]; !ok {
					dfsGetKeys(child, pathKey)
				}
			}
		}
	}

	dfsGetKeys(t.root, Bytes{})
	return keys
}

// GetPrefixKeys returns all the keys that exist in the trie such that each key
// contains a specified prefix. Keys are retrieved by performing a DFS on the
// trie where at each node we keep track of the current path (key) and prefix
// traversed thusfar. If a node has a Value the full path (key) is appended to
// a list. After the trie search is exhausted, the final list is returned.
func (t *Trie) GetPrefixKeys(prefix Bytes) []Bytes {
	visited := make(map[*TrieNode]bool)
	var keys []Bytes

	if len(prefix) == 0 {
		return keys
	}

	var dfsGetPrefixKeys func(n *TrieNode, prefixIdx int, key Bytes)
	dfsGetPrefixKeys = func(n *TrieNode, prefixIdx int, key Bytes) {
		if n != nil {
			pathKey := append(key, n.symbol)

			if prefixIdx == len(prefix) || n.symbol == (prefix[prefixIdx]-byte('0')) {
				visited[n] = true

				if n.Value != nil {
					fullKey := make(Bytes, len(pathKey))

					// Copy the contents of the current path (key) to a new key
					// so future recursive calls will contain the correct
					// bytes.
					copy(fullKey, pathKey)
					keys = append(keys, fullKey)
				}

				if prefixIdx < len(prefix) {
					prefixIdx++
				}
				for _, child := range n.children {
					if child != nil {
						if _, ok := visited[child]; !ok {
							dfsGetPrefixKeys(child, prefixIdx, pathKey)
						}
					}
				}
			}
		}
	}

	// Find starting node from the root's children
	if n := t.root.children[prefix[0]-byte('0')]; n != nil {
		dfsGetPrefixKeys(n, 0, Bytes{})
	}

	return keys
}

// GetPrefixValues returns all the values that exist in the trie such that each
// key that corresponds to that Value contains a specified prefix. Values are
// retrieved by performing a DFS on the trie where at each node we check if the
// prefix is exhausted or matches thusfar and the current node has a Value. If
// the current node has a Value, it is appended to a list. After the trie
// search is exhausted, the final list is returned.
func (t *Trie) GetPrefixValues(prefix Bytes) []interface{} {
	visited := make(map[*TrieNode]bool)
	var values []interface{}

	if len(prefix) == 0 {
		return values
	}

	var dfsGetPrefixValues func(n *TrieNode, prefixIdx int)
	dfsGetPrefixValues = func(n *TrieNode, prefixIdx int) {
		if n != nil {
			if prefixIdx == len(prefix) || n.symbol == prefix[prefixIdx] {
				visited[n] = true

				if n.Value != nil {
					values = append(values, n.Value)
				}

				if prefixIdx < len(prefix) {
					prefixIdx++
				}

				for _, child := range n.children {
					if _, ok := visited[child]; !ok {
						dfsGetPrefixValues(child, prefixIdx)
					}
				}
			}
		}
	}

	//// Find starting node from the root's children
	if n := t.root.children[prefix[0]-byte('0')]; n != nil {
		dfsGetPrefixValues(n, 0)
	}

	return values
}

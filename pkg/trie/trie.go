package trie

import (
	"bytes"
	"sync"
)

// Bytes reflects a type alias for a byte slice
type Bytes []byte

const zeroAscii = byte('0')

// Node implements a node that the Trie is composed of. Each node contains
// a symbol.
type Node struct {
	parent   *Node
	children []*Node
	symbol   byte
	Value    interface{}
	root     bool
}

// Trie implements a thread-safe search tree that stores key Value pairs.
type Trie struct {
	rw   sync.RWMutex
	root *Node
	size int
}

// NewTrie returns a new initialized empty Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &Node{root: true, children: make([]*Node, 10)},
		size: 0,
	}
}

func newNode(symbol byte, parent *Node) *Node {
	return &Node{children: make([]*Node, 10), symbol: symbol, parent: parent}
}

// Size returns the total number of nodes in the trie.
func (t *Trie) Size() int {
	t.rw.RLock()
	defer t.rw.RUnlock()
	return t.size
}

func convert(s string) Bytes {
	var bs Bytes
	for i := 0; i < len(s); i++ {
		bs = append(bs, s[i] - byte('0'))
	}
	return bs
}

// Insert inserts a key Value pair into the trie. If the key already exists,
// the Value is updated.
func (t *Trie) Insert(sKey string, value interface{}) {
	key := convert(sKey)
	t.rw.Lock()
	defer t.rw.Unlock()
	if bytes.Equal(key, Bytes("")) {
		return
	}

	currNode := t.root

	for _, sym := range key {
		symbol := sym
		if currNode.children[symbol] == nil {
			currNode.children[symbol] = newNode(symbol, currNode)
		}

		currNode = currNode.children[symbol]
	}

	// Only increase size if the key Value pair is new, otherwise we consider
	// the operation as an update.
	if currNode.Value == nil {
		t.size++
	}

	currNode.Value = value
}

func (t *Trie) Delete(sKey string) (value *interface{}, deleted bool) {
	key := convert(sKey)
	t.rw.RLock()
	defer t.rw.RUnlock()

	currNode := t.root

	for _, symbol := range key {
		if currNode.children[symbol] == nil {
			return nil, false
		}
		currNode = currNode.children[symbol]
	}

	if currNode.Value != nil {
		t.size--
	}

	pTmpValue := currNode.Value
	parent := currNode.parent
	currNode.Value = nil

	for !hasChildren(parent.children){
		if parent.root {
			break
		}
		tmpPar := parent.parent
		parent = &Node{}
		parent = tmpPar
	}

	return &pTmpValue, true
}

func hasChildren(nodes []*Node) bool {
	hasChildren := false
	for _, node := range nodes {
		if node != nil && node.Value != nil{
			hasChildren = true
			break
		}
	}
	return hasChildren
}

// Search attempts to search for a Value in the trie given a key.
func (t *Trie) Search(sKey string) (*interface{}, bool) {
	key := convert(sKey)
	t.rw.RLock()
	defer t.rw.RUnlock()

	currNode := t.root

	for _, symbol := range key {
		if currNode.children[symbol] == nil {
			return nil, false
		}

		currNode = currNode.children[symbol]
	}
	if currNode.Value == nil {
		return nil, false
	}

	return &currNode.Value, true
}

// GetAllKeys returns all the keys that exist in the trie. Keys are retrieved
// by performing a DFS on the trie.
func (t *Trie) GetAllKeys() []string {
	visited := make(map[*Node]bool)
	var keys []Bytes

	var dfsGetKeys func(n *Node, key Bytes)
	dfsGetKeys = func(n *Node, key Bytes) {
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
	var strS []string
	for _, key := range keys {
		tmp := ""
		for _, b := range key {
			tmp += string(b + zeroAscii)
		}
		strS = append(strS, tmp)
	}

	return strS
}

// GetPrefixKeys returns all the keys that exist in the trie  Keys are retrieved
// by performing a DFS on the trie.
func (t *Trie) GetPrefixKeys(sPrefix string) []string {
	prefix := convert(sPrefix)
	visited := make(map[*Node]bool)
	var keys []Bytes

	if len(prefix) == 0 {
		return []string{}
	}

	var dfsGetPrefixKeys func(n *Node, prefixIdx int, key Bytes)
	dfsGetPrefixKeys = func(n *Node, prefixIdx int, key Bytes) {
		if n != nil {
			pathKey := append(key, n.symbol)

			if prefixIdx == len(prefix) || n.symbol == (prefix[prefixIdx]) {
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
	if n := t.root.children[prefix[0]]; n != nil {
		dfsGetPrefixKeys(n, 0, Bytes{})
	}
	var strS []string
	for _, key := range keys {
		if len(key) >= len(prefix) {
			tmp := ""
			for _, b := range key {
				tmp += string(b + zeroAscii)
			}
			strS = append(strS, tmp)
		}
	}
	return strS
}

// GetPrefixValues returns all the values that exist in the trie with given prefix
// Values retrieved by performing a DFS on the trie.
func (t *Trie) GetPrefixValues(sPrefix string) []interface{} {
	prefix := convert(sPrefix)
	visited := make(map[*Node]bool)
	var values []interface{}

	if len(prefix) == 0 {
		return values
	}

	var dfsGetPrefixValues func(n *Node, prefixIdx int)
	dfsGetPrefixValues = func(n *Node, prefixIdx int) {
		if n != nil {
			if prefixIdx == len(prefix) || n.symbol == (prefix[prefixIdx]) {
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
	if n := t.root.children[prefix[0]]; n != nil {
		dfsGetPrefixValues(n, 0)
	}

	return values
}

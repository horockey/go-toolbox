package b_tree

import (
	"github.com/horockey/go-toolbox/datastructs/pkg/comparer"
	"github.com/horockey/go-toolbox/datastructs/trees"
)

type bTree[K, V any] struct {
	comparer comparer.Comparer[K]

	root *node[K, V]
	size int
}

// Creates new b-tree with string key type.
func New[V any]() *bTree[string, V] {
	return &bTree[string, V]{
		comparer: comparer.NewStringComparer(),
	}
}

// Creates new b-tree with custom key type.
// Corresponding Comparer required.
func NewWithCustomKeyType[K, V any](comp comparer.Comparer[K]) *bTree[K, V] {
	return &bTree[K, V]{
		comparer: comp,
	}
}

func (tree *bTree[K, V]) Get(key K) (V, error) {
	n, idx := tree.get(key)
	if n == nil {
		return *new(V), trees.NotFoundError[K]{GivenKey: key}
	}

	return n.values[idx], nil
}

func (tree *bTree[K, V]) Insert(key K, val V) error {
	// TODO
	return nil
}

func (tree *bTree[K, V]) Remove(key K) error {
	// TODO
	return nil
}

func (tree *bTree[K, V]) Keys() []K {
	// TODO
	return nil
}

func (tree *bTree[K, V]) Clear() {
	tree.root = nil
	tree.size = 0
}

func (tree *bTree[K, V]) get(key K) (n *node[K, V], idx int) {
	cur := tree.root
	for cur != nil {
		for keyIdx := 0; keyIdx < cur.keysCount; keyIdx++ {
			nodeKey := cur.keys[keyIdx]

			switch tree.comparer.Compare(key, nodeKey) {
			case 0:
				return cur, keyIdx
			case 1:
				cur = cur.children[keyIdx]
			case -1:
				if keyIdx == cur.keysCount-1 {
					cur = cur.children[keyIdx+1]
				}
			}
		}
	}
	return nil, 0
}

package red_black_tree

import (
	"errors"
	"sync"

	"github.com/horockey/go-toolbox/datastructs/pkg/comparer"
)

var ErrNotFound error = errors.New("unable to find element for given key")

type redBlackTree[K, V any] struct {
	mu sync.RWMutex

	comparer comparer.Comparer[K]

	root *node[K, V]
	size uint
}

// Creates new red-black tree with string key type.
func New[V any]() *redBlackTree[string, V] {
	return &redBlackTree[string, V]{
		comparer: comparer.NewStringComparer(),
	}
}

// Creates new red-black tree with custom key type.
// Corresponding Comparer required.
func NewWithCustomKey[K, V any](comp comparer.Comparer[K]) *redBlackTree[K, V] {
	return &redBlackTree[K, V]{
		comparer: comp,
	}
}

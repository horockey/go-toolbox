package avl_tree

import (
	"errors"
	"fmt"
	"sync"

	"github.com/horockey/go-toolbox/datastructs/internal/comparer"
)

var ErrNotFound error = errors.New("unable to find element for given key")

type avlTree[K any, V any] struct {
	mu sync.RWMutex

	comp comparer.Comparer[K]

	root *node[K, V]
	size uint
}

// Creates new AVL tree with string key type.
func New[V any]() *avlTree[string, V] {
	return &avlTree[string, V]{
		comp: &comparer.StringComparer{},
	}
}

// Creates new AVL tree with custom key type.
// Corresponding Comparer required.
func NewWithCustomKey[K any, V any](comp comparer.Comparer[K]) *avlTree[K, V] {
	return &avlTree[K, V]{
		comp: comp,
	}
}

// Adds new element to AVL tree.
// Fixes balance, if necessary.
func (t *avlTree[K, V]) Insert(key K, val V) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.insert(t.root, key, val); err != nil {
		return fmt.Errorf("inserting new element: %w", err)
	}

	t.size++
	return nil
}

// If tree contains given key, returns corresponding value.
// Returns ErrNotFounf otherwise.
func (t *avlTree[K, V]) Get(key K) (V, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	n := t.get(key)
	if n == nil {
		return *new(V), ErrNotFound
	}

	return n.Value, nil
}

func (t *avlTree[K, V]) insert(subroot *node[K, V], key K, val V) error {
	if subroot == nil {
		if subroot == t.root {
			t.root = &node[K, V]{
				Key:   key,
				Value: val,
			}
			return nil
		}
		return errors.New("given subroot is nil")
	}

	var nextNode *node[K, V]
	var nextNodeIsLeft bool
	switch t.comp.Compare(key, subroot.Key) {
	case -1, 0:
		nextNode = subroot.right
		nextNodeIsLeft = false
	case 1:
		nextNode = subroot.left
		nextNodeIsLeft = true
	}

	if nextNode != nil {
		if err := t.insert(nextNode, key, val); err != nil {
			return err
		}
	} else {
		nextNode = &node[K, V]{
			Key:    key,
			Value:  val,
			parent: subroot,
		}
		switch nextNodeIsLeft {
		case true:
			subroot.left = nextNode
		case false:
			subroot.right = nextNode
		}
	}

	t.balance(subroot)

	return nil
}

func (t *avlTree[K, V]) get(key K) *node[K, V] {
	cur := t.root
	for cur != nil {
		switch t.comp.Compare(key, cur.Key) {
		case -1:
			cur = cur.right
		case 1:
			cur = cur.left
		case 0:
			return cur
		}
	}
	return nil
}

func (t *avlTree[K, V]) balance(subroot *node[K, V]) {
	subroot.fixHeight()
	switch subroot.balanceFactor() {
	case 2:
		if subroot.right.balanceFactor() < 0 {
			t.rightRotation(subroot.right)
		}
		t.leftRotation(subroot)
	case -2:
		if subroot.left.balanceFactor() > 0 {
			t.leftRotation(subroot.left)
		}
		t.rightRotation(subroot)
	default:
		return
	}
}

func (t *avlTree[K, V]) leftRotation(subroot *node[K, V]) {
	temp := subroot.right
	if temp == nil {
		return
	}

	subroot.right = temp.left
	subroot.right.parent = subroot

	temp.left = subroot
	temp.parent = subroot.parent
	if subroot.parent != nil {
		switch subroot {
		case subroot.parent.left:
			subroot.parent.left = temp
		case subroot.parent.right:
			subroot.parent.right = temp
		}
	} else {
		t.root = temp
	}

	subroot.parent = temp

	subroot.fixHeight()
	temp.fixHeight()
}

func (t *avlTree[K, V]) rightRotation(subroot *node[K, V]) {
	temp := subroot.left
	if temp == nil {
		return
	}

	subroot.left = temp.right
	subroot.left.parent = subroot

	temp.right = subroot
	temp.parent = subroot.parent
	if subroot.parent != nil {
		switch subroot {
		case subroot.parent.left:
			subroot.parent.left = temp
		case subroot.parent.right:
			subroot.parent.right = temp
		}
	} else {
		t.root = temp
	}

	subroot.parent = temp

	subroot.fixHeight()
	temp.fixHeight()
}

package avl_tree

import (
	"errors"
	"fmt"
	"sync"

	"github.com/horockey/go-toolbox/datastructs/comparer"
	"github.com/horockey/go-toolbox/datastructs/trees"
)

type avlTree[K, V any] struct {
	mu sync.RWMutex

	comparer comparer.Comparer[K]

	root *node[K, V]
	size uint
}

// Creates new AVL tree with string key type.
func New[V any]() *avlTree[string, V] {
	return &avlTree[string, V]{
		comparer: comparer.NewStringComparer(),
	}
}

// Creates new AVL tree with custom key type.
// Corresponding Comparer required.
func NewWithCustomKey[K, V any](comp comparer.Comparer[K]) *avlTree[K, V] {
	return &avlTree[K, V]{
		comparer: comp,
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
// Returns ErrNotFound otherwise.
func (t *avlTree[K, V]) Get(key K) (V, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	n := t.get(key)
	if n == nil {
		return *new(V), trees.NotFoundError[K]{GivenKey: key}
	}

	return n.Value, nil
}

// If tree contains given key, removes KV pair from itself.
// Returns ErrNotFound otherwise.
func (t *avlTree[K, V]) Remove(key K) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	n := t.get(key)
	if n == nil {
		return trees.NotFoundError[K]{GivenKey: key}
	}
	t.remove(n)

	t.size--
	return nil
}

// Removes all keys from tree.
func (t *avlTree[K, V]) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.root = nil
	t.size = 0
}

// Returns all keys that the tree contains in ascending order.
// Order is defined by comparer.
func (t *avlTree[K, V]) Keys() []K {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.keys(t.root)
}

// Returns number of nodes in tree.
func (t *avlTree[K, V]) Size() uint {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.size
}

// Return the height of tree root.
func (t *avlTree[K, V]) Height() uint {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.root == nil {
		return 0
	}

	return t.root.height
}

func (t *avlTree[K, V]) insert(subroot *node[K, V], key K, val V) error {
	if subroot == nil {
		if subroot == t.root {
			t.root = newNode(key, val)
			return nil
		}
		return errors.New("given subroot is nil")
	}

	var nextNode *node[K, V]
	var nextNodeIsLeft bool
	switch t.comparer.Compare(key, subroot.Key) {
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
		nextNode = newNode(key, val)
		nextNode.parent = subroot
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

func (t *avlTree[K, V]) remove(subroot *node[K, V]) {
	if subroot == nil {
		return
	}

	if subroot.isLeaf() ||
		subroot.hasOnlyLeft() ||
		subroot.hasOnlyRight() {
		var child *node[K, V]
		if subroot.hasOnlyLeft() {
			child = subroot.left
		}
		if subroot.hasOnlyRight() {
			child = subroot.right
		}

		if subroot.parent == nil {
			t.root = child
			return
		}

		switch subroot {
		case subroot.parent.right:
			subroot.parent.right = child
		case subroot.parent.left:
			subroot.parent.left = child
		}

		if child != nil {
			child.parent = subroot.parent
		}
	} else {
		n := t.removeMin(subroot.right)
		// keep all pointers on their place, just susbtitute the payload
		subroot.Key = n.Key
		subroot.Value = n.Value
	}

	for cur := subroot.parent; cur != nil; cur = cur.parent {
		t.balance(cur)
	}
}

func (t *avlTree[K, V]) removeMin(subroot *node[K, V]) *node[K, V] {
	if subroot.isLeaf() || subroot.hasOnlyRight() {
		// is the most left node
		t.remove(subroot)
		return subroot
	}

	n := t.removeMin(subroot.left)
	t.balance(subroot)
	return n
}

func (t *avlTree[K, V]) get(key K) *node[K, V] {
	cur := t.root
	for cur != nil {
		switch t.comparer.Compare(key, cur.Key) {
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

func (t *avlTree[K, V]) keys(subroot *node[K, V]) []K {
	if subroot == nil {
		return []K{}
	}

	res := append(t.keys(subroot.left), subroot.Key)
	res = append(res, t.keys(subroot.right)...)
	return res
}

func (t *avlTree[K, V]) balance(subroot *node[K, V]) {
	subroot.fixHeight()
	switch subroot.balanceFactor() {
	case 2:
		if subroot.right.balanceFactor() < 0 {
			// turns to greater left rotation
			t.rightRotation(subroot.right)
		}
		t.leftRotation(subroot)
	case -2:
		if subroot.left.balanceFactor() > 0 {
			// turns to greater right rotation
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
	if subroot.right != nil {
		subroot.right.parent = subroot
	}

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
	if subroot.left != nil {
		subroot.left.parent = subroot
	}

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

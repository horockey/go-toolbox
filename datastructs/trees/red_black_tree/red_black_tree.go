package red_black_tree

import (
	"errors"
	"fmt"
	"sync"

	"github.com/horockey/go-toolbox/datastructs/pkg/comparer"
	"github.com/horockey/go-toolbox/datastructs/trees"
)

type redBlackTree[K, V any] struct {
	mu sync.RWMutex

	comparer comparer.Comparer[K]

	root *node[K, V]
	size int
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

func (tree *redBlackTree[K, V]) Get(key K) (V, error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	n, err := tree.getNode(key)
	if err != nil {
		return *new(V), fmt.Errorf("getting element from tree: %w", err)
	}

	return n.Value, nil
}

func (tree *redBlackTree[K, V]) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	return tree.size
}

func (tree *redBlackTree[K, V]) Insert(key K, val V) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	if err := tree.insert(tree.root, key, val); err != nil {
		return fmt.Errorf("inserting new element: %w", err)
	}

	tree.size++

	return nil
}

func (tree *redBlackTree[K, V]) Remove(key K) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	n, err := tree.getNode(key)
	if err != nil {
		return fmt.Errorf("searching element in tree: %w", err)
	}

	tree.remove(n)
	tree.size--

	return nil
}

func (tree *redBlackTree[K, V]) insert(subroot *node[K, V], key K, val V) error {
	if subroot == nil {
		if subroot == tree.root {
			tree.root = newNode(key, val)
			tree.root.color = ColorBlack
			return nil
		}
		return errors.New("given subroot is nil")
	}

	son := &node[K, V]{}
	var dir trees.Direction
	switch tree.comparer.Compare(key, subroot.Key) {
	case 0:
		return trees.AlreadyExistsError[K]{GivenKey: key}
	case 1:
		son = subroot.left
		dir = trees.DirectionLeft
	case -1:
		son = subroot.right
		dir = trees.DirectionRight
	}

	if son == nil {
		son = newNode(key, val)
		son.direction = dir
		son.parent = subroot
		switch dir {
		case trees.DirectionLeft:
			subroot.left = son
		case trees.DirectionRight:
			subroot.right = son
		}

		return nil
	}

	err := tree.insert(son, key, val)
	if err != nil {
		return err
	}

	tree.ballanceAfterInsertion(son)

	return nil
}

func (tree *redBlackTree[K, V]) remove(n *node[K, V]) {
	rmNodeColor := n.color
	child := &node[K, V]{}

	defer func() {
		if rmNodeColor == ColorBlack {
			tree.ballanceTreeAfterRemoval(child)
		}
	}()

	transplantNode := func(from, to *node[K, V]) {
		if to == tree.root {
			tree.root = from
			return
		}
		switch n.direction {
		case trees.DirectionLeft:
			to.parent.left = from
		case trees.DirectionRight:
			to.parent.right = from
		}
		from.parent = to.parent
	}

	if n.hasOneChild() || n.hasNoChildren() {

		switch {
		case n.hasOnlyLeft():
			child = n.left
		default:
			child = n.right
		}

		transplantNode(n, child)
		return
	}

	minNode := tree.getLeftmost(n.right)
	n.Key = minNode.Key
	n.Value = minNode.Value
	rmNodeColor = minNode.color

	switch {
	case n.hasOnlyLeft():
		child = minNode.left
	default:
		child = minNode.right
	}
	transplantNode(minNode, child)
}

func (tree *redBlackTree[K, V]) ballanceAfterInsertion(newNode *node[K, V]) {
	if newNode == tree.root {
		newNode.color = ColorBlack
		return
	}

	if newNode.parent.color == ColorBlack {
		// no need to ballance
		return
	}

	if newNode.uncle() != nil && newNode.uncle().color == ColorRed {
		newNode.parent.color = ColorBlack
		newNode.uncle().color = ColorBlack
		if newNode.grandParent() != tree.root {
			newNode.grandParent().color = ColorRed
		}
		return
	}

	// zigzag to line
	if newNode.direction != newNode.parent.direction {
		temp := newNode.parent
		if newNode.direction == trees.DirectionRight {
			tree.leftRotation(newNode.parent)
		} else {
			tree.rightRotation(newNode.parent)
		}
		newNode = temp
	}

	newNode.parent.color = ColorBlack
	switch newNode.parent.direction {
	case trees.DirectionLeft:
		tree.rightRotation(newNode.grandParent())
	case trees.DirectionRight:
		tree.leftRotation(newNode.parent)
	}
}

func (tree *redBlackTree[K, V]) ballanceTreeAfterRemoval(rmNode *node[K, V]) {
	for rmNode != tree.root && rmNode.color == ColorBlack {
		brother := rmNode.brother()
		switch rmNode.direction {
		case trees.DirectionLeft:
			if brother.color == ColorRed {
				brother.color = ColorBlack
				rmNode.parent.color = ColorRed
				tree.leftRotation(rmNode.parent)
				brother = rmNode.parent.right
			}
			if brother.left.color == ColorBlack && brother.right.color == ColorBlack {
				brother.color = ColorRed
				rmNode = rmNode.parent
			} else {
				if brother.right.color == ColorBlack {
					brother.left.color = ColorBlack
					brother.color = ColorRed
					tree.rightRotation(brother)
					brother = rmNode.parent.right
				}
				brother.color = rmNode.parent.color
				rmNode.parent.color = ColorBlack
				rmNode.right.color = ColorBlack
				tree.leftRotation(rmNode.parent)
				rmNode = tree.root
			}
		case trees.DirectionRight:
			if brother.color == ColorRed {
				brother.color = ColorBlack
				rmNode.parent.color = ColorRed
				tree.rightRotation(rmNode.parent)
				brother = rmNode.parent.left
			}
			if brother.left.color == ColorBlack && brother.right.color == ColorBlack {
				brother.color = ColorRed
				rmNode = rmNode.parent
			} else {
				if brother.left.color == ColorBlack {
					brother.right.color = ColorBlack
					brother.color = ColorRed
					tree.leftRotation(brother)
					brother = rmNode.parent.left
				}
				brother.color = rmNode.parent.color
				rmNode.parent.color = ColorBlack
				rmNode.left.color = ColorBlack
				tree.rightRotation(rmNode.parent)
				rmNode = tree.root
			}
		}
	}

	rmNode.color = ColorBlack
}

func (t *redBlackTree[K, V]) getLeftmost(subroot *node[K, V]) *node[K, V] {
	if subroot.hasNoChildren() || subroot.hasOnlyRight() {
		// is the most left node
		return subroot
	}

	n := t.getLeftmost(subroot.left)
	return n
}

func (tree *redBlackTree[K, V]) getNode(key K) (*node[K, V], error) {
	cur := tree.root
	for cur != nil {
		switch tree.comparer.Compare(key, cur.Key) {
		case 0:
			return cur, nil
		case 1:
			cur = cur.left
		case -1:
			cur = cur.right
		}
	}
	return nil, trees.NotFoundError[K]{GivenKey: key}
}

func (t *redBlackTree[K, V]) leftRotation(subroot *node[K, V]) {
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
}

func (t *redBlackTree[K, V]) rightRotation(subroot *node[K, V]) {
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
}

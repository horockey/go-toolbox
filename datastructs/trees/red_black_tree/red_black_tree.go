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

func (tree *redBlackTree[K, V]) Clear() {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	tree.root = nil
	tree.size = 0
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

		tree.ballanceAfterInsertion(son)

		return nil
	}

	err := tree.insert(son, key, val)
	if err != nil {
		return err
	}

	return nil
}

func (tree *redBlackTree[K, V]) remove(n *node[K, V]) {
	rmNode := n

	defer func() {
		if rmNode.color == ColorBlack {
			tree.ballanceAfterRemoval(rmNode)
		}
	}()

	transplantNode := func(from, to *node[K, V]) {
		if to == nil {
			return
		}

		if to == tree.root {
			tree.root = from
			if from != nil {
				from.parent = nil
				from.direction = trees.DirectionNoDir
			}
			return
		}

		switch to.direction {
		case trees.DirectionLeft:
			to.parent.left = from
		case trees.DirectionRight:
			to.parent.right = from
		}
		if from != nil {
			from.parent = to.parent
			from.direction = to.direction
		}
	}

	if n.hasOneChild() || n.hasNoChildren() {
		var child *node[K, V]
		switch {
		case n.hasOnlyLeft():
			child = n.left
		default:
			child = n.right
		}

		transplantNode(child, n)
		return
	}

	minNode := tree.getLeftmost(n.right)
	n.Key = minNode.Key
	n.Value = minNode.Value
	rmNode = minNode

	var child *node[K, V]
	switch {
	case minNode.hasOnlyLeft():
		child = minNode.left
	default:
		child = minNode.right
	}
	transplantNode(child, minNode)
}

func (tree *redBlackTree[K, V]) ballanceAfterInsertion(newNode *node[K, V]) {
	defer func() {
		tree.root.color = ColorBlack
	}()

	if newNode == tree.root {
		return
	}

	for newNode != tree.root && newNode.parent.color == ColorRed {
		if newNode.uncle() != nil && newNode.uncle().color == ColorRed {
			newNode.parent.color = ColorBlack
			newNode.uncle().color = ColorBlack
			newNode.grandParent().color = ColorRed
			newNode = newNode.grandParent()
			continue
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
		newNode.grandParent().color = ColorRed
		switch newNode.parent.direction {
		case trees.DirectionLeft:
			tree.rightRotation(newNode.grandParent())
		case trees.DirectionRight:
			tree.leftRotation(newNode.grandParent())
		}
	}
}

func (tree *redBlackTree[K, V]) ballanceAfterRemoval(rmNode *node[K, V]) {
	for rmNode != tree.root && rmNode.color == ColorBlack {
		brother := rmNode.brother()
		switch rmNode.direction {
		default:
			rmNode = tree.root
			break
		case trees.DirectionLeft:
			if brother.color == ColorRed {
				brother.color = ColorBlack
				rmNode.parent.color = ColorRed
				tree.leftRotation(rmNode.parent)
				brother = rmNode.parent.right
			}
			if (brother.left == nil || brother.left.color == ColorBlack) &&
				(brother.right == nil || brother.right.color == ColorBlack) {
				brother.color = ColorRed
				rmNode = rmNode.parent
				continue
			}
			if brother.right == nil || brother.right.color == ColorBlack {
				if brother.left != nil {
					brother.left.color = ColorBlack
				}
				brother.color = ColorRed
				tree.rightRotation(brother)
				brother = rmNode.parent.right
			}
			brother.color = rmNode.parent.color
			rmNode.parent.color = ColorBlack
			if brother.right != nil {
				brother.right.color = ColorBlack
			}
			tree.leftRotation(rmNode.parent)
			rmNode = tree.root

		case trees.DirectionRight:
			if brother.color == ColorRed {
				brother.color = ColorBlack
				rmNode.parent.color = ColorRed
				tree.rightRotation(rmNode.parent)
				brother = rmNode.parent.left
			}
			if (brother.left == nil || brother.left.color == ColorBlack) &&
				(brother.right == nil || brother.right.color == ColorBlack) {
				brother.color = ColorRed
				rmNode = rmNode.parent
				continue
			}
			if brother.left == nil || brother.left.color == ColorBlack {
				if brother.right != nil {
					brother.right.color = ColorBlack
				}
				brother.color = ColorRed
				tree.leftRotation(brother)
				brother = rmNode.parent.left
			}
			brother.color = rmNode.parent.color
			rmNode.parent.color = ColorBlack
			if brother.left != nil {
				brother.left.color = ColorBlack
			}
			tree.rightRotation(rmNode.parent)
			rmNode = tree.root
		}

	}

	if rmNode != nil {
		rmNode.color = ColorBlack
	}
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
		subroot.right.direction = trees.DirectionRight
	}

	temp.left = subroot
	temp.parent = subroot.parent
	if subroot.parent != nil {
		switch subroot.direction {
		case trees.DirectionLeft:
			subroot.parent.left = temp
		case trees.DirectionRight:
			subroot.parent.right = temp
		}

		temp.direction = subroot.direction
	} else {
		t.root = temp
		temp.direction = trees.DirectionNoDir
	}

	subroot.direction = trees.DirectionLeft
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
		subroot.left.direction = trees.DirectionLeft
	}

	temp.right = subroot
	temp.parent = subroot.parent
	if subroot.parent != nil {
		switch subroot.direction {
		case trees.DirectionLeft:
			subroot.parent.left = temp
		case trees.DirectionRight:
			subroot.parent.right = temp
		}

		temp.direction = subroot.direction
	} else {
		t.root = temp
		temp.direction = trees.DirectionNoDir
	}

	subroot.direction = trees.DirectionRight
	subroot.parent = temp
}

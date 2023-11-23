package avl_tree

import (
	"github.com/horockey/go-toolbox/math"
)

type node[K, V any] struct {
	Key   K
	Value V

	left   *node[K, V]
	right  *node[K, V]
	parent *node[K, V]

	height uint
}

func newNode[K, V any](key K, val V) *node[K, V] {
	return &node[K, V]{
		Key:    key,
		Value:  val,
		height: 1,
	}
}

func (n *node[K, V]) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *node[K, V]) hasOnlyLeft() bool {
	return n.left != nil && n.right == nil
}

func (n *node[K, V]) hasOnlyRight() bool {
	return n.left == nil && n.right != nil
}

func (n *node[K, V]) balanceFactor() int {
	l, r := 0, 0
	if n.left != nil {
		l = int(n.left.height)
	}
	if n.right != nil {
		r = int(n.right.height)
	}
	return r - l
}

func (n *node[K, V]) fixHeight() {
	if n.isLeaf() {
		n.height = 1
		return
	}

	var l, r uint = 0, 0
	if n.left != nil {
		l = n.left.height
	}
	if n.right != nil {
		r = n.right.height
	}

	n.height = math.Max(l, r) + 1
}

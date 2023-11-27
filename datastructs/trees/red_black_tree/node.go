package red_black_tree

import "github.com/horockey/go-toolbox/datastructs/trees"

type node[K, V any] struct {
	Key   K
	Value V

	color color

	left      *node[K, V]
	right     *node[K, V]
	parent    *node[K, V]
	direction trees.Direction
}

func newNode[K, V any](key K, val V) *node[K, V] {
	return &node[K, V]{
		Key:   key,
		Value: val,
		color: ColorRed,
	}
}

func (n *node[K, V]) grandParent() *node[K, V] {
	if n.parent == nil {
		return nil
	}
	return n.parent.parent
}

func (n *node[K, V]) brother() *node[K, V] {
	if n.parent == nil {
		return nil
	}

	switch n.direction {
	case trees.DirectionLeft:
		return n.parent.right
	case trees.DirectionRight:
		return n.parent.left
	}

	return nil
}

func (n *node[K, V]) uncle() *node[K, V] {
	grand := n.grandParent()
	if grand == nil {
		return nil
	}

	switch n.parent.direction {
	case trees.DirectionLeft:
		return grand.right
	case trees.DirectionRight:
		return grand.left
	}

	return nil
}

func (n *node[K, V]) hasNoChildren() bool {
	return n.left == nil && n.right == nil
}

func (n *node[K, V]) hasOneChild() bool {
	nodeCount := 0
	if n.left != nil {
		nodeCount++
	}
	if n.right != nil {
		nodeCount++
	}
	return nodeCount == 1
}

func (n *node[K, V]) hasOnlyLeft() bool {
	return n.left != nil && n.right == nil
}

func (n *node[K, V]) hasOnlyRight() bool {
	return n.left == nil && n.right != nil
}

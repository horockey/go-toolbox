package red_black_tree

type Node[K, V any] struct {
	Key   K
	Value V
}

type node[K, V any] struct {
	Key   K
	Value V

	color  color
	left   *node[K, V]
	right  *node[K, V]
	parent *node[K, V]
}

func (n *node[K, V]) ToPublic() *Node[K, V] {
	return &Node[K, V]{
		Key:   n.Key,
		Value: n.Value,
	}
}

func newNode[K, V any](key K, val V) *node[K, V] {
	return &node[K, V]{
		Key:   key,
		Value: val,
	}
}

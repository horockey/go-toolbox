package b_tree

type node[K, V any] struct {
	keys      []K
	values    []V
	keysCount int

	children   []*node[K, V]
	childCount int
}

func newNode[K, V any](kvs ...kv[K, V]) *node[K, V] {
	n := node[K, V]{
		keys:     make([]K, len(kvs)),
		values:   make([]V, len(kvs)),
		children: make([]*node[K, V], len(kvs)-1),
	}

	for _, kv := range kvs {
		n.keys = append(n.keys, kv.key)
		n.values = append(n.values, kv.value)
		n.keysCount++
	}

	return &n
}

func (n *node[K, V]) isLeaf() bool {
	isLeaf := true
	for idx := 0; idx < n.childCount && isLeaf; idx++ {
		isLeaf = isLeaf && n.children[idx] == nil
	}

	return isLeaf
}

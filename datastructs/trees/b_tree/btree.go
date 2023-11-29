package btree

import (
	"cmp"
	"fmt"
	"io"
	"strings"
	"sync"
)

const (
	DefaultFreeListSize = 32
)

type BTree BTreeG[Item]

func New(degree int) *BTree {
	return (*BTree)(NewG[Item](degree, itemLess))
}

type FreeListG[T any] struct {
	mu       sync.Mutex
	freelist []*node[T]
}

func NewFreeListG[T any](size int) *FreeListG[T] {
	return &FreeListG[T]{freelist: make([]*node[T], 0, size)}
}

func (f *FreeListG[T]) freeNode(n *node[T]) (out bool) {
	f.mu.Lock()
	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	f.mu.Unlock()
	return
}

type ItemIteratorG[T any] func(item T) bool

func Less[T cmp.Ordered]() LessFunc[T] {
	return func(a, b T) bool { return a < b }
}

func NewOrderedG[T cmp.Ordered](degree int) *BTreeG[T] {
	return NewG[T](degree, Less[T]())
}

func NewG[T any](degree int, less LessFunc[T]) *BTreeG[T] {
	return NewWithFreeListG(degree, less, NewFreeListG[T](DefaultFreeListSize))
}

func NewWithFreeListG[T any](degree int, less LessFunc[T], f *FreeListG[T]) *BTreeG[T] {
	if degree <= 1 {
		panic("bad degree")
	}
	return &BTreeG[T]{
		degree: degree,
		cow:    &copyOnWriteContext[T]{freelist: f, less: less},
	}
}

func min[T any](n *node[T]) (_ T, found bool) {
	if n == nil {
		return
	}
	for len(n.children) > 0 {
		n = n.children[0]
	}
	if len(n.items) == 0 {
		return
	}
	return n.items[0], true
}

func max[T any](n *node[T]) (_ T, found bool) {
	if n == nil {
		return
	}
	for len(n.children) > 0 {
		n = n.children[len(n.children)-1]
	}
	if len(n.items) == 0 {
		return
	}
	return n.items[len(n.items)-1], true
}

type toRemove int

const (
	removeItem toRemove = iota
	removeMin
	removeMax
)

func (n *node[T]) remove(item T, minItems int, typ toRemove) (_ T, _ bool) {
	var i int
	var found bool
	switch typ {
	case removeMax:
		if len(n.children) == 0 {
			return n.items.pop(), true
		}
		i = len(n.items)
	case removeMin:
		if len(n.children) == 0 {
			return n.items.removeAt(0), true
		}
		i = 0
	case removeItem:
		i, found = n.items.find(item, n.cow.less)
		if len(n.children) == 0 {
			if found {
				return n.items.removeAt(i), true
			}
			return
		}
	default:
		panic("invalid type")
	}
	if len(n.children[i].items) <= minItems {
		return n.growChildAndRemove(i, item, minItems, typ)
	}
	child := n.mutableChild(i)
	if found {
		out := n.items[i]
		var zero T
		n.items[i], _ = child.remove(zero, minItems, removeMax)
		return out, true
	}
	return child.remove(item, minItems, typ)
}

func (n *node[T]) growChildAndRemove(i int, item T, minItems int, typ toRemove) (T, bool) {
	if i > 0 && len(n.children[i-1].items) > minItems {
		child := n.mutableChild(i)
		stealFrom := n.mutableChild(i - 1)
		stolenItem := stealFrom.items.pop()
		child.items.insertAt(0, n.items[i-1])
		n.items[i-1] = stolenItem
		if len(stealFrom.children) > 0 {
			child.children.insertAt(0, stealFrom.children.pop())
		}
	} else if i < len(n.items) && len(n.children[i+1].items) > minItems {
		child := n.mutableChild(i)
		stealFrom := n.mutableChild(i + 1)
		stolenItem := stealFrom.items.removeAt(0)
		child.items = append(child.items, n.items[i])
		n.items[i] = stolenItem
		if len(stealFrom.children) > 0 {
			child.children = append(child.children, stealFrom.children.removeAt(0))
		}
	} else {
		if i >= len(n.items) {
			i--
		}
		child := n.mutableChild(i)
		mergeItem := n.items.removeAt(i)
		mergeChild := n.children.removeAt(i + 1)
		child.items = append(child.items, mergeItem)
		child.items = append(child.items, mergeChild.items...)
		child.children = append(child.children, mergeChild.children...)
		n.cow.freeNode(mergeChild)
	}
	return n.remove(item, minItems, typ)
}

type direction int

const (
	descend = direction(-1)
	ascend  = direction(+1)
)

type optionalItem[T any] struct {
	item  T
	valid bool
}

func optional[T any](item T) optionalItem[T] {
	return optionalItem[T]{item: item, valid: true}
}

func empty[T any]() optionalItem[T] {
	return optionalItem[T]{}
}

func (n *node[T]) iterate(dir direction, start, stop optionalItem[T], includeStart bool, hit bool, iter ItemIteratorG[T]) (bool, bool) {
	var ok, found bool
	var index int
	switch dir {
	case ascend:
		if start.valid {
			index, _ = n.items.find(start.item, n.cow.less)
		}
		for i := index; i < len(n.items); i++ {
			if len(n.children) > 0 {
				if hit, ok = n.children[i].iterate(dir, start, stop, includeStart, hit, iter); !ok {
					return hit, false
				}
			}
			if !includeStart && !hit && start.valid && !n.cow.less(start.item, n.items[i]) {
				hit = true
				continue
			}
			hit = true
			if stop.valid && !n.cow.less(n.items[i], stop.item) {
				return hit, false
			}
			if !iter(n.items[i]) {
				return hit, false
			}
		}
		if len(n.children) > 0 {
			if hit, ok = n.children[len(n.children)-1].iterate(dir, start, stop, includeStart, hit, iter); !ok {
				return hit, false
			}
		}
	case descend:
		if start.valid {
			index, found = n.items.find(start.item, n.cow.less)
			if !found {
				index = index - 1
			}
		} else {
			index = len(n.items) - 1
		}
		for i := index; i >= 0; i-- {
			if start.valid && !n.cow.less(n.items[i], start.item) {
				if !includeStart || hit || n.cow.less(start.item, n.items[i]) {
					continue
				}
			}
			if len(n.children) > 0 {
				if hit, ok = n.children[i+1].iterate(dir, start, stop, includeStart, hit, iter); !ok {
					return hit, false
				}
			}
			if stop.valid && !n.cow.less(stop.item, n.items[i]) {
				return hit, false //	continue
			}
			hit = true
			if !iter(n.items[i]) {
				return hit, false
			}
		}
		if len(n.children) > 0 {
			if hit, ok = n.children[0].iterate(dir, start, stop, includeStart, hit, iter); !ok {
				return hit, false
			}
		}
	}
	return hit, true
}

func (n *node[T]) print(w io.Writer, level int) {
	fmt.Fprintf(w, "%sNODE:%v\n", strings.Repeat("  ", level), n.items)
	for _, c := range n.children {
		c.print(w, level+1)
	}
}

type BTreeG[T any] struct {
	degree int
	length int
	root   *node[T]
	cow    *copyOnWriteContext[T]
}

type LessFunc[T any] func(a, b T) bool

type copyOnWriteContext[T any] struct {
	freelist *FreeListG[T]
	less     LessFunc[T]
}

func (t *BTreeG[T]) Clone() (t2 *BTreeG[T]) {
	cow1, cow2 := *t.cow, *t.cow
	out := *t
	t.cow = &cow1
	out.cow = &cow2
	return &out
}

func (t *BTreeG[T]) maxItems() int {
	return t.degree*2 - 1
}

func (t *BTreeG[T]) minItems() int {
	return t.degree - 1
}

func (c *copyOnWriteContext[T]) newNode() (n *node[T]) {
	n = c.freelist.newNode()
	n.cow = c
	return
}

type freeType int

const (
	ftFreelistFull freeType = iota
	ftStored
	ftNotOwned
)

func (c *copyOnWriteContext[T]) freeNode(n *node[T]) freeType {
	if n.cow == c {
		n.items.truncate(0)
		n.children.truncate(0)
		n.cow = nil
		if c.freelist.freeNode(n) {
			return ftStored
		} else {
			return ftFreelistFull
		}
	} else {
		return ftNotOwned
	}
}

func (t *BTreeG[T]) ReplaceOrInsert(item T) (_ T, _ bool) {
	if t.root == nil {
		t.root = t.cow.newNode()
		t.root.items = append(t.root.items, item)
		t.length++
		return
	} else {
		t.root = t.root.mutableFor(t.cow)
		if len(t.root.items) >= t.maxItems() {
			item2, second := t.root.split(t.maxItems() / 2)
			oldroot := t.root
			t.root = t.cow.newNode()
			t.root.items = append(t.root.items, item2)
			t.root.children = append(t.root.children, oldroot, second)
		}
	}
	out, outb := t.root.insert(item, t.maxItems())
	if !outb {
		t.length++
	}
	return out, outb
}

func (t *BTreeG[T]) Delete(item T) (T, bool) {
	return t.deleteItem(item, removeItem)
}

func (t *BTreeG[T]) DeleteMin() (T, bool) {
	var zero T
	return t.deleteItem(zero, removeMin)
}

func (t *BTreeG[T]) DeleteMax() (T, bool) {
	var zero T
	return t.deleteItem(zero, removeMax)
}

func (t *BTreeG[T]) deleteItem(item T, typ toRemove) (_ T, _ bool) {
	if t.root == nil || len(t.root.items) == 0 {
		return
	}
	t.root = t.root.mutableFor(t.cow)
	out, outb := t.root.remove(item, t.minItems(), typ)
	if len(t.root.items) == 0 && len(t.root.children) > 0 {
		oldroot := t.root
		t.root = t.root.children[0]
		t.cow.freeNode(oldroot)
	}
	if outb {
		t.length--
	}
	return out, outb
}

func (t *BTreeG[T]) AscendRange(greaterOrEqual, lessThan T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, optional[T](greaterOrEqual), optional[T](lessThan), true, false, iterator)
}

func (t *BTreeG[T]) AscendLessThan(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, empty[T](), optional(pivot), false, false, iterator)
}

func (t *BTreeG[T]) AscendGreaterOrEqual(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, optional[T](pivot), empty[T](), true, false, iterator)
}

func (t *BTreeG[T]) Ascend(iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, empty[T](), empty[T](), false, false, iterator)
}

func (t *BTreeG[T]) DescendRange(lessOrEqual, greaterThan T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, optional[T](lessOrEqual), optional[T](greaterThan), true, false, iterator)
}

func (t *BTreeG[T]) DescendLessOrEqual(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, optional[T](pivot), empty[T](), true, false, iterator)
}

func (t *BTreeG[T]) DescendGreaterThan(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, empty[T](), optional[T](pivot), false, false, iterator)
}

func (t *BTreeG[T]) Descend(iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, empty[T](), empty[T](), false, false, iterator)
}

func (t *BTreeG[T]) Get(key T) (_ T, _ bool) {
	if t.root == nil {
		return
	}
	return t.root.get(key)
}

func (t *BTreeG[T]) Min() (_ T, _ bool) {
	return min(t.root)
}

func (t *BTreeG[T]) Max() (_ T, _ bool) {
	return max(t.root)
}

func (t *BTreeG[T]) Has(key T) bool {
	_, ok := t.Get(key)
	return ok
}

func (t *BTreeG[T]) Len() int {
	return t.length
}

func (t *BTreeG[T]) Clear(addNodesToFreelist bool) {
	if t.root != nil && addNodesToFreelist {
		t.root.reset(t.cow)
	}
	t.root, t.length = nil, 0
}

func (n *node[T]) reset(c *copyOnWriteContext[T]) bool {
	for _, child := range n.children {
		if !child.reset(c) {
			return false
		}
	}
	return c.freeNode(n) != ftFreelistFull
}

var itemLess LessFunc[Item] = func(a, b Item) bool {
	return a.Less(b)
}

type FreeList FreeListG[Item]

func NewFreeList(size int) *FreeList {
	return (*FreeList)(NewFreeListG[Item](size))
}

func NewWithFreeList(degree int, f *FreeList) *BTree {
	return (*BTree)(NewWithFreeListG[Item](degree, itemLess, (*FreeListG[Item])(f)))
}

type ItemIterator ItemIteratorG[Item]

func (t *BTree) Clone() (t2 *BTree) {
	return (*BTree)((*BTreeG[Item])(t).Clone())
}

func (t *BTree) Remove(item Item) Item {
	i, _ := (*BTreeG[Item])(t).Delete(item)
	return i
}

func (t *BTree) DeleteMax() Item {
	i, _ := (*BTreeG[Item])(t).DeleteMax()
	return i
}

func (t *BTree) DeleteMin() Item {
	i, _ := (*BTreeG[Item])(t).DeleteMin()
	return i
}

func (t *BTree) Get(key Item) Item {
	i, _ := (*BTreeG[Item])(t).Get(key)
	return i
}

func (t *BTree) Max() Item {
	i, _ := (*BTreeG[Item])(t).Max()
	return i
}

func (t *BTree) Min() Item {
	i, _ := (*BTreeG[Item])(t).Min()
	return i
}

func (t *BTree) Has(key Item) bool {
	return (*BTreeG[Item])(t).Has(key)
}

func (t *BTree) Insert(item Item) Item {
	i, _ := (*BTreeG[Item])(t).ReplaceOrInsert(item)
	return i
}

func (t *BTree) AscendRange(greaterOrEqual, lessThan Item, iterator ItemIterator) {
	(*BTreeG[Item])(t).AscendRange(greaterOrEqual, lessThan, (ItemIteratorG[Item])(iterator))
}

func (t *BTree) AscendLessThan(pivot Item, iterator ItemIterator) {
	(*BTreeG[Item])(t).AscendLessThan(pivot, (ItemIteratorG[Item])(iterator))
}

func (t *BTree) AscendGreaterOrEqual(pivot Item, iterator ItemIterator) {
	(*BTreeG[Item])(t).AscendGreaterOrEqual(pivot, (ItemIteratorG[Item])(iterator))
}

func (t *BTree) Ascend(iterator ItemIterator) {
	(*BTreeG[Item])(t).Ascend((ItemIteratorG[Item])(iterator))
}

func (t *BTree) DescendRange(lessOrEqual, greaterThan Item, iterator ItemIterator) {
	(*BTreeG[Item])(t).DescendRange(lessOrEqual, greaterThan, (ItemIteratorG[Item])(iterator))
}

func (t *BTree) DescendLessOrEqual(pivot Item, iterator ItemIterator) {
	(*BTreeG[Item])(t).DescendLessOrEqual(pivot, (ItemIteratorG[Item])(iterator))
}

func (t *BTree) DescendGreaterThan(pivot Item, iterator ItemIterator) {
	(*BTreeG[Item])(t).DescendGreaterThan(pivot, (ItemIteratorG[Item])(iterator))
}

func (t *BTree) Descend(iterator ItemIterator) {
	(*BTreeG[Item])(t).Descend((ItemIteratorG[Item])(iterator))
}

func (t *BTree) Size() int {
	return (*BTreeG[Item])(t).Len()
}

func (t *BTree) Clear(addNodesToFreelist bool) {
	(*BTreeG[Item])(t).Clear(addNodesToFreelist)
}

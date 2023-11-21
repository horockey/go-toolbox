package comparer

var _ Comparer[any] = &pureComparer[any]{}

type pureComparer[T any] struct {
	compFunc func(a, b T) int
}

func NewPureComparer[T any](compFunc func(a, b T) int) *pureComparer[T] {
	return &pureComparer[T]{
		compFunc: compFunc,
	}
}

func (comp *pureComparer[T]) Compare(a, b T) int {
	return comp.compFunc(a, b)
}

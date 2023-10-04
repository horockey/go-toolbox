package comparer

type Comparer[T any] interface {
	// Returns
	// -1 if a is greater than b
	// 0 if they are equal
	// 1 if b is greater then a
	Compare(a, b T) int
}

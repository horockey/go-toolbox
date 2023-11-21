package comparer

type Comparer[T any] interface {
	// Returns
	// -1 if a > b
	// 0 if a = b
	// 1 if a < b
	Compare(a, b T) int
}

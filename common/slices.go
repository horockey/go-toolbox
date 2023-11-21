package common

// Returns concatentaion of given slices.
func ConcatSlices[T any](sls ...[]T) []T {
	res := []T{}
	for _, sl := range sls {
		res = append(res, sl...)
	}
	return res
}

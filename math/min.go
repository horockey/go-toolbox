package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Returns minimum of given numbers, or NaN, if numbers is empty.
func Min[T constraints.Float | constraints.Integer](numbers ...T) T {
	if len(numbers) == 0 {
		return T(math.NaN())
	}
	extr := numbers[0]
	for _, number := range numbers {
		extr = min2(extr, number)
	}
	return extr
}

func min2[T constraints.Float | constraints.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

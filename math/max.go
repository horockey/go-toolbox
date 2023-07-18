package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Returns maximum of given numbers, or NaN, if numbers is empty.
func Max[T constraints.Float | constraints.Integer](numbers ...T) T {
	if len(numbers) == 0 {
		return T(math.NaN())
	}
	extr := numbers[0]
	for _, number := range numbers {
		extr = max2(extr, number)
	}
	return extr
}

func max2[T constraints.Float | constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

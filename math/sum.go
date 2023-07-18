package math

import "golang.org/x/exp/constraints"

// Returns sum of given numbers.
func Sum[T constraints.Float | constraints.Integer](numbers ...T) T {
	var sum T = 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

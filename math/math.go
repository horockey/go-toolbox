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

// Returns sum of given numbers.
func Sum[T constraints.Float | constraints.Integer](numbers ...T) T {
	var sum T = 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// Returns greater commn divisor (GCD) of given numbers, or NaN, if numbers is empty.
func GCD[T constraints.Integer](numbers ...T) T {
	if len(numbers) == 0 {
		return T(math.NaN())
	}
	if len(numbers) == 1 {
		return numbers[0]
	}
	lastGCD := Min(numbers...)
	for idx := 0; idx < len(numbers)-1; idx++ {
		lastGCD = min2(lastGCD, gcd2(numbers[idx], numbers[idx+1]))
	}
	return lastGCD
}

func min2[T constraints.Float | constraints.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max2[T constraints.Float | constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func gcd2[T constraints.Integer](a, b T) T {
	a, b = max2(a, b), min2(a, b)
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

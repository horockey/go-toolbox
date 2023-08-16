package string_comparer_test

import (
	"testing"

	"github.com/horockey/go-toolbox/datastructs/pkg/comparer/string_comparer"
	"github.com/stretchr/testify/require"
)

var cases = []struct {
	l           string
	r           string
	expectation int
}{
	{"1", "1", 0},
	{"1", "2", 1},
	{"1", "0", -1},
	{"10", "2", 1},
	{"A", "a", 1},
}

func TestStringComparer(t *testing.T) {
	comp := string_comparer.New()

	for _, c := range cases {
		res := comp.Compare(c.l, c.r)
		require.Equal(t, c.expectation, res)
	}
}

package red_black_tree_test

import (
	"strconv"
	"testing"

	"github.com/horockey/go-toolbox/datastructs/trees/red_black_tree"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	tree := red_black_tree.New[int]()

	const N = 10

	for i := 0; i < N; i++ {
		err := tree.Insert(strconv.Itoa(i), i)
		require.NoError(t, err)
	}

	require.Equal(t, N, tree.Size())
}

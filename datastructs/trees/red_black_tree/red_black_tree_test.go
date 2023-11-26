package red_black_tree

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUD(t *testing.T) {
	tree := New[int]()

	const N = 10

	t.Run("test_insert", func(t *testing.T) {
		for i := 0; i < N; i++ {
			err := tree.Insert(strconv.Itoa(i), i)
			require.NoError(t, err)

			_, err = blackHeight(tree.root)
			require.NoError(t, err)
			require.Equal(t, ColorBlack, tree.root.color)
		}

		require.Equal(t, N, tree.Size())
	})

	t.Run("test_read", func(t *testing.T) {
		for i := 0; i < N; i++ {
			val, err := tree.Get(strconv.Itoa(i))
			require.NoError(t, err)
			require.Equal(t, i, val)
		}

		require.Equal(t, N, tree.Size())
	})

	t.Run("test_remove", func(t *testing.T) {
		for i := 0; i < N; i++ {
			err := tree.Remove(strconv.Itoa(i))
			require.NoError(t, err)

			_, err = blackHeight(tree.root)
			require.NoError(t, err)
			if tree.root != nil {
				require.Equal(t, ColorBlack, tree.root.color)
			}
		}

		require.Equal(t, 0, tree.Size())
	})
}

func TestClear(t *testing.T) {
	tree := New[int]()

	err := tree.Insert("1", 1)
	require.NoError(t, err)
	require.Equal(t, 1, tree.Size())

	tree.Clear()

	require.Equal(t, 0, tree.Size())
	require.Nil(t, tree.root)
}

func blackHeight[K, V any](subroot *node[K, V]) (int, error) {
	if subroot == nil {
		return 0, nil
	}

	if subroot.hasNoChildren() {
		return 1, nil
	}

	leftBlackHeight, err := blackHeight(subroot.left)
	if err != nil {
		return 0, err
	}
	if subroot.left == nil || subroot.left.color == ColorBlack {
		leftBlackHeight++
	}

	rightBlackHeight, err := blackHeight(subroot.right)
	if err != nil {
		return 0, err
	}
	if subroot.right == nil || subroot.right.color == ColorBlack {
		rightBlackHeight++
	}

	if leftBlackHeight != rightBlackHeight {
		return 0, fmt.Errorf(
			"tree disballanced in node with key %+v(%s): %d/%d",
			subroot.Key,
			subroot.color.String(),
			leftBlackHeight,
			rightBlackHeight,
		)
	}

	return leftBlackHeight, nil
}

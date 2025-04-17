package avl_tree_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/horockey/go-toolbox/datastructs/comparer"
	"github.com/horockey/go-toolbox/datastructs/trees/avl_tree"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	// in ballanced tree root's height H = floor(log2(n)) +-2
	cases := []struct {
		numberOfElements uint
		expectedHeight   uint
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{10, 4},
		{1_000, 11},
		{100_000, 17},
	}

	for _, c := range cases {
		tree := avl_tree.New[int]()
		for i := 0; i < int(c.numberOfElements); i++ {
			err := tree.Insert(strconv.Itoa(i), i)
			require.NoError(t, err)
		}
		require.InDelta(t, c.expectedHeight, tree.Height(), 2.0)
	}
}

func TestRemove(t *testing.T) {
	const N = 10

	tree := avl_tree.New[int]()
	for i := 0; i < N; i++ {
		err := tree.Insert(strconv.Itoa(i), i)
		require.NoError(t, err)
	}

	for i := 0; i < N; i++ {
		err := tree.Remove(strconv.Itoa(i))
		require.NoError(t, err)
		require.InDelta(
			t,
			math.Floor(max(0, math.Log2(float64(tree.Size())))),
			tree.Height(),
			2.0,
		)
	}
}

func TestGet(t *testing.T) {
	tree := avl_tree.New[int]()
	expected := struct {
		key string
		val int
	}{
		key: "1",
		val: 1,
	}

	err := tree.Insert(expected.key, expected.val)
	require.NoError(t, err)
	val, err := tree.Get(expected.key)
	require.NoError(t, err)
	require.Equal(t, expected.val, val)
}

func TestCustomKey(t *testing.T) {
	type Key struct {
		Foo string
		Bar int
	}

	keyCompFunc := func(a, b Key) int {
		if a.Foo < b.Foo {
			return -1
		} else if a.Foo > b.Foo {
			return 1
		} else {
			if a.Bar < b.Bar {
				return -1
			} else if a.Bar > b.Bar {
				return 1
			} else {
				return 0
			}
		}
	}

	tree := avl_tree.NewWithCustomKey[Key, string](comparer.NewPureComparer(keyCompFunc))

	expected := struct {
		Key   Key
		Value string
	}{
		Key:   Key{Foo: "foo", Bar: -17},
		Value: "testVal",
	}

	err := tree.Insert(expected.Key, expected.Value)
	require.NoError(t, err)
	val, err := tree.Get(expected.Key)
	require.NoError(t, err)
	require.Equal(t, expected.Value, val)
}

package btree

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const degree int = 3

type Key int

func (k Key) Less(other Item) bool {
	otherK, ok := other.(Key)
	if !ok {
		panic(errors.New("assertion failed"))
	}

	return k < otherK
}

func TestCRUD(t *testing.T) {
	const N int = 10

	tr := New(degree)

	t.Run("test insert", func(t *testing.T) {
		for i := 0; i < N; i++ {
			tr.Insert(Key(i))
		}

		require.Equal(t, N, tr.Size())
	})

	t.Run("test get", func(t *testing.T) {
		for i := 0; i < N; i++ {
			it := tr.Get(Key(i))
			require.Equal(t, Key(i), it.(Key))
		}
	})

	t.Run("test remove", func(t *testing.T) {
		for i := 0; i < N; i++ {
			tr.Remove(Key(i))
		}
		require.Equal(t, 0, tr.Size())
	})
}

package vector_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ykhrustalev/go-vector"
	"testing"
)

func requireError(t *testing.T, actualErr, expectedError error) {
	require.Error(t, actualErr)
	require.Equal(t, expectedError, actualErr)
}

func requireVector(t *testing.T, v *vector.Vector, expectedCap int, expectedItems []int) {
	require.Equal(t, expectedCap, v.Cap())
	require.Equal(t, expectedItems, v.Slice())
}

func TestFrom(t *testing.T) {
	v := vector.From(1, 2, 3, 4, 5)
	requireVector(t, v, 10, []int{1, 2, 3, 4, 5})

	v.Append(6)
	v.Append(7)
	v.Append(8)
	v.Append(9)
	v.Append(10)
	v.Append(11)
	requireVector(t, v, 20, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
}

func TestNew(t *testing.T) {
	v := vector.New()
	requireVector(t, v, 10, nil)

	v.Append(1)
	requireVector(t, v, 10, []int{1})

	for i := 2; i < 12; i++ {
		v.Append(i)
	}
	requireVector(t, v, 20, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})
}

func TestNewWithCap(t *testing.T) {
	v := vector.NewWithCap(2, 3)
	requireVector(t, v, 2, nil)

	v.Append(1)
	requireVector(t, v, 2, []int{1})
	v.Append(2)
	requireVector(t, v, 2, []int{1, 2})

	v.Append(3)
	requireVector(t, v, 6, []int{1, 2, 3})
}

func TestVector_Add(t *testing.T) {
	v := vector.NewWithCap(2, 3)

	requireError(t, v.Add(1, 11), vector.ErrInvalidIndex)

	require.NoError(t, v.Add(0, 11))
	requireVector(t, v, 2, []int{11})

	require.NoError(t, v.Add(0, 12))
	requireVector(t, v, 2, []int{12, 11})

	require.NoError(t, v.Add(1, 13))
	requireVector(t, v, 6, []int{12, 13, 11})

	require.NoError(t, v.Add(2, 14))
	requireVector(t, v, 6, []int{12, 13, 14, 11})

	requireError(t, v.Add(-1, 11), vector.ErrInvalidIndex)
	requireError(t, v.Add(10, 11), vector.ErrInvalidIndex)
}

func TestVector_Set(t *testing.T) {
	v := vector.NewWithCap(2, 3)

	requireError(t, v.Set(0, 11), vector.ErrInvalidIndex)

	v.AppendAll(1, 2, 3)

	require.NoError(t, v.Set(0, 9))
	requireVector(t, v, 6, []int{9, 2, 3})

	requireError(t, v.Set(-1, 11), vector.ErrInvalidIndex)
	requireError(t, v.Set(99, 11), vector.ErrInvalidIndex)
}

func TestVector_Append(t *testing.T) {
	v := vector.NewWithCap(2, 3)

	v.Append(10)
	v.Append(11)
	require.Equal(t, []int{10, 11}, v.Slice())
	require.Equal(t, v.Cap(), 2)

	v.Append(12)
	require.Equal(t, []int{10, 11, 12}, v.Slice())
	require.Equal(t, v.Cap(), 6)
}

func TestVector_AppendAll(t *testing.T) {
	v := vector.NewWithCap(2, 3)

	v.AppendAll(10, 11)
	require.Equal(t, []int{10, 11}, v.Slice())
	require.Equal(t, v.Cap(), 2)

	v.AppendAll(12)
	require.Equal(t, []int{10, 11, 12}, v.Slice())
	require.Equal(t, v.Cap(), 6)

	v.AppendAll(13, 14, 15, 16, 17, 18, 19, 20)
	require.Equal(t, []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, v.Slice())
	require.Equal(t, v.Cap(), 18)
}

func TestVector_Cap(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		v := vector.New()
		require.Equal(t, 10, v.Cap())
	})

	t.Run("with cap1", func(t *testing.T) {
		v := vector.NewWithCap(2, 3)
		require.Equal(t, 2, v.Cap())
	})

	t.Run("with cap2", func(t *testing.T) {
		v := vector.NewWithCap(99, 3)
		require.Equal(t, 99, v.Cap())
	})
}

func TestVector_Clear(t *testing.T) {
	v := vector.From(1, 2, 3, 4, 5)
	requireVector(t, v, 10, []int{1, 2, 3, 4, 5})

	v.Clear()
	requireVector(t, v, 10, nil)
}

func TestVector_Len(t *testing.T) {
	v := vector.From(1, 2, 3, 4, 5)
	require.Equal(t, 5, v.Len())

	v.Append(6)
	require.Equal(t, 6, v.Len())
}

func TestVector_Peek(t *testing.T) {
	requireSuccess := func(t *testing.T, v *vector.Vector, index int, expectedItem int) {
		item, err := v.Peek(index)
		require.NoError(t, err)
		require.Equal(t, expectedItem, item)
	}

	requireError := func(t *testing.T, v *vector.Vector, index int) {
		_, err := v.Peek(index)
		requireError(t, err, vector.ErrInvalidIndex)
	}

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4, 5)

		requireSuccess(t, v, 0, 1)
		requireSuccess(t, v, 2, 3)
		requireSuccess(t, v, 4, 5)

		requireError(t, v, -1)
		requireError(t, v, 99)
	})

	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		requireError(t, v, 0)
	})
}

func TestVector_IndexOf(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		require.Equal(t, -1, v.IndexOf(1))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4, 5, 1, 2, 3, 4, 5)
		require.Equal(t, -1, v.IndexOf(0))
		require.Equal(t, 0, v.IndexOf(1))
		require.Equal(t, 2, v.IndexOf(3))
		require.Equal(t, -1, v.IndexOf(6))
	})
}

func TestVector_Slice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		require.Equal(t, []int(nil), v.Slice())
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3)
		require.Equal(t, []int{1, 2, 3}, v.Slice())
	})
}

func TestVector_Clone(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		clone := v.Clone()
		requireVector(t, clone, 10, []int(nil))
		v.Append(1)
		requireVector(t, clone, 10, []int(nil))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3)
		clone := v.Clone()
		requireVector(t, clone, 10, []int{1, 2, 3})
		v.Append(1)
		requireVector(t, clone, 10, []int{1, 2, 3})
	})
}

func TestVector_Remove(t *testing.T) {
	requireSuccess := func(
		t *testing.T,
		v *vector.Vector,
		index int,
		expectedItem int,
		expectedCap int,
		expectedItems []int,
	) {
		item, err := v.Remove(index)
		require.NoError(t, err)
		requireVector(t, v, expectedCap, expectedItems)
		require.Equal(t, expectedItem, item)
	}

	requireError := func(t *testing.T, v *vector.Vector, index int) {
		_, err := v.Remove(index)
		requireError(t, err, vector.ErrInvalidIndex)
	}

	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		requireError(t, v, 0)
		requireError(t, v, 1)
		requireError(t, v, -1)
	})

	t.Run("empty", func(t *testing.T) {
		v := vector.From(10, 20, 30, 40, 50, 60, 70, 80, 90)
		requireError(t, v, -1)
		requireError(t, v, 99)
		requireSuccess(t, v, 8, 90, 10, []int{10, 20, 30, 40, 50, 60, 70, 80})
		requireSuccess(t, v, 7, 80, 10, []int{10, 20, 30, 40, 50, 60, 70})
		requireSuccess(t, v, 0, 10, 10, []int{20, 30, 40, 50, 60, 70})
		requireSuccess(t, v, 0, 20, 10, []int{30, 40, 50, 60, 70})
		requireSuccess(t, v, 1, 40, 10, []int{30, 50, 60, 70})
		requireSuccess(t, v, 2, 60, 10, []int{30, 50, 70})
	})
}

func TestVector_Each(t *testing.T) {
	requireSuccess := func(t *testing.T, v *vector.Vector, expectedItems []int) {
		var actual []int
		v.Each(func(index, item int) bool {
			actual = append(actual, item)
			return true
		})
		require.Equal(t, expectedItems, actual)
	}

	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		requireSuccess(t, v, []int(nil))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3)
		requireSuccess(t, v, []int{1, 2, 3})
	})

	t.Run("stop", func(t *testing.T) {
		v := vector.From(1, 2, 3)

		var actual []int
		v.Each(func(index, item int) bool {
			actual = append(actual, item)
			if item == 2 {
				return false
			}
			return true
		})
		require.Equal(t, actual, []int{1, 2})
	})
}

func TestVector_InnerProduct(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v1 := vector.New()
		v2 := vector.New()

		actual, err := v1.InnerProduct(v2)
		require.NoError(t, err)
		require.Equal(t, 0, actual)
	})

	t.Run("with items", func(t *testing.T) {
		v1 := vector.From(1, 2, 3)
		v2 := vector.From(2, 3, 4)

		actual, err := v1.InnerProduct(v2)
		require.NoError(t, err)
		require.Equal(t, 1*2+2*3+3*4, actual)
	})

	t.Run("different size", func(t *testing.T) {
		v1 := vector.New()
		v2 := vector.From(2, 3, 4)

		_, err := v1.InnerProduct(v2)
		requireError(t, err, vector.ErrSizeDiffers)
	})
}

func TestVector_Any(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := vector.New()

		require.False(t, v.Any(func(item int) bool { return true }))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4)

		require.False(t, v.Any(func(item int) bool { return false }))
		require.True(t, v.Any(func(item int) bool { return item == 3 }))
		require.True(t, v.Any(func(item int) bool { return item > 3 }))
	})
}

func TestVector_All(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := vector.New()

		require.False(t, v.All(func(item int) bool { return true }))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4)

		require.False(t, v.All(func(item int) bool { return false }))
		require.False(t, v.All(func(item int) bool { return item == 3 }))
		require.True(t, v.All(func(item int) bool { return item > 0 }))
	})
}

func TestVector_RemoveIf(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := vector.New()

		v.RemoveIf(func(item int) bool { return true })
		requireVector(t, v, 10, []int(nil))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4)

		v.RemoveIf(func(item int) bool { return item < 3 })
		requireVector(t, v, 10, []int{3, 4})

		v.RemoveIf(func(item int) bool { return item > 3 })
		requireVector(t, v, 10, []int{3})

		v.RemoveIf(func(item int) bool { return true })
		requireVector(t, v, 10, []int(nil))
	})
}

func TestVector_Equal(t *testing.T) {
	t.Run("equal empty", func(t *testing.T) {
		v1 := vector.New()
		v2 := vector.New()

		require.True(t, v1.Equal(v2))
	})

	t.Run("equal with items", func(t *testing.T) {
		v1 := vector.From(1, 2, 3)
		v2 := vector.From(1, 2, 3)

		require.True(t, v1.Equal(v2))
	})

	t.Run("not equal one empty", func(t *testing.T) {
		v1 := vector.New()
		v2 := vector.From(1, 2, 3)

		require.False(t, v1.Equal(v2))
	})

	t.Run("not equal diff len", func(t *testing.T) {
		v1 := vector.From(1, 2)
		v2 := vector.From(1, 2, 3)

		require.False(t, v1.Equal(v2))
	})

	t.Run("not equal same len", func(t *testing.T) {
		v1 := vector.From(1, 2, 4)
		v2 := vector.From(1, 2, 3)

		require.False(t, v1.Equal(v2))
	})
}

func TestVector_Accumulate(t *testing.T) {
	sum := func(a, b int) int { return a + b }

	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		assert.Equal(t, vector.New().Slice(), v.Accumulate(sum).Slice())
	})

	t.Run("single items", func(t *testing.T) {
		v := vector.From(9)
		assert.Equal(t, vector.From(9).Slice(), v.Accumulate(sum).Slice())
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4, 5)
		assert.Equal(t, vector.From(1, 3, 6, 10, 15).Slice(), v.Accumulate(sum).Slice())
	})
}

func TestVector_Reduce(t *testing.T) {
	sum := func(a, b int) int { return a + b }

	t.Run("empty", func(t *testing.T) {
		v := vector.New()
		assert.Equal(t, 0, v.Reduce(sum))
	})

	t.Run("single items", func(t *testing.T) {
		v := vector.From(9)
		assert.Equal(t, 9, v.Reduce(sum))
	})

	t.Run("with items", func(t *testing.T) {
		v := vector.From(1, 2, 3, 4, 5)
		assert.Equal(t, 15, v.Reduce(sum))
	})
}

package slices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatten(t *testing.T) {
	t.Run("flattens slices of ints", func(t *testing.T) {
		input := [][]int{{1, 2}, {3, 4, 5}, {}}
		expected := []int{1, 2, 3, 4, 5}
		result := Flatten(input)
		assert.Equal(t, expected, result)
	})

	t.Run("flattens slices of strings", func(t *testing.T) {
		input := [][]string{{"a", "b"}, {}, {"c"}}
		expected := []string{"a", "b", "c"}
		result := Flatten(input)
		assert.Equal(t, expected, result)
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		input := [][]int{}
		expected := []int{}
		result := Flatten(input)
		assert.Equal(t, expected, result)
	})

	t.Run("returns empty slice for slices of empty slices", func(t *testing.T) {
		input := [][]int{{}, {}, {}}
		expected := []int{}
		result := Flatten(input)
		assert.Equal(t, expected, result)
	})

	t.Run("preserves order of elements", func(t *testing.T) {
		input := [][]int{{1}, {2, 3}, {4}}
		expected := []int{1, 2, 3, 4}
		result := Flatten(input)
		assert.Equal(t, expected, result)
	})
}

func TestCompareStringArraysMatching(t *testing.T) {
	areEqual := EqualsIgnoreOrder([]string{"Hello", "World"}, []string{"Hello", "World"})
	assert.True(t, areEqual, "The arrays are equal")
	areEqual = EqualsIgnoreOrder([]string{"World", "Hello"}, []string{"Hello", "World"})
	assert.True(t, areEqual, "The arrays are equal")
	areEqual = EqualsIgnoreOrder([]string{"Hello", "Hello"}, []string{"Hello", "Hello"})
	assert.True(t, areEqual, "The arrays are equal")
	areEqual = EqualsIgnoreOrder([]string{}, []string{})
	assert.True(t, areEqual, "The arrays are equal")

	areEqual = EqualsIgnoreOrder([]int{1, 3, 3, 7}, []int{1, 3, 3, 7})
	assert.True(t, areEqual, "The arrays are equal")
	areEqual = EqualsIgnoreOrder([]int{3, 7, 3, 1}, []int{1, 3, 3, 7})
	assert.True(t, areEqual, "The arrays are equal")
	areEqual = EqualsIgnoreOrder([]int{3, 3, 3, 3}, []int{3, 3, 3, 3})
	assert.True(t, areEqual, "The arrays are equal")
	areEqual = EqualsIgnoreOrder([]int{}, []int{})
	assert.True(t, areEqual, "The arrays are equal")
}

func TestCompareStringArraysNotMatchingSize(t *testing.T) {
	areEqual := EqualsIgnoreOrder([]string{"Hello"}, []string{"Hello", "World"})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]string{"World", "Hello"}, []string{"Hello", "World", "Test"})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]string{"World", "Hello", "Here"}, []string{"Hello", "World"})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]string{"World", "Hello", "Hello", "Here"}, []string{"World", "Hello", "Here", "Here"})
	assert.False(t, areEqual, "The arrays are not equal")

	areEqual = EqualsIgnoreOrder([]int{1}, []int{1, 3})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]int{1, 3}, []int{1, 3, 3})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]int{1, 3, 3}, []int{1, 3})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]int{1, 2, 2, 3}, []int{1, 2, 3, 3})
	assert.False(t, areEqual, "The arrays are not equal")
}

func TestCompareStringArraysNotMatching(t *testing.T) {
	areEqual := EqualsIgnoreOrder([]string{"Hello", "World1"}, []string{"Hello", "World"})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]string{"World", "Hello"}, []string{"hello", "world"})
	assert.False(t, areEqual, "The arrays are not equal")

	areEqual = EqualsIgnoreOrder([]int{1, 3, 3, 7}, []int{1, 2, 3, 7})
	assert.False(t, areEqual, "The arrays are not equal")
	areEqual = EqualsIgnoreOrder([]int{1, 7, 7, 7}, []int{1, 2, 3, 7})
	assert.False(t, areEqual, "The arrays are not equal")
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int{1, 2, 3}, 2))
	assert.False(t, Contains([]int{1, 2, 3}, 4))
	assert.False(t, Contains([]int{}, 1))
	assert.True(t, Contains([]int{5}, 5))
	assert.False(t, Contains([]int{5}, 6))
}

func TestContainsCount(t *testing.T) {
	assert.Equal(t, 3, ContainsCount([]string{"a", "b", "a", "c", "a"}, "a"))
	assert.Equal(t, 0, ContainsCount([]string{"x", "y", "z"}, "a"))
	assert.Equal(t, 0, ContainsCount([]string{}, "a"))
	assert.Equal(t, 1, ContainsCount([]string{"a"}, "a"))
}

func TestIndexOf(t *testing.T) {
	assert.Equal(t, 1, IndexOf([]int{10, 20, 30}, 20))
	assert.Equal(t, -1, IndexOf([]int{10, 20, 30}, 40))
	assert.Equal(t, -1, IndexOf([]int{}, 10))
	assert.Equal(t, 0, IndexOf([]int{5, 6, 7}, 5))
}

func TestIndexOfAll(t *testing.T) {
	assert.ElementsMatch(t, []int{0, 2, 4}, IndexOfAll([]int{1, 2, 1, 3, 1}, 1))
	assert.Empty(t, IndexOfAll([]int{4, 5, 6}, 1))
	assert.Empty(t, IndexOfAll([]int{}, 1))
	assert.Equal(t, []int{0}, IndexOfAll([]int{7}, 7))
}

func TestUnique(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, Unique([]string{"a", "b", "a", "c", "b"}))
	assert.Equal(t, []string{"x", "y", "z"}, Unique([]string{"x", "y", "z"}))
	assert.Empty(t, Unique([]string{}))
	assert.Equal(t, []string{"a"}, Unique([]string{"a", "a", "a"}))
}

func TestReverse(t *testing.T) {
	var sNil []int
	Reverse(sNil)
	assert.Empty(t, sNil)

	sEmpty := []int{}
	Reverse(sEmpty)
	assert.Empty(t, sEmpty)

	sSingle := []int{42}
	Reverse(sSingle)
	assert.Equal(t, []int{42}, sSingle)

	sliceItems := []int{1, 2, 3}
	Reverse(sliceItems)
	assert.Equal(t, []int{3, 2, 1}, sliceItems)

}

func TestChunk(t *testing.T) {
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, Chunk([]int{1, 2, 3, 4, 5}, 2))
	assert.Equal(t, [][]int{{1, 2}, {3, 4}}, Chunk([]int{1, 2, 3, 4}, 2))
	assert.Equal(t, [][]int{{1, 2}}, Chunk([]int{1, 2}, 5))
	assert.Empty(t, Chunk([]int{}, 3))
}

func TestChunkPanics(t *testing.T) {
	assert.Panics(t, func() {
		Chunk([]int{1, 2, 3}, 0)
	})
}

func TestChunk_BoundaryCondition(t *testing.T) {

	// Slice length is exactly divisible by chunk size
	slice := []int{1, 2, 3, 4}
	size := 2

	// end == len(s) will be true on the last iteration
	// So the condition `end > len(s)` should be false
	expected := [][]int{{1, 2}, {3, 4}}

	actual := Chunk(slice, size)
	assert.Equal(t, expected, actual, "Chunk should split evenly without exceeding slice bounds")
}

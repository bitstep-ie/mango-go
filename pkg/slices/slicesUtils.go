// Package slices provides helpful functions for slices
package slices

// EqualsIgnoreOrder compares two splices ignoring the order of elements
// It returns true if both slices contain the same elements (and same frequency of same elements), regardless of order
func EqualsIgnoreOrder[ComparableSlice ~[]Type, Type comparable](s1 ComparableSlice, s2 ComparableSlice) bool {
	if len(s1) != len(s2) {
		return false
	}

	frequency := make(map[Type]int)
	for _, key := range s1 {
		frequency[key]++
	}

	for _, key := range s2 {
		if frequency[key] == 0 {
			return false
		}
		frequency[key]--
	}

	return true
}

// Contains returns true if value is found
// Empty slice always returns false
func Contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// ContainsCount gives back the count of value findings
func ContainsCount[T comparable](slice []T, value T) int {
	count := 0
	for _, x := range slice {
		if x == value {
			count++
		}
	}
	return count
}

// IndexOf returns index of first match. Similar to strings.Index
// Empty slice / no finding returns -1
func IndexOf[T comparable](slice []T, value T) int {
	for i, x := range slice {
		if x == value {
			return i
		}
	}
	return -1
}

// IndexOfAll returns all positions of value findings
func IndexOfAll[T comparable](slice []T, value T) []int {
	var indices []int
	for i, x := range slice {
		if x == value {
			indices = append(indices, i)
		}
	}
	return indices
}

// Unique returns a new slice containing only the unique values from slice (deduplication)
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Reverse flips the contents of the slice
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Chunk returns a slice of slices of size
// Each of size specified - last may be smaller
// Panics when size <= 0
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		panic("Chunk: size must be greater than 0")
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

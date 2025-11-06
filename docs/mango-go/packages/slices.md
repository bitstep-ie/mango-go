# mango4go - slices

The `slices` package is aimed to be a useful package for utilities working with `slices`.


## EqualsIgnoreOrder
It compares two splices ignoring the order of elements. Returns true if both slices contain the same elements (and same frequency of same elements), regardless of order of elements.

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/slices"

// ... rest of your code

// compare slices as needed and save the result
areEqual := slices.EqualsIgnoreOrder([]string{"World", "Hello"}, []string{"hello", "world"})
```

## Contains
Returns `true` if `value` is found. *Empty* slice always returns `false`.


## ContainsCount
Returns the `count` of `value` findings in the slice.


## IndexOf 
Returns `index` of *first match*. Similar to `strings.Index`. *Empty* slice / *no finding* returns -1.


## IndexOfAll
Returns all *indexes* of `value` findings, empty slice if no findings.
### Example
`IndexOfAll([]int{1, 2, 1, 3, 1}, 1))` 
will return 
`[]int{0, 2, 4}` which represent all the positions where the value for search (`1` can be found)


## Unique (deduplication)
Returns a new slice containing *only the unique values* from slice.


## Reverse
Flips the contents of the slice.


## Chunk 
Returns a slice of slices of size (last may be smaller). **Panics** when `size <= 0`.

### Example
`Chunk([]int{1, 2, 3, 4, 5}, 2))`
will return:
`[][]int{{1, 2}, {3, 4}, {5}}`
In plain english, get slices of size 2 from this larger slice.
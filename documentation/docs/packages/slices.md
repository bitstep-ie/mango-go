# `pkg/slices`

Generic helpers for common slice patterns—membership, counting, chunking, reversing, and deduplication—without re-writing loops at every callsite.

## Quick Start

```go
import mangoslices "github.com/bitstep-ie/mango-go/pkg/slices"

names := []string{"amy", "ben", "amy"}
if mangoslices.Contains(names, "ben") {
    idx := mangoslices.IndexOf(names, "amy") // 0
    freq := mangoslices.ContainsCount(names, "amy") // 2
    unique := mangoslices.Unique(names) // ["amy","ben"]
    chunks := mangoslices.Chunk(names, 2) // [["amy","ben"],["amy"]]
    _ = []any{idx, freq, unique, chunks}
}
```

## API Reference

| Function | Description |
| --- | --- |
| `EqualsIgnoreOrder(a, b)` | compare two slices while ignoring element order but respecting multiplicity |
| `Contains(slice, value)` | membership test |
| `ContainsCount(slice, value)` | count occurrences |
| `IndexOf(slice, value)` | index of first match (`-1` if missing) |
| `IndexOfAll(slice, value)` | every index that matches |
| `Unique(slice)` | deduplicated copy preserving first-seen order |
| `Reverse(slice)` | in-place reversal |
| `Chunk(slice, size)` | split into evenly sized batches (last chunk may be smaller) |

All helpers are generic (Go 1.18+) and therefore type-safe across ints, structs, strings, etc. Functions that mutate (`Reverse`) do so in-place; the rest allocate new slices as needed.

## Examples

### EqualsIgnoreOrder

```go
a := []int{1, 2, 2, 3}
b := []int{3, 2, 1, 2}
same := mangoslices.EqualsIgnoreOrder(a, b) // true
```

### IndexOfAll

```go
positions := mangoslices.IndexOfAll([]string{"ca", "ba", "ca"}, "ca")
// []int{0, 2}
```

### Chunk (Batching)

```go
jobs := []int{1, 2, 3, 4, 5}
batches := mangoslices.Chunk(jobs, 2)
// [][]int{{1,2}, {3,4}, {5}}
```

> `Chunk` panics when `size <= 0`; validate user input before calling it.

## Tips

- These helpers are a drop-in replacement for small anonymous loops, making tests and production code easier to read.
- Combine `Unique` + `Chunk` to dedupe and batch IDs before fan-out.
- `EqualsIgnoreOrder` is helpful in tests for comparing slices where ordering is non-deterministic.

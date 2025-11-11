# `pkg/random`

Convenience helpers for building fixtures, tokens, and fuzz data. Functions span deterministic math/rand utilities and crypto-safe password generation.

## Quick Start

```go
import (
    "math/rand"
    "time"
    mangorand "github.com/bitstep-ie/mango-go/pkg/random"
)

func init() {
    rand.Seed(time.Now().UnixNano()) // seed math/rand once per process
}

func example() {
    promo := mangorand.String(8)
    orderID := mangorand.Number[int](1000, 9999)
    mustRun := mangorand.Bool()
    secret := mangorand.Password(24, mangorand.PasswordOptions{
        Letters: true,
        Digits:  true,
        Symbols: true,
        Exclude: "O0l1",
    })
    _ = []any{promo, orderID, mustRun, secret}
}
```

## API Highlights

### Numbers & Booleans

- `Number[T Num](min, max T)` – inclusive range for ints/uints, `[min,max)` for floats. Automatically swaps min/max.
- `Sign()` – returns `+1` or `-1`.
- `Bool()` – coin flip.
- `Choice(slice []T)` – pick a random element, panics on empty slice.

### Bytes & Strings

- `Byte()` – cryptographically secure random byte.
- `String(n)` – alphanumeric (upper+lowercase, digits).
- `Alpha(n)` / `Numeric(n)` – letters-only or digits-only.
- `FromCharset(n, charset)` – supply your own rune set.

### Passwords

```go
pw := mangorand.Password(16, mangorand.PasswordOptions{
    Letters: true,
    Digits:  true,
    Symbols: false,
})
```

Options allow mixing letters, digits, symbols, plus excluding ambiguous characters. Generation uses `crypto/rand` so the output is safe for access tokens and human passwords. Panics if the resulting charset would be empty.

### Time Helpers

- `Date(min, max time.Time)` – inclusive random timestamp.
- `Duration(min, max time.Duration)` – inclusive random duration.

Both helpers swap arguments when `min > max`, making it easy to mock “recent” or “soon” values without extra branching.

## Tips

- Only the functions under “Numbers & Booleans” rely on `math/rand`; seed it once during startup for non-deterministic sequences.
- `Choice` and `Password` panic when misused; wrap them in helper functions if you need error returns instead.
- Combine `PasswordOptions{Symbols: true}` with `testutils.ContainsAllRunes` to assert generated strings hit every required character class in tests.

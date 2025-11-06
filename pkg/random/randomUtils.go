package random

import (
	cryptoRand "crypto/rand"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

const (
	alphaLowercase = "abcdefghijklmnopqrstuvwxyz"    // Lowercase Alpha
	alphaUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"    // Upercase Alpha
	alpha          = alphaLowercase + alphaUppercase // All alpha
	numeric        = "0123456789"
	symbols        = "!@#$%^&*()-_=+[]{}<>?/|~"
	alphaNumeric   = alpha + numeric // All Alpha + Numeric
)

// Num is a simple numeric constraint
type Num interface {
	~int | ~int64 | ~uint | ~float32 | ~float64
}

// Number returns a random number in the range [min, max] inclusive for int
// and [min, max) for floats
// Safely swaps max & min
// Limitations:
// - uses math/rand (not safe for concurrent use, not for security-sensitive use)
// - overflow warning on extreme ranges
func Number[T Num](min, max T) T {
	if min > max {
		min, max = max, min
	}
	switch any(min).(type) {
	case int:
		return min + T(rand.Intn(int(max-min+1)))
	case int64:
		return min + T(rand.Int63n(int64(max-min+1)))
	case uint:
		return min + T(uint(rand.Intn(int(max-min+1))))
	case float32:
		return min + T(rand.Float32()*float32(max-min))
	case float64:
		return min + T(rand.Float64()*float64(max-min))
	default:
		panic("unsupported type")
	}
}

// Sign returns either +1 or -1 randomly
func Sign() int {
	return 1 - 2*rand.Intn(2)
}

// Bool returns either true or false randomly
func Bool() bool {
	if rand.Intn(2) == 0 {
		return true
	}
	return false
}

// Choice picks a random element from a non-empty slice
// Panics if the slice is empty.
func Choice[T any](slice []T) T {
	if len(slice) == 0 {
		panic("Choice - empty slice")
	}
	return slice[rand.Intn(len(slice))]
}

// Byte returns a single secure random byte [0,255]
func Byte() byte {
	n, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(256))
	return byte(n.Int64())
}

// ----- Strings -----

// String returns a random alphanumeric string of length n
// includes both lowercase and uppercase letters
func String(n int) string {
	return FromCharset(n, alphaNumeric)
}

// Alpha returns a random string of letters of length n
// includes both lowercase and uppercase letters
func Alpha(n int) string {
	return FromCharset(n, alpha)
}

// Numeric returns a random string of digits of length n
func Numeric(n int) string {
	return FromCharset(n, numeric)
}

// FromCharset returns a random string from your specific charset
func FromCharset(n int, charset string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// PasswordOptions defines character sets for generating passwords
type PasswordOptions struct {
	Letters bool   // include a-zA-Z
	Digits  bool   // include 0-9
	Symbols bool   // include special characters
	Exclude string // any chars to exclude
}

// Password generates a random password of length n according to options.
// Panic if no charsets selected (or effectively not selected - e.g: excluding all options)
func Password(n int, opts PasswordOptions) string {
	if n <= 0 {
		return ""
	}

	charset := ""
	if opts.Letters {
		charset += alpha
	}
	if opts.Digits {
		charset += numeric
	}
	if opts.Symbols {
		charset += symbols
	}
	if opts.Exclude != "" {
		for _, exclude := range opts.Exclude {
			charset = strings.ReplaceAll(charset, string(exclude), "")
		}
	}

	if len(charset) == 0 {
		panic("random.password: no character sets selected")
	}

	password := make([]byte, n)
	for i := range password {
		idx, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[idx.Int64()]
	}
	return string(password)
}

// ----- Time / Dates -----

// Date returns a random time.Time between min and max (inclusive)
// Safely swaps min and max
func Date(min, max time.Time) time.Time {
	if min.After(max) {
		min, max = max, min
	}

	// Convert to Unix nanoseconds for integer math
	minUnix := min.UnixNano()
	maxUnix := max.UnixNano()

	delta := maxUnix - minUnix
	return time.Unix(0, minUnix+rand.Int63n(delta+1)).UTC()
}

// Duration returns a random time.Duration between min and max (inclusive)
// Safely swaps min and max
func Duration(min, max time.Duration) time.Duration {
	if min > max {
		min, max = max, min
	}

	delta := max - min
	return min + time.Duration(rand.Int63n(int64(delta)+1))
}

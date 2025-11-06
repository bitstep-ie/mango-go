package random

import (
	"github.com/bitstep-ie/mango-go/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"unicode"
)

func TestRandomNumber_Int(t *testing.T) {
	min, max := 5, 10
	for i := 0; i < 1000; i++ {
		n := Number(min, max)
		assert.GreaterOrEqual(t, n, min, "int: value below min")
		assert.LessOrEqual(t, n, max, "int: value above max")
	}
}

func TestRandomNumber_Int64(t *testing.T) {
	min, max := int64(100), int64(200)
	for i := 0; i < 1000; i++ {
		n := Number(min, max)
		assert.GreaterOrEqual(t, n, min, "int64: value below min")
		assert.LessOrEqual(t, n, max, "int64: value above max")
	}
}

func TestRandomNumber_Uint(t *testing.T) {
	min, max := uint(1), uint(100)
	for i := 0; i < 1000; i++ {
		n := Number(min, max)
		assert.GreaterOrEqual(t, n, min, "uint: value below min")
		assert.LessOrEqual(t, n, max, "uint: value above max")
	}
}

func TestRandomNumber_Float32(t *testing.T) {
	min, max := float32(1.5), float32(2.5)
	for i := 0; i < 1000; i++ {
		n := Number(min, max)
		assert.GreaterOrEqual(t, n, min, "float32: value below min")
		assert.LessOrEqual(t, n, max, "float32: value above max")
	}
}

func TestRandomNumber_Float64(t *testing.T) {
	min, max := 0.1, 0.2
	for i := 0; i < 1000; i++ {
		n := Number(min, max)
		assert.GreaterOrEqual(t, n, min, "float64: value below min")
		assert.LessOrEqual(t, n, max, "float64: value above max")
	}
}

func TestRandomNumber_MinEqualsMax(t *testing.T) {
	assert.Equal(t, 5, Number(5, 5), "int: min == max")
	assert.Equal(t, int64(42), Number(int64(42), int64(42)), "int64: min == max")
	assert.Equal(t, uint(7), Number(uint(7), uint(7)), "uint: min == max")
	assert.Equal(t, float32(3.14), Number(float32(3.14), float32(3.14)), "float32: min == max")
	assert.Equal(t, 2.718, Number(2.718, 2.718), "float64: min == max")
}

func TestRandomNumber_MinGreaterThanMax(t *testing.T) {
	n := Number(10, 5)
	assert.GreaterOrEqual(t, n, 5, "int: value below swapped min")
	assert.LessOrEqual(t, n, 10, "int: value above swapped max")
}

func TestRandomSign_ValidValues(t *testing.T) {
	for i := 0; i < 1000; i++ {
		val := Sign()
		assert.Contains(t, []int{-1, 1}, val, "Sign should only return -1 or 1")
	}
}

func TestRandomSign_Distribution(t *testing.T) {
	counts := map[int]int{-1: 0, 1: 0}
	for i := 0; i < 10000; i++ {
		val := Sign()
		counts[val]++
	}

	// Check both values are reasonably represented
	assert.Greater(t, counts[-1], 4000, "Expected at least ~40% -1s")
	assert.Greater(t, counts[1], 4000, "Expected at least ~40% 1s")
}

func TestRandomBool_ValidValues(t *testing.T) {
	for i := 0; i < 1000; i++ {
		val := Bool()
		assert.IsType(t, true, val, "Bool should return a boolean")
	}
}

func TestRandomBool_Distribution(t *testing.T) {
	counts := map[bool]int{true: 0, false: 0}
	for i := 0; i < 10000; i++ {
		val := Bool()
		counts[val]++
	}

	// Check both values are reasonably represented
	assert.Greater(t, counts[true], 4000, "Expected at least ~40% trues")
	assert.Greater(t, counts[false], 4000, "Expected at least ~40% falses")
}

func TestRandomChoice_ValidSelection(t *testing.T) {
	sample := []string{"apple", "banana", "cherry"}
	for i := 0; i < 1000; i++ {
		val := Choice(sample)
		assert.Contains(t, sample, val, "Choice should return an element from the slice")
	}
}

func TestRandomChoice_PanicOnEmptySlice(t *testing.T) {
	assert.Panics(t, func() {
		Choice([]string{})
	}, "Choice should panic on empty slice")
}

func TestRandomFromCharset_LengthAndCharset(t *testing.T) {
	charset := "abc123"
	for i := 0; i < 100; i++ {
		str := FromCharset(10, charset)
		assert.Len(t, str, 10, "FromCharset should return string of correct length")
		for _, ch := range str {
			assert.Contains(t, charset, string(ch), "Character should be from the charset")
		}
	}
}

func TestRandomString_AlphaNumericContent(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := String(20)
		assert.Len(t, str, 20)
		for _, ch := range str {
			assert.True(t, unicode.IsLetter(ch) || unicode.IsDigit(ch), "String should contain only alphanumeric characters")
		}
	}
}

func TestRandomAlpha_OnlyLetters(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := Alpha(15)
		assert.Len(t, str, 15)
		for _, ch := range str {
			assert.True(t, unicode.IsLetter(ch), "Alpha should contain only letters")
		}
	}
}

func TestRandomNumeric_OnlyDigits(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := Numeric(12)
		assert.Len(t, str, 12)
		for _, ch := range str {
			assert.True(t, unicode.IsDigit(ch), "Numeric should contain only digits")
		}
	}
}

func TestRandomFromCharset_Distribution(t *testing.T) {
	charset := "abc"
	counts := map[rune]int{'a': 0, 'b': 0, 'c': 0}
	total := 10000
	for i := 0; i < total; i++ {
		str := FromCharset(1, charset)
		counts[rune(str[0])]++
	}

	for _, ch := range charset {
		assert.Greater(t, counts[ch], total/5, "Each character should appear reasonably often")
	}
}

func TestDate_WithinBounds(t *testing.T) {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 1000; i++ {
		d := Date(min, max)
		assert.False(t, d.Before(min), "Date should not be before min")
		assert.False(t, d.After(max), "Date should not be after max")
	}
}

func TestDate_SwapsMinMax(t *testing.T) {
	min := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	max := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 1000; i++ {
		d := Date(min, max)
		assert.False(t, d.Before(max), "Date should not be before swapped min")
		assert.False(t, d.After(min), "Date should not be after swapped max")
	}
}

func TestDate_EqualMinMax(t *testing.T) {
	timestamp := time.Date(2021, 5, 5, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		d := Date(timestamp, timestamp)
		assert.Equal(t, timestamp, d, "Date should equal min/max when they are the same")
	}
}

func TestDuration_WithinBounds(t *testing.T) {
	min := 10 * time.Second
	max := 1 * time.Minute
	for i := 0; i < 1000; i++ {
		d := Duration(min, max)
		assert.GreaterOrEqual(t, d, min, "Duration should not be less than min")
		assert.LessOrEqual(t, d, max, "Duration should not be greater than max")
	}
}

func TestDuration_SwapsMinMax(t *testing.T) {
	min := 1 * time.Minute
	max := 10 * time.Second
	for i := 0; i < 1000; i++ {
		d := Duration(min, max)
		assert.GreaterOrEqual(t, d, max, "Duration should not be less than swapped min")
		assert.LessOrEqual(t, d, min, "Duration should not be greater than swapped max")
	}
}

func TestDuration_EqualMinMax(t *testing.T) {
	dur := 42 * time.Second
	for i := 0; i < 10; i++ {
		d := Duration(dur, dur)
		assert.Equal(t, dur, d, "Duration should equal min/max when they are the same")
	}
}

func TestByte_ValidRange(t *testing.T) {
	for i := 0; i < 1000; i++ {
		b := Byte()
		assert.GreaterOrEqual(t, int(b), 0, "Byte should be >= 0")
		assert.LessOrEqual(t, int(b), 255, "Byte should be <= 255")
	}
}

func TestPassword_Length(t *testing.T) {
	opts := PasswordOptions{Letters: true}
	pwd := Password(16, opts)
	assert.Len(t, pwd, 16, "Password should have correct length")
}

func TestPassword_LettersOnly(t *testing.T) {
	opts := PasswordOptions{Letters: true}
	pwd := Password(20, opts)
	for _, ch := range pwd {
		assert.True(t, unicode.IsLetter(ch), "Password should contain only letters")
	}
}

func TestPassword_DigitsOnly(t *testing.T) {
	opts := PasswordOptions{Digits: true}
	pwd := Password(20, opts)
	for _, ch := range pwd {
		assert.True(t, unicode.IsDigit(ch), "Password should contain only digits")
	}
}

func TestPassword_SymbolsOnly(t *testing.T) {
	opts := PasswordOptions{Symbols: true}
	pwd := Password(20, opts)
	for _, ch := range pwd {
		assert.Contains(t, symbols, string(ch), "Password should contain only symbols")
	}
}

func TestPassword_LettersAndDigits(t *testing.T) {
	opts := PasswordOptions{Letters: true, Digits: true}
	pwd := Password(30, opts)
	for _, ch := range pwd {
		assert.True(t, unicode.IsLetter(ch) || unicode.IsDigit(ch), "Password should contain letters or digits")
	}
}

func TestPassword_ExcludeCharacters(t *testing.T) {
	opts := PasswordOptions{Letters: true, Digits: true, Exclude: "aeiouAEIOU012345"}
	pwd := Password(50, opts)
	for _, ch := range pwd {
		assert.NotContains(t, opts.Exclude, string(ch), "Password should not contain excluded characters")
	}
}

func TestPassword_EmptyLength(t *testing.T) {
	opts := PasswordOptions{Letters: true}
	pwd := Password(0, opts)
	assert.Equal(t, "", pwd, "Password of length 0 should be empty")
}

func TestPassword_PanicOnEmptyCharset(t *testing.T) {
	opts := PasswordOptions{Letters: false, Digits: false, Symbols: false}
	assert.Panics(t, func() {
		Password(10, opts)
	}, "Password should panic if no character sets are selected")
}

func TestPassword_PanicOnExcludedAll(t *testing.T) {
	opts := PasswordOptions{
		Letters: true,
		Digits:  true,
		Symbols: true,
		Exclude: alpha + numeric + symbols,
	}
	assert.Panics(t, func() {
		Password(10, opts)
	}, "Password should panic if all characters are excluded")
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectedLen int
		contains    []string // optional: other strings it should include
	}{
		{"AlphaLowercase", alphaLowercase, 26, nil},
		{"AlphaUppercase", alphaUppercase, 26, nil},
		{"Alpha", alpha, 52, []string{alphaLowercase, alphaUppercase}},
		{"Numeric", numeric, 10, nil},
		{"Symbols", symbols, -1, nil}, // length unknown, just ensure not empty
		{"AlphaNumeric", alphaNumeric, 62, []string{alpha, numeric}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedLen > 0 {
				require.Equal(t, tt.expectedLen, len(tt.value), "length mismatch")
			} else {
				require.NotEmpty(t, tt.value, "value should not be empty")
			}

			for _, s := range tt.contains {
				testutils.ContainsAllRunes(t, tt.value, s, tt.name+" should contain all characters from sub-string")
			}
		})
	}
}

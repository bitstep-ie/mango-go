package env

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Deprecated alias smoke tests – ensure legacy methods still behave correctly
// ---------------------------------------------------------------------------

func TestDeprecated_EnvOrDefault(t *testing.T) {
	const key = "DEPRECATED_ENV_OR_DEFAULT"

	unsetEnv(t, key)
	assert.Equal(t, "fallback", EnvOrDefault(key, "fallback"))

	setEnv(t, key, "real")
	assert.Equal(t, "real", EnvOrDefault(key, "fallback"))
}

func TestDeprecated_MustEnv(t *testing.T) {
	const key = "DEPRECATED_MUST_ENV"

	setEnv(t, key, "present")
	assert.Equal(t, "present", MustEnv(key))

	unsetEnv(t, key)
	assert.Panics(t, func() { MustEnv(key) })
}

func TestDeprecated_EnvAsInt(t *testing.T) {
	const key = "DEPRECATED_ENV_AS_INT"

	unsetEnv(t, key)
	assert.Equal(t, 7, EnvAsInt(key, 7))

	setEnv(t, key, "99")
	assert.Equal(t, 99, EnvAsInt(key, 7))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { EnvAsInt(key, 7) })
}

func TestDeprecated_MustEnvAsInt(t *testing.T) {
	const key = "DEPRECATED_MUST_ENV_AS_INT"

	setEnv(t, key, "42")
	assert.Equal(t, 42, MustEnvAsInt(key))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { MustEnvAsInt(key) })

	unsetEnv(t, key)
	assert.Panics(t, func() { MustEnvAsInt(key) })
}

func TestDeprecated_EnvAsBool(t *testing.T) {
	const key = "DEPRECATED_ENV_AS_BOOL"

	unsetEnv(t, key)
	assert.True(t, EnvAsBool(key, true))

	setEnv(t, key, "false")
	assert.False(t, EnvAsBool(key, true))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { EnvAsBool(key, true) })
}

func TestDeprecated_MustEnvAsBool(t *testing.T) {
	const key = "DEPRECATED_MUST_ENV_AS_BOOL"

	setEnv(t, key, "true")
	assert.True(t, MustEnvAsBool(key))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { MustEnvAsBool(key) })

	unsetEnv(t, key)
	assert.Panics(t, func() { MustEnvAsBool(key) })
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func TestStringDefault(t *testing.T) {
	const key = "STRING_DEFAULT"

	unsetEnv(t, key)
	v, err := StringDefault(key, "fallback")
	assert.NoError(t, err)
	assert.Equal(t, "fallback", v)

	setEnv(t, key, "hello")
	v, err = StringDefault(key, "fallback")
	assert.NoError(t, err)
	assert.Equal(t, "hello", v)
}

func TestString(t *testing.T) {
	const key = "STRING"

	setEnv(t, key, "world")
	v, err := String(key)
	assert.NoError(t, err)
	assert.Equal(t, "world", v)

	unsetEnv(t, key)
	_, err = String(key)
	assert.Error(t, err)
}

func TestMustString(t *testing.T) {
	const key = "MUST_STRING"

	setEnv(t, key, "required")
	assert.Equal(t, "required", MustString(key))

	unsetEnv(t, key)
	assert.Panics(t, func() { MustString(key) })
}

// ---------------------------------------------------------------------------
// Int
// ---------------------------------------------------------------------------

func TestIntDefault(t *testing.T) {
	const key = "INT_DEFAULT"

	unsetEnv(t, key)
	v, err := IntDefault(key, 10)
	assert.NoError(t, err)
	assert.Equal(t, 10, v)

	setEnv(t, key, "55")
	v, err = IntDefault(key, 10)
	assert.NoError(t, err)
	assert.Equal(t, 55, v)

	setEnv(t, key, "not-an-int")
	_, err = IntDefault(key, 10)
	assert.Error(t, err)
}

func TestInt(t *testing.T) {
	const key = "INT"

	setEnv(t, key, "123")
	v, err := Int(key)
	assert.NoError(t, err)
	assert.Equal(t, 123, v)

	setEnv(t, key, "bad")
	_, err = Int(key)
	assert.Error(t, err)

	unsetEnv(t, key)
	_, err = Int(key)
	assert.Error(t, err)
}

func TestMustInt(t *testing.T) {
	const key = "MUST_INT"

	setEnv(t, key, "7")
	assert.Equal(t, 7, MustInt(key))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { MustInt(key) })

	unsetEnv(t, key)
	assert.Panics(t, func() { MustInt(key) })
}

// ---------------------------------------------------------------------------
// Bool
// ---------------------------------------------------------------------------

func TestBoolDefault(t *testing.T) {
	const key = "BOOL_DEFAULT"

	unsetEnv(t, key)
	v, err := BoolDefault(key, true)
	assert.NoError(t, err)
	assert.True(t, v)

	setEnv(t, key, "false")
	v, err = BoolDefault(key, true)
	assert.NoError(t, err)
	assert.False(t, v)

	// strconv.ParseBool accepts "1" as true
	setEnv(t, key, "1")
	v, err = BoolDefault(key, false)
	assert.NoError(t, err)
	assert.True(t, v)

	setEnv(t, key, "not-a-bool")
	_, err = BoolDefault(key, false)
	assert.Error(t, err)
}

func TestBool(t *testing.T) {
	const key = "BOOL"

	setEnv(t, key, "true")
	v, err := Bool(key)
	assert.NoError(t, err)
	assert.True(t, v)

	setEnv(t, key, "bad")
	_, err = Bool(key)
	assert.Error(t, err)

	unsetEnv(t, key)
	_, err = Bool(key)
	assert.Error(t, err)
}

func TestMustBool(t *testing.T) {
	const key = "MUST_BOOL"

	setEnv(t, key, "false")
	assert.False(t, MustBool(key))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { MustBool(key) })

	unsetEnv(t, key)
	assert.Panics(t, func() { MustBool(key) })
}

// ---------------------------------------------------------------------------
// Duration
// ---------------------------------------------------------------------------

func TestDurationDefault(t *testing.T) {
	const key = "DURATION_DEFAULT"

	unsetEnv(t, key)
	v, err := DurationDefault(key, 5*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 5*time.Second, v)

	setEnv(t, key, "2m")
	v, err = DurationDefault(key, 5*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Minute, v)

	setEnv(t, key, "bad-duration")
	_, err = DurationDefault(key, 5*time.Second)
	assert.Error(t, err)
}

func TestDuration(t *testing.T) {
	const key = "DURATION"

	setEnv(t, key, "90m")
	v, err := Duration(key)
	assert.NoError(t, err)
	assert.Equal(t, 90*time.Minute, v)

	setEnv(t, key, "bad")
	_, err = Duration(key)
	assert.Error(t, err)

	unsetEnv(t, key)
	_, err = Duration(key)
	assert.Error(t, err)
}

func TestMustDuration(t *testing.T) {
	const key = "MUST_DURATION"

	setEnv(t, key, "1h")
	assert.Equal(t, time.Hour, MustDuration(key))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { MustDuration(key) })

	unsetEnv(t, key)
	assert.Panics(t, func() { MustDuration(key) })
}

// ---------------------------------------------------------------------------
// StringSplice
// ---------------------------------------------------------------------------

func TestStringSpliceDefault(t *testing.T) {
	const key = "STRING_SPLICE_DEFAULT"

	unsetEnv(t, key)
	v, err := StringSpliceDefault(key, ",", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, v)

	setEnv(t, key, "x,y,z")
	v, err = StringSpliceDefault(key, ",", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"x", "y", "z"}, v)

	// custom separator
	setEnv(t, key, "alpha|beta|gamma")
	v, err = StringSpliceDefault(key, "|", []string{"default"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"alpha", "beta", "gamma"}, v)
}

func TestStringSplice(t *testing.T) {
	const key = "STRING_SPLICE"

	setEnv(t, key, "one,two,three")
	v, err := StringSplice(key, ",")
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two", "three"}, v)

	unsetEnv(t, key)
	_, err = StringSplice(key, ",")
	assert.Error(t, err)
}

func TestMustStringSplice(t *testing.T) {
	const key = "MUST_STRING_SPLICE"

	setEnv(t, key, "red,green,blue")
	assert.Equal(t, []string{"red", "green", "blue"}, MustStringSplice(key, ","))

	unsetEnv(t, key)
	assert.Panics(t, func() { MustStringSplice(key, ",") })
}

// ---------------------------------------------------------------------------
// Generic: AsDefault / As / MustAs
// ---------------------------------------------------------------------------

func TestAsDefault(t *testing.T) {
	const key = "AS_DEFAULT"

	unsetEnv(t, key)
	v, err := AsDefault(key, 99, strconv.Atoi)
	assert.NoError(t, err)
	assert.Equal(t, 99, v)

	setEnv(t, key, "42")
	v, err = AsDefault(key, 99, strconv.Atoi)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	setEnv(t, key, "bad")
	_, err = AsDefault(key, 99, strconv.Atoi)
	assert.Error(t, err)
}

func TestAs(t *testing.T) {
	const key = "AS"

	setEnv(t, key, "7")
	v, err := As(key, strconv.Atoi)
	assert.NoError(t, err)
	assert.Equal(t, 7, v)

	setEnv(t, key, "bad")
	_, err = As(key, strconv.Atoi)
	assert.Error(t, err)

	unsetEnv(t, key)
	_, err = As(key, strconv.Atoi)
	assert.Error(t, err)
}

func TestMustAs(t *testing.T) {
	const key = "MUST_AS"

	setEnv(t, key, "24")
	assert.Equal(t, 24, MustAs(key, strconv.Atoi))

	setEnv(t, key, "bad")
	assert.Panics(t, func() { MustAs(key, strconv.Atoi) })

	unsetEnv(t, key)
	assert.Panics(t, func() { MustAs(key, strconv.Atoi) })
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	assert.NoError(t, os.Unsetenv(key))
	t.Cleanup(func() {
		_ = os.Unsetenv(key)
	})
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	assert.NoError(t, os.Setenv(key, value))
	t.Cleanup(func() {
		_ = os.Unsetenv(key)
	})
}

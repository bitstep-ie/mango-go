// Package env has simple functions for environment variables manipulation
package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const envVarNotSetMessage = "environment variable %s not set"

// Deprecated: use StringDefault instead. Will be removed in v1.0.0
// EnvOrDefault will return
// the value of the environment variable key (envKey),
// or if not found (eq "") returns the default (defaultVal)
func EnvOrDefault(envKey, defaultVal string) string {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal
	}
	return value
}

// Deprecated: use String or MustString instead. Will be removed in v1.0.0
// MustEnv will return
// the value of the environment variable key (envKey)
// or panic if this is not found (eq "")
func MustEnv(envKey string) string {
	value := os.Getenv(envKey)
	if value == "" {
		panic(fmt.Sprintf(envVarNotSetMessage, envKey))
	}
	return value
}

// Deprecated: use IntDefault instead. Will be removed in v1.0.0
// EnvAsInt returns the envKey as int
// or the default value if missing (eq "")
// Panics if the value of the envKey is NOT an integer
func EnvAsInt(envKey string, defaultVal int) int {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not an integer", envKey))
	}
	return i
}

// Deprecated: use Int or MustInt instead. Will be removed in v1.0.0
// MustEnvAsInt returns the envKey as int
// or panics if missing (eq "")
// Panics if the value of the envKey is NOT an integer
func MustEnvAsInt(envKey string) int {
	value := os.Getenv(envKey)
	if value == "" {
		panic(fmt.Sprintf(envVarNotSetMessage, envKey))
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not an integer", envKey))
	}
	return i
}

// Deprecated: use BoolDefault instead. Will be removed in v1.0.0
// EnvAsBool returns the envKey as bool
// or the default value if missing (eq "")
// Panics if the value of the envKey is NOT a bool (parsing with strconv.ParseBool)
func EnvAsBool(envKey string, defaultVal bool) bool {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a boolean", envKey))
	}
	return b
}

// Deprecated: use Bool or MustBool instead. Will be removed in v1.0.0
// MustEnvAsBool returns the envKey as bool
// Panics if envKey is NOT set
// Panics if the value of the envKey is NOT a bool (parsing with strconv.ParseBool)
func MustEnvAsBool(envKey string) bool {
	value := os.Getenv(envKey)
	if value == "" {
		panic(fmt.Sprintf(envVarNotSetMessage, envKey))
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a boolean", envKey))
	}
	return b
}

// StringDefault returns the value of the environment variable envKey as a string.
// If the variable is not set (eq ""), it returns defaultVal and a nil error.
// A non-empty value is always returned without error regardless of its content.
func StringDefault(envKey string, defaultVal string) (string, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal, nil
	}
	return value, nil
}

// String returns the value of the environment variable envKey as a string.
// Returns an error if the variable is not set (eq "").
func String(envKey string) (string, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return "", fmt.Errorf(envVarNotSetMessage, envKey)
	}
	return value, nil
}

// MustString returns the value of the environment variable envKey as a string.
// Panics if the variable is not set (eq "").
func MustString(envKey string) string {
	value, err := String(envKey)
	if err != nil {
		panic(err)
	}
	return value
}

// IntDefault returns the value of the environment variable envKey as an int.
// If the variable is not set (eq ""), it returns defaultVal and a nil error.
// Returns an error if the value cannot be parsed as an integer.
func IntDefault(envKey string, defaultVal int) (int, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal, nil
	}
	parsedEnvValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not an integer", envKey)
	}
	return parsedEnvValue, nil
}

// Int returns the value of the environment variable envKey as an int.
// Returns an error if the variable is not set (eq "") or cannot be parsed as an integer.
func Int(envKey string) (int, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return 0, fmt.Errorf(envVarNotSetMessage, envKey)
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not an integer", envKey)
	}
	return parsedValue, nil
}

// MustInt returns the value of the environment variable envKey as an int.
// Panics if the variable is not set (eq "") or cannot be parsed as an integer.
func MustInt(envKey string) int {
	value, err := Int(envKey)
	if err != nil {
		panic(err)
	}
	return value
}

// BoolDefault returns the value of the environment variable envKey as a bool.
// If the variable is not set (eq ""), it returns defaultVal and a nil error.
// Parsing is delegated to strconv.ParseBool, which accepts: 1, t, T, TRUE, true, True,
// 0, f, F, FALSE, false, False. Returns an error for any other value.
func BoolDefault(envKey string, defaultVal bool) (bool, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal, nil
	}
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s is not a boolean", envKey)
	}
	return parsedValue, nil
}

// Bool returns the value of the environment variable envKey as a bool.
// Returns an error if the variable is not set (eq "") or cannot be parsed as a bool.
// Parsing is delegated to strconv.ParseBool.
func Bool(envKey string) (bool, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return false, fmt.Errorf(envVarNotSetMessage, envKey)
	}
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s is not a boolean", envKey)
	}
	return parsedValue, nil
}

// MustBool returns the value of the environment variable envKey as a bool.
// Panics if the variable is not set (eq "") or cannot be parsed as a bool.
// Parsing is delegated to strconv.ParseBool.
func MustBool(envKey string) bool {
	value, err := Bool(envKey)
	if err != nil {
		panic(err)
	}
	return value
}

// DurationDefault returns the value of the environment variable envKey as a time.Duration.
// If the variable is not set (eq ""), it returns defaultVal and a nil error.
// Parsing is delegated to time.ParseDuration. Valid units: ns, us (or µs), ms, s, m, h.
// Returns an error if the value cannot be parsed as a duration.
func DurationDefault(envKey string, defaultVal time.Duration) (time.Duration, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal, nil
	}
	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a duration", envKey)
	}
	return parsedValue, nil
}

// Duration returns the value of the environment variable envKey as a time.Duration.
// Returns an error if the variable is not set (eq "") or cannot be parsed as a duration.
// Parsing is delegated to time.ParseDuration. Valid units: ns, us (or µs), ms, s, m, h.
func Duration(envKey string) (time.Duration, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return 0, fmt.Errorf(envVarNotSetMessage, envKey)
	}
	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a duration", envKey)
	}
	return parsedValue, nil
}

// MustDuration returns the value of the environment variable envKey as a time.Duration.
// Panics if the variable is not set (eq "") or cannot be parsed as a duration.
// Parsing is delegated to time.ParseDuration. Valid units: ns, us (or µs), ms, s, m, h.
func MustDuration(envKey string) time.Duration {
	value, err := Duration(envKey)
	if err != nil {
		panic(err)
	}
	return value
}

// StringSpliceDefault returns the value of the environment variable envKey split
// into a []string using sep as the delimiter.
// If the variable is not set (eq ""), it returns defaultVal and a nil error.
// No trimming or filtering of elements is applied; the raw strings.Split result is returned.
func StringSpliceDefault(envKey string, sep string, defaultVal []string) ([]string, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal, nil
	}
	return strings.Split(value, sep), nil
}

// StringSplice returns the value of the environment variable envKey split
// into a []string using sep as the delimiter.
// Returns an error if the variable is not set (eq "").
// No trimming or filtering of elements is applied; the raw strings.Split result is returned.
func StringSplice(envKey string, sep string) ([]string, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return nil, fmt.Errorf(envVarNotSetMessage, envKey)
	}
	return strings.Split(value, sep), nil
}

// MustStringSplice returns the value of the environment variable envKey split
// into a []string using sep as the delimiter.
// Panics if the variable is not set (eq "").
// No trimming or filtering of elements is applied; the raw strings.Split result is returned.
func MustStringSplice(envKey string, sep string) []string {
	value, err := StringSplice(envKey, sep)
	if err != nil {
		panic(err)
	}
	return value
}

// AsDefault returns the value of the environment variable envKey parsed as T using parser.
// If the variable is not set (eq ""), it returns defaultVal and a nil error.
// Returns an error if parser returns an error for the raw string value.
// parser must be a function of the form func(string) (T, error).
func AsDefault[T any](envKey string, defaultVal T, parser func(string) (T, error)) (T, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultVal, nil
	}

	parsedValue, err := parser(value)
	if err != nil {
		return *new(T), fmt.Errorf("environment variable %s could not be parsed: %v+", envKey, err)
	}
	return parsedValue, nil
}

// As returns the value of the environment variable envKey parsed as T using parser.
// Returns an error if the variable is not set (eq "") or parser returns an error.
// parser must be a function of the form func(string) (T, error).
func As[T any](envKey string, parser func(string) (T, error)) (T, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return *new(T), fmt.Errorf(envVarNotSetMessage, envKey)
	}

	parsedValue, err := parser(value)
	if err != nil {
		return *new(T), fmt.Errorf("environment variable %s could not be parsed: %v+", envKey, err)
	}
	return parsedValue, nil
}

// MustAs returns the value of the environment variable envKey parsed as T using parser.
// Panics if the variable is not set (eq "") or parser returns an error.
// parser must be a function of the form func(string) (T, error).
func MustAs[T any](envKey string, parser func(string) (T, error)) T {
	value, err := As[T](envKey, parser)
	if err != nil {
		panic(err)
	}
	return value
}

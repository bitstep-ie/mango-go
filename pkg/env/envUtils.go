// Package env has simple functions for environment variables manipulation
package env

import (
	"fmt"
	"os"
	"strconv"
)

const envVarNotSetMessage = "environment variable %s not set"

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

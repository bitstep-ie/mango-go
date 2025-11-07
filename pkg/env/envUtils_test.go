package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnvOrDefault(t *testing.T) {
	const key = "ENV_OR_DEFAULT"
	const defaultVal = "default"

	err := os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := EnvOrDefault(key, defaultVal); got != defaultVal {
		t.Errorf("Expected default value %q, got %q", defaultVal, got)
	}

	expected := "actual"
	err = os.Setenv(key, expected)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := EnvOrDefault(key, defaultVal); got != expected {
		t.Errorf("Expected env value %q, got %q", expected, got)
	}
	err = os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
}

func TestMustEnv(t *testing.T) {
	const key = "MUST_ENV"
	expected := "required"

	err := os.Setenv(key, expected)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := MustEnv(key); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
	err = os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when env var is missing")
		}
	}()
	MustEnv(key)
}

func TestMustEnv_PanicWhenMissing(t *testing.T) {
	const key = "MUST_ENV_MISSING"
	err := os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when environment variable %q is missing", key)
		}
	}()
	MustEnv(key)
}

func TestEnvAsInt(t *testing.T) {
	const key = "ENV_AS_INT"
	const defaultVal = 42

	err := os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := EnvAsInt(key, defaultVal); got != defaultVal {
		t.Errorf("Expected default int %d, got %d", defaultVal, got)
	}

	err = os.Setenv(key, "100")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := EnvAsInt(key, defaultVal); got != 100 {
		t.Errorf("Expected 100, got %d", got)
	}

	err = os.Setenv(key, "not_an_int")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-integer value")
		}
	}()
	EnvAsInt(key, defaultVal)
	err = os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
}

func TestMustEnvAsInt(t *testing.T) {
	const key = "MUST_ENV_AS_INT"

	err := os.Setenv(key, "200")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := MustEnvAsInt(key); got != 200 {
		t.Errorf("Expected 200, got %d", got)
	}

	err = os.Setenv(key, "not_an_int")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-integer value")
		}
	}()
	MustEnvAsInt(key)

	err = os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing env var")
		}
	}()
	MustEnvAsInt(key)
}

func TestMustEnvAsInt_PanicWhenMissing(t *testing.T) {
	const key = "MUST_ENV_AS_INT_MISSING"
	err := os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when environment variable %q is missing", key)
		}
	}()
	MustEnvAsInt(key)
}

func TestEnvAsBool(t *testing.T) {
	const key = "ENV_AS_BOOL"
	const defaultVal = true

	err := os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := EnvAsBool(key, defaultVal); got != defaultVal {
		t.Errorf("Expected default bool %v, got %v", defaultVal, got)
	}

	err = os.Setenv(key, "false")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := EnvAsBool(key, defaultVal); got != false {
		t.Errorf("Expected false, got %v", got)
	}

	err = os.Setenv(key, "not_a_bool")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-boolean value")
		}
	}()
	EnvAsBool(key, defaultVal)
	err = os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
}

func TestMustEnvAsBool(t *testing.T) {
	const key = "MUST_ENV_AS_BOOL"

	err := os.Setenv(key, "true")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	if got := MustEnvAsBool(key); got != true {
		t.Errorf("Expected true, got %v", got)
	}

	err = os.Setenv(key, "not_a_bool")
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-boolean value")
		}
	}()
	MustEnvAsBool(key)

	err = os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing env var")
		}
	}()
	MustEnvAsBool(key)
}

func TestMustEnvAsBool_PanicWhenMissing(t *testing.T) {
	const key = "MUST_ENV_AS_BOOL_MISSING"
	err := os.Unsetenv(key)
	if err != nil {
		assert.NoError(t, err, "No error expected")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when environment variable %q is missing", key)
		}
	}()
	MustEnvAsBool(key)
}

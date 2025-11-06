package env

import (
	"os"
	"testing"
)

func TestEnvOrDefault(t *testing.T) {
	const key = "ENV_OR_DEFAULT"
	const defaultVal = "default"

	os.Unsetenv(key)
	if got := EnvOrDefault(key, defaultVal); got != defaultVal {
		t.Errorf("Expected default value %q, got %q", defaultVal, got)
	}

	expected := "actual"
	os.Setenv(key, expected)
	if got := EnvOrDefault(key, defaultVal); got != expected {
		t.Errorf("Expected env value %q, got %q", expected, got)
	}
	os.Unsetenv(key)
}

func TestMustEnv(t *testing.T) {
	const key = "MUST_ENV"
	expected := "required"

	os.Setenv(key, expected)
	if got := MustEnv(key); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
	os.Unsetenv(key)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when env var is missing")
		}
	}()
	MustEnv(key)
}

func TestMustEnv_PanicWhenMissing(t *testing.T) {
	const key = "MUST_ENV_MISSING"
	os.Unsetenv(key)

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

	os.Unsetenv(key)
	if got := EnvAsInt(key, defaultVal); got != defaultVal {
		t.Errorf("Expected default int %d, got %d", defaultVal, got)
	}

	os.Setenv(key, "100")
	if got := EnvAsInt(key, defaultVal); got != 100 {
		t.Errorf("Expected 100, got %d", got)
	}

	os.Setenv(key, "not_an_int")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-integer value")
		}
	}()
	EnvAsInt(key, defaultVal)
	os.Unsetenv(key)
}

func TestMustEnvAsInt(t *testing.T) {
	const key = "MUST_ENV_AS_INT"

	os.Setenv(key, "200")
	if got := MustEnvAsInt(key); got != 200 {
		t.Errorf("Expected 200, got %d", got)
	}

	os.Setenv(key, "not_an_int")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-integer value")
		}
	}()
	MustEnvAsInt(key)

	os.Unsetenv(key)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing env var")
		}
	}()
	MustEnvAsInt(key)
}

func TestMustEnvAsInt_PanicWhenMissing(t *testing.T) {
	const key = "MUST_ENV_AS_INT_MISSING"
	os.Unsetenv(key)

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

	os.Unsetenv(key)
	if got := EnvAsBool(key, defaultVal); got != defaultVal {
		t.Errorf("Expected default bool %v, got %v", defaultVal, got)
	}

	os.Setenv(key, "false")
	if got := EnvAsBool(key, defaultVal); got != false {
		t.Errorf("Expected false, got %v", got)
	}

	os.Setenv(key, "not_a_bool")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-boolean value")
		}
	}()
	EnvAsBool(key, defaultVal)
	os.Unsetenv(key)
}

func TestMustEnvAsBool(t *testing.T) {
	const key = "MUST_ENV_AS_BOOL"

	os.Setenv(key, "true")
	if got := MustEnvAsBool(key); got != true {
		t.Errorf("Expected true, got %v", got)
	}

	os.Setenv(key, "not_a_bool")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-boolean value")
		}
	}()
	MustEnvAsBool(key)

	os.Unsetenv(key)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing env var")
		}
	}()
	MustEnvAsBool(key)
}

func TestMustEnvAsBool_PanicWhenMissing(t *testing.T) {
	const key = "MUST_ENV_AS_BOOL_MISSING"
	os.Unsetenv(key)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when environment variable %q is missing", key)
		}
	}()
	MustEnvAsBool(key)
}

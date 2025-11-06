package testutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustMakeTempFile(t *testing.T) {
	tempDir := t.TempDir()

	file := MustMakeTempFile(t, tempDir, "testFile-*.txt")

	assert.FileExists(t, file.Name())
}

func TestMustMakeTempFileFailing(t *testing.T) {
	assertionFailed := checkForAssertion(func(t *testing.T) {
		_ = MustMakeTempFile(t, "/root/", "testFile-*.txt")
	})

	if !assertionFailed {
		t.Errorf("Expected the test function to assert, but it did not.")
	}
}

func TestAssertInValidUUID(t *testing.T) {
	assertionFailed := checkForAssertion(func(t *testing.T) {
		AssertValidUUID(t, "invalid", "fieldName")
	})

	if !assertionFailed {
		t.Errorf("Expected the test function to assert, but it did not.")
	}
}

func TestAssertValidUUID(t *testing.T) {
	assertionFailed := checkForAssertion(func(t *testing.T) {
		AssertValidUUID(t, "f7198919-bb3e-4ff0-9746-72d913f8a812", "fieldName")
	})

	if assertionFailed {
		t.Errorf("Expected the test function to NOT assert, but it did.")
	}
}

func checkForAssertion(testFunc func(t *testing.T)) (asserted bool) {
	defer func() {
		if r := recover(); r != nil { // Recover from panic caused by FailNow() or Fatal()
			asserted = true
		}
	}()

	t := &testing.T{}
	testFunc(t) // Run the test function in a subtest to isolate its behavior
	return t.Failed()
}

func TestContainsAllRunes(t *testing.T) {
	str := "abcdef"
	chars := "ace"
	// Should pass, no test failures
	ContainsAllRunes(t, str, chars)
	// You can assert no failures were recorded
	assert.False(t, t.Failed(), "expected no failures")
}

// Package testutils contains helping methods to be used in tests
package testutils

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// MustMakeTempFile creates a temporary file in the tempDir using the fileNamePattern
// Failure to create the temporary file will fail the test using assert.Fail
// Returns the file created
func MustMakeTempFile(t *testing.T, tempDir string, fileNamePattern string) *os.File {
	textFile, err := os.CreateTemp(tempDir, fileNamePattern)
	if err != nil {
		assert.Fail(t, "Failed to create temp file: %v", err)
	}
	defer func(textFile *os.File) {
		_ = textFile.Close()
	}(textFile)
	return textFile
}

// AssertValidUUID asserts if the value is a valid UUID
// value is the value to be checked
// fieldName is the fieldName printed in the error message if the value is NOT a valid UUID
func AssertValidUUID(t *testing.T, value string, fieldName string) {
	_, err := uuid.Parse(value)
	assert.NoError(t, err, "The "+fieldName+" should be a valid UUID")
}

// ContainsAllRunes helper to assert all runes in `chars` exist in `str`
func ContainsAllRunes(t *testing.T, str string, chars string, msgAndArgs ...interface{}) {
	for _, c := range chars {
		assert.Contains(t, str, string(c), msgAndArgs...)
	}
}

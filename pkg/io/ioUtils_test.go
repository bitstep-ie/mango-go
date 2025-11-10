package io

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bitstep-ie/mango-go/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// mockCloser simulates an io.Closer with a controllable error response
type mockCloser struct {
	closed bool
	err    error
}

func (m *mockCloser) Close() error {
	m.closed = true
	return m.err
}

func TestSafeClose(t *testing.T) {
	t.Run("closes a normal closer", func(t *testing.T) {
		c := &mockCloser{}
		SafeClose(c)
		assert.True(t, c.closed, "expected Close to be called")
	})

	t.Run("handles nil without panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			SafeClose(nil)
		}, "SafeClose(nil) should not panic")
	})

	t.Run("ignores error from closer", func(t *testing.T) {
		c := &mockCloser{err: errors.New("boom")}
		assert.NotPanics(t, func() {
			SafeClose(c)
		}, "SafeClose should ignore errors from Close")
		assert.True(t, c.closed, "expected Close to be called even if it errors")
	})
}

func TestSafeClosePrint(t *testing.T) {
	t.Run("should close without error", func(t *testing.T) {
		// capture log output
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(io.Discard)

		c := &mockCloser{err: nil}
		SafeClosePrint(c)

		assert.Empty(t, buf.String(), "expected no log output on successful close")
	})

	t.Run("should log error if close fails", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(io.Discard)

		c := &mockCloser{err: errors.New("boom")}
		SafeClosePrint(c)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "close error: boom", "expected error to be logged")
	})
}

func TestBackupFilesMatching(t *testing.T) {

	tempDir := t.TempDir()

	textFile := testutils.MustMakeTempFile(t, tempDir, "tempfile-*.txt")
	noBackupFile := testutils.MustMakeTempFile(t, tempDir, "tempfile-*.tmp")

	path, _ := filepath.Abs(filepath.Dir(textFile.Name()))

	errs := BackupFilesWithExt(path, []string{".txt"})
	if errs != nil {
		t.Errorf("Failed to backup the file with ext %s", errs)
	}

	assert.FileExists(t, noBackupFile.Name(), "File NOT for backup was moved/deleted, expected to still exist")
	assert.FileExists(t, textFile.Name(), "Expected to be backedup to continue to exist")
	assert.FileExists(t, textFile.Name()+".bak", "Backup file of matching is created")

}

func TestBackupFilesMatchingFailToWalk(t *testing.T) {
	// Save the original function to restore later
	originalWalkFunc := walkFunc

	// Override walkFunc to simulate failure
	walkFunc = func(root string, fn filepath.WalkFunc) error {
		return fmt.Errorf("failed to walk")
	}

	tempDir := t.TempDir()

	err := BackupFilesWithExt(tempDir, []string{".txt"})
	if err != nil {
		assert.Equal(t, "failed to walk", err.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	walkFunc = originalWalkFunc
}

func TestBackupFilesMatchingFailToCopyFiles(t *testing.T) {
	// Save the original function to restore later
	originalCopyFunc := copyFunc

	// Override walkFunc to simulate failure
	copyFunc = func(dst io.Writer, src io.Reader) (written int64, err error) {
		return 0, fmt.Errorf("failed to copy")
	}

	tempDir := t.TempDir()

	textFile := testutils.MustMakeTempFile(t, tempDir, "tempfile-*.txt")

	path, _ := filepath.Abs(filepath.Dir(textFile.Name()))

	err := BackupFilesWithExt(path, []string{".txt"})
	if err != nil {
		assert.Equal(t, "failed to copy", err.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	copyFunc = originalCopyFunc
}

func TestDeleteFileWithExt(t *testing.T) {
	tempDir := t.TempDir()

	file := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt")
	fileNoMatchingExt := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt.bak")

	path, _ := filepath.Abs(filepath.Dir(file.Name()))

	errs := DeleteFileWithExt(path, []string{".txt"})
	if errs != nil {
		t.Errorf("Failed to delete the file with ext")
	}

	assert.NoFileExists(t, file.Name(), "Expected to be deleted test file was NOT deleted")
	assert.FileExists(t, fileNoMatchingExt.Name(), "Expected to NOT be deleted test file WAS deleted")
}

func TestDeleteFileWithExtFailsToDelete(t *testing.T) {
	tempDir := t.TempDir()

	file := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt")

	path, _ := filepath.Abs(filepath.Dir(file.Name()))

	// Save the original function to restore later
	originalRemoveFunc := removeFunc

	// Override removeFunc to simulate failure
	removeFunc = func(name string) error {
		return fmt.Errorf("failed to remove file")
	}

	errs := DeleteFileWithExt(path, []string{".txt"})
	if errs != nil {
		assert.Equal(t, "failed to delete file "+file.Name()+": failed to remove file", errs.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	// Restore the original remove function
	removeFunc = originalRemoveFunc
}

func TestDeleteFileWithExtFailsToWalk(t *testing.T) {
	// Save the original function to restore later
	originalWalkFunc := walkFunc

	// Override walkFunc to simulate failure
	walkFunc = func(root string, fn filepath.WalkFunc) error {
		return fmt.Errorf("failed to walk")
	}

	errs := DeleteFileWithExt("/", []string{".txt"})
	if errs != nil {
		assert.Equal(t, "failed to walk", errs.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	// Restore the original remove function
	walkFunc = originalWalkFunc
}

func TestRestoreAllBakFiles(t *testing.T) {
	tempDir := t.TempDir()

	backedFile := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt.bak")
	nonBackedFile := testutils.MustMakeTempFile(t, tempDir, "testFile-*.tmp")

	path, _ := filepath.Abs(filepath.Dir(backedFile.Name()))

	errs := RestoreAllBakFiles(path)
	if errs != nil {
		t.Errorf("Failed to restore the bak file")
	}

	assert.FileExists(t, nonBackedFile.Name(), "File was NOT restored")
	assert.NoFileExists(t, backedFile.Name(), "No bak file exists anymore")
	assert.FileExists(t, strings.Trim(backedFile.Name(), ".bak"), "Restored file exists")
}

func TestRestoreAllBakFilesWhenExistingOriginalFileWithDifferentContent(t *testing.T) {
	tempDir := t.TempDir()

	backedFile := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt.bak")
	err := os.WriteFile(backedFile.Name(), []byte("my backedup data"), 0644)
	if err != nil {
		t.Fatalf("Failed to write the backedup file")
	}
	originalFilename := strings.Replace(backedFile.Name(), ".bak", "", 1)
	err = os.WriteFile(originalFilename, []byte("my original data"), 0644)
	if err != nil {
		t.Fatalf("Error writing file: %s", err)
	}

	nonBackedFile := testutils.MustMakeTempFile(t, tempDir, "testFile-*.tmp")

	path, _ := filepath.Abs(filepath.Dir(backedFile.Name()))

	errs := RestoreAllBakFiles(path)
	if errs != nil {
		t.Errorf("Failed to restore the bak file")
	}

	assert.FileExists(t, nonBackedFile.Name(), "File was NOT restored")
	assert.NoFileExists(t, backedFile.Name(), "No bak file exists anymore")
	assert.FileExists(t, strings.Trim(backedFile.Name(), ".bak"), "Restored file exists")
	originalContent, err := os.ReadFile(originalFilename)
	if err != nil {
		t.Fatalf("Failed to read the original file")
	}
	assert.Equal(t, originalContent, []byte("my backedup data"), "Restored file content")
}

func TestRestoreAllBakFilesFailsToWalk(t *testing.T) {
	// Save the original function to restore later
	originalWalkFunc := walkFunc

	// Override walkFunc to simulate failure
	walkFunc = func(root string, fn filepath.WalkFunc) error {
		return fmt.Errorf("failed to walk")
	}

	errs := RestoreAllBakFiles("path")

	if errs != nil {
		assert.Equal(t, "failed to walk", errs.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	// restore original walkFunc
	walkFunc = originalWalkFunc
}

func TestRestoreAllBakFilesFailsToCopyFiles(t *testing.T) {
	// Save the original function to restore later
	originalCopyFunc := copyFunc

	// Override walkFunc to simulate failure
	copyFunc = func(dst io.Writer, src io.Reader) (written int64, err error) {
		return 0, fmt.Errorf("failed to copy")
	}

	tempDir := t.TempDir()

	backedFile := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt.bak")

	path, _ := filepath.Abs(filepath.Dir(backedFile.Name()))

	errs := RestoreAllBakFiles(path)

	if errs != nil {
		assert.Equal(t, "failed to copy", errs.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	// restore original walkFunc
	copyFunc = originalCopyFunc
}

func TestRestoreAllBakFilesFailsToDelete(t *testing.T) {
	// Save the original function to restore later
	originalRemoveFunc := removeFunc

	// Override walkFunc to simulate failure
	removeFunc = func(name string) error {
		return fmt.Errorf("failed to remove")
	}

	tempDir := t.TempDir()

	backedFile := testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt.bak")

	path, _ := filepath.Abs(filepath.Dir(backedFile.Name()))

	errs := RestoreAllBakFiles(path)

	if errs != nil {
		assert.Equal(t, "failed to delete file "+backedFile.Name()+": failed to remove", errs.Error())
	} else {
		t.Error("Expected to return an error, but it did not")
	}

	// restore original removeFunc
	removeFunc = originalRemoveFunc
}

// TestGetFilesWithExt ensures that duplicate extensions will still return only the matching and no duplicates
func TestGetFilesWithExtDuplicateExtensions(t *testing.T) {
	tempDir := t.TempDir()

	_ = testutils.MustMakeTempFile(t, tempDir, "testFile-*.txt")
	ext, err := getFilesWithExt(tempDir, []string{".txt", ".txt"})
	if err != nil {
		t.Errorf("Failed to get the files with ext .txt - %s", err.Error())
	}
	assert.Equal(t, 1, len(ext), "Wrong number of files found")
}

// TestGetFilesWithExt ensures that duplicate extensions will still return only the matching and no duplicates
func TestGetFilesWithExtFailFindingDir(t *testing.T) {

	ext, err := getFilesWithExt("unknownTestFolder", []string{".txt", ".txt"})
	if err == nil {
		t.Errorf("Expected to fail as folder does not exist")
	} else {
		assert.Equal(t, 0, len(ext), "Wrong number of files found")
		if runtime.GOOS == "windows" {
			assert.Contains(t, err.Error(), "unknownTestFolder: The system cannot find the file specified.")
		} else {
			assert.Equal(t, "lstat unknownTestFolder: no such file or directory", err.Error())
		}
	}
}

func TestDeleteFilesNonExistentFile(t *testing.T) {
	err := deleteFiles([]string{"unknownTestFile"})
	if err == nil {
		t.Errorf("Expected to fail as file does not exist")
	} else {
		if runtime.GOOS == "windows" {
			assert.Equal(t, "failed to delete file unknownTestFile: remove unknownTestFile: The system cannot find the file specified.", err.Error())
		} else {
			assert.Equal(t, "failed to delete file unknownTestFile: remove unknownTestFile: no such file or directory", err.Error())
		}
	}
}

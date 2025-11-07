// Package io contains methods that can be useful for manipulating files/folders
package io

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// This is to enable tests, I hate it!
var (
	removeFunc = os.Remove
	walkFunc   = filepath.Walk
	openFunc   = os.Open
	createFunc = os.Create
	copyFunc   = io.Copy
)

// DeleteFileWithExt will delete all the files from the dir with extensions in the list
// extensions must be in the format .<ext> for matching
// Exact matches only, and the file end extension will be used to compare
// For example if the file found in dir is named file-name.txt.bak .bak is the determined extension for comparison
func DeleteFileWithExt(dir string, extensions []string) error {
	files, err := getFilesWithExt(dir, extensions)
	if err != nil {
		return err
	}
	return deleteFiles(files)
}

// BackupFilesWithExt take an inline (same folder) copy (adding suffix of .bak to the original names) of the files with the extensions
// It does not delete the original files
func BackupFilesWithExt(dir string, extensions []string) error {
	files, err := getFilesWithExt(dir, extensions)
	if err != nil {
		return err
	}
	err = backupFiles(files)
	if err != nil {
		return err
	}
	return nil
}

// RestoreAllBakFiles copies all the bak files into their respective original (dropping .bak extension) files and deletes the .bak files
func RestoreAllBakFiles(dir string) error {
	files, err := getFilesWithExt(dir, []string{".bak"})
	if err != nil {
		return err
	}

	err = restoreBakFiles(files)
	if err != nil {
		return err
	}
	return nil
}

func deleteFiles(files []string) error {
	for _, file := range files {
		err := removeFunc(file)
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %v", file, err)
		}
	}
	return nil
}

func getFilesWithExt(dir string, extensions []string) ([]string, error) {
	var files []string
	err := walkFunc(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			for _, ext := range extensions {
				if filepath.Ext(path) == ext {
					files = append(files, path)
					break
				}
			}
		}

		return nil
	})

	return files, err
}

func backupFiles(files []string) error {
	for _, file := range files {
		err := copyFile(file+".bak", file)
		if err != nil {
			return err
		}
	}
	return nil
}

func restoreBakFiles(files []string) error {
	for _, file := range files {
		err := copyFile(strings.TrimSuffix(file, ".bak"), file)
		if err != nil {
			return err
		}
	}

	err := deleteFiles(files)
	if err != nil {
		return err
	}
	return nil
}

func copyFile(dest string, source string) error {
	var sourceFile, err = openFunc(source)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			_ = fmt.Errorf("issue handling %s file", sourceFile.Name())
		}
	}(sourceFile)

	destFile, err := createFunc(dest)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {
			_ = fmt.Errorf("issue handling %s file", sourceFile.Name())

		}
	}(destFile)

	_, err = copyFunc(destFile, sourceFile) // should probably check the number of bytes copied
	if err != nil {
		return err
	}

	return nil
}

# mango4go - io

The `io` package is aimed to be a useful package for utilities working with the filesystem.

## DeleteFileWithExt

DeleteFileWithExt will delete all the files from the specified dir that have extensions matching in the list. The extensions list must contain values in the format .<ext>. It supports exact matches only, of the final file extension.

For example if the file found in dir is named file-name.txt.bak .bak is the determined extension for comparison.

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/io"

// ... rest of your code

// delete all the .txt files
err := DeleteFileWithExt(dirPath, []string{".txt"})
if errs != nil {
	// handle the error
}
```


## BackupFilesWithExt

BackupFilesWithExt will create an inline (same folder) copy (adding suffix of .bak to the original names) of the files with the extensions.

It does not delete the original files

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/io"

// ... rest of your code

// backup all .txt files
err := BackupFilesWithExt(dirPath, []string{".txt"})
if errs != nil {
	// handle the error
}
// all the files .txt from dirPath will be copied into .txt.bak
```


## RestoreAllBakFiles

RestoreAllBakFiles copies all the bak files into their respective original (dropping .bak extension) files and deletes the .bak files

Existing original files will be overwritten.


### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/io"

// ... rest of your code

// backup all .txt files
err := BackupFilesWithExt(dirPath, []string{".txt"})
if errs != nil {
	// handle the error
}
// all the files .txt from dirPath will be copied into .txt.bak
```

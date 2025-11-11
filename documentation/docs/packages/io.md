# `pkg/io`

Small filesystem helpers for bulk operations (delete, backup, restore) based on file extensions. The package walks a
directory tree, so you can prepare rollbacks or cleanups with a single call.

## Quick Start

```go
import mangoio "github.com/bitstep-ie/mango-go/pkg/io"

func rotateConfigs(dir string) error {
if err := mangoio.BackupFilesWithExt(dir, []string{".yaml", ".json"}); err != nil {
return err
}

// … mutate configs here …

return mangoio.RestoreAllBakFiles(dir)
}
```

## API Overview

| Function                                     | Description                                                |
|----------------------------------------------|------------------------------------------------------------|
| `SafeClose(c)`                               | safely close closer c, no output                           |
| `SafeClosePrint(c)`                          | safely close closer c, and output the error if any         |
| `DeleteFileWithExt(dir, []string{".log"})`   | remove every matching file (no prompt)                     |
| `BackupFilesWithExt(dir, []string{".conf"})` | copy `foo.conf → foo.conf.bak`                             |
| `RestoreAllBakFiles(dir)`                    | copy `*.bak` back to originals and remove the `.bak` files |

> Extensions must include the leading dot (`.log`). Matches are based on the final path extension (`file.txt.bak` is
> treated as `.bak`).

## Usage Patterns

### Cleaning Generated Files

```go
func cleanArtifacts(dir string) error {
return mangoio.DeleteFileWithExt(dir, []string{".tmp", ".bak"})
}
```

### Safe Inline Backups

```go
func upgrade(dir string) error {
// 1. duplicate target files
if err := mangoio.BackupFilesWithExt(dir, []string{".json"}); err != nil {
return err
}

// 2. apply your upgrade logic
if err := runMigration(dir); err != nil {
// roll back
_ = mangoio.RestoreAllBakFiles(dir)
return err
}

// 3. cleanup backups if desired
return mangoio.DeleteFileWithExt(dir, []string{".bak"})
}
```

## Testing Hooks

Internally the package swaps `os.Remove`, `filepath.Walk`, and copy operations with overridable function variables to
simplify unit tests. In production you never touch these, but the pattern means you can inject fakes when writing tests
for your code that depends on `pkg/io`.

## Tips

- Run these helpers on the narrowest directory possible to avoid traversing large trees.
- All operations are best-effort: once a delete or copy fails, the function returns the first error with the file path.
- `RestoreAllBakFiles` overwrites existing originals. Use `BackupFilesWithExt` immediately before mutations if you need
  a guaranteed rollback path.

# `pkg/testutils`

Testing helpers that reduce boilerplate when working with temporary files and structured assertions.

```go
import (
    "os"
    "testing"
    testutils "github.com/bitstep-ie/mango-go/pkg/testutils"
)

func TestProcessor(t *testing.T) {
    tmp := testutils.MustMakeTempFile(t, t.TempDir(), "payload-*.json")
    defer os.Remove(tmp.Name())

    // write fixture data, run code under test…

    testutils.AssertValidUUID(t, got.LogID, "logId")
    testutils.ContainsAllRunes(t, got.Token, "!@#", "token must include punctuation")
}
```

## Helpers

### `MustMakeTempFile(t, dir, pattern) *os.File`

- Creates a file using `os.CreateTemp`.
- Fails the test via `assert.Fail` if creation fails.
- Closes the returned file automatically; reopen with `os.Open` if you need to append later.

### `AssertValidUUID(t, value, fieldName)`

- Delegates to `uuid.Parse`.
- Fails the test with a friendly message that includes `fieldName`.

### `ContainsAllRunes(t, str, chars, msgAndArgs...)`

- Ensures every rune in `chars` exists at least once in `str`.
- Useful for asserting generated passwords or tokens meet composition rules.

## Tips

- The helpers pull in `stretchr/testify/assert`, so you can mix these with other testify assertions without additional setup.
- When using `MustMakeTempFile`, prefer `t.TempDir()` so Go cleans up the directory after the test run.

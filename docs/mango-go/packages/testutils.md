# mango4go - testutils

The `testutils` package is to contain useful functions for testing


## MustMakeTempFile
It will create a temporary file in the dir using the fileNamePattern.
Failure to create the temp file will result in t.Fatalf

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/testutils"

// ... rest of your code

// create the temporary file
tempFile := testutils.MustMakeTempFile(t, t., "tempfile-*.txt")
```


## AssertValidUUID
Validates the value sent in is a valid uuid, by attempting to parse the value using `github.com/google/uuid` `Parse`, and any errors are reported as assertion.
The fieldName passed in will be used in error message printout.

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/testutils"

// ... rest of your code

// assert valid UUID
invalidUuidValue := "invalidUUID"
validUuidValue := "f7198919-bb3e-4ff0-9746-72d913f8a812"
testutils.AssertValidUUID(t, invalidUuidValue, "correlationId") // this will assert
testutils.AssertValidUUID(t, validUuidValue, "correlationId") // this will pass
```

## ContainsAllRunes
Validates all runes in `chars` exist in `str`.


### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/testutils"

// ... rest of your code

// assert chars in str
str := "abcdef"
chars := "ace"
testutils.ContainsAllRunes(t, str, chars, "Message for failure")
```

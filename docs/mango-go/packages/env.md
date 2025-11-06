# mango4go - env

The `env` package provides small functions for environment variables manipulation.


## EnvOrDefault
It returns the value of the environment variable, or the default specified.

*Note:* the check for environment variable existence is done by comparison to empty string (`""`). 

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/env"

// ... rest of your code

// get the environment variable
value := env.EnvOrDefault("ENV_KEY", "some-default-value")
// now value will have ENV_KEY value if present, or "some-default-value" if ENV_KEY doesn't exist in the environment
```


## MustEnv
It returns the value of the environment variable or panics

*Note:* the check for environment variable existence is done by comparison to empty string (`""`).

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/env"

// ... rest of your code

// get the environment variable
value := env.MustEnv("ENV_KEY")
// now value will have ENV_KEY value
// panics otherwise
```


## EnvAsInt
It returns the value of the environment variable as `integer`, and the default value if it doesn't exist.

Panic if the value is **NOT** an `integer`.

*Note:* the check for environment variable existence is done by comparison to empty string (`""`).

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/env"

// ... rest of your code

// get the environment variable
value := env.EnvAsInt("ENV_KEY", 7)
// now value will have ENV_KEY value as integer (panic if not an integer)
// or the default value 7 if it doesn't exist
```



## MustEnvAsInt
It returns the value of the environment variable as `integer`, and panics if it doesn't exist.

Panic if the value is **NOT** an `integer`.

*Note:* the check for environment variable existence is done by comparison to empty string (`""`).

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/env"

// ... rest of your code

// get the environment variable
value := env.MustEnvAsInt("ENV_KEY")
// now value will have ENV_KEY value as integer (panic if not an integer)
// or panics if it doesn't exist
```


## EnvAsBool
It returns the value of the environment variable as `bool`, and default value if it doesn't exist.

Panic if the value is **NOT** a `bool`.

*Note:* the check for environment variable existence is done by comparison to empty string (`""`).

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/env"

// ... rest of your code

// get the environment variable
value := env.EnvAsBool("ENV_KEY", false)
// now value will have ENV_KEY value as bool (panic if not a bool)
// or default value if it doesn't exist
```


## MustEnvAsBool
It returns the value of the environment variable as `bool`, and panics if it doesn't exist.

Panic if the value is **NOT** a `bool`.

*Note:* the check for environment variable existence is done by comparison to empty string (`""`).

### How to use it?

```go language=go
// Import the library package desired
import "github.com/bitstep-ie/mango-go/env"

// ... rest of your code

// get the environment variable
value := env.MustEnvAsBool("ENV_KEY")
// now value will have ENV_KEY value as bool (panic if not a bool)
// or panics if it doesn't exist
```
# `pkg/env`

Helpers for reading configuration from environment variables with sensible defaults and explicit “must” variants that panic when a value is missing or malformed.

## Quick Start

```go
import mangoenv "github.com/bitstep-ie/mango-go/pkg/env"

type Config struct {
    Port        int
    EnableHTTPS bool
    Secret      string
}

func Load() Config {
    return Config{
        Port:        mangoenv.EnvAsInt("PORT", 8080),
        EnableHTTPS: mangoenv.EnvAsBool("ENABLE_HTTPS", false),
        Secret:      mangoenv.MustEnv("API_SECRET"),
    }
}
```

## API Cheatsheet

| Function | Purpose |
| --- | --- |
| `EnvOrDefault(key, fallback string)` | string with fallback |
| `MustEnv(key)` | string or panic if empty |
| `EnvAsInt(key, fallback int)` / `MustEnvAsInt(key)` | parse integer values |
| `EnvAsBool(key, fallback bool)` / `MustEnvAsBool(key)` | parse boolean values |

Each helper treats “missing” as `""` and panics with descriptive messages for invalid conversions.

## Examples

### Strings

```go
name := mangoenv.EnvOrDefault("SERVICE_NAME", "checkout")
token := mangoenv.MustEnv("BEARER_TOKEN") // panic if unset
```

### Integers

```go
maxConnections := mangoenv.EnvAsInt("MAX_CONN", 10)
timeout := mangoenv.MustEnvAsInt("REQUEST_TIMEOUT_SECONDS")
```

### Booleans

```go
debug := mangoenv.EnvAsBool("DEBUG", false)
tlsOnly := mangoenv.MustEnvAsBool("TLS_ONLY")
```

## Tips

- Use `EnvAs*` when you can tolerate defaults (local dev) and `MustEnv*` for production-critical knobs.
- Panics happen on invalid formats (e.g., `EnvAsInt("PORT")` with `PORT=abc`). Keep these calls near bootstrapping code so the service fails fast.
- Wrap lookups in a struct constructor (see Quick Start) to centralize configuration logic.

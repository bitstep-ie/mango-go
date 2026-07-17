# `pkg/env`

Helpers for reading configuration from environment variables. Every helper has three flavours:

| Flavour | Signature pattern | Behaviour when missing or invalid |
| --- | --- | --- |
| `*Default` | `Foo(key, default) (T, error)` | returns the default on missing, error on invalid format |
| bare | `Foo(key) (T, error)` | returns error on missing **or** invalid format |
| `Must*` | `MustFoo(key) T` | panics on missing **or** invalid format |

This makes it easy to choose the right trade-off: propagate errors in library code, use defaults for optional config, panic fast for required production values.

## Quick Start

```go
import mangoenv "github.com/bitstep-ie/mango-go/pkg/env"

type Config struct {
    Port        int
    EnableHTTPS bool
    Secret      string
}

func Load() (Config, error) {
    port, err := mangoenv.IntDefault("PORT", 8080)
    if err != nil {
        return Config{}, err
    }
    https, err := mangoenv.BoolDefault("ENABLE_HTTPS", false)
    if err != nil {
        return Config{}, err
    }
    return Config{
        Port:        port,
        EnableHTTPS: https,
        Secret:      mangoenv.MustString("API_SECRET"), // panic if unset
    }, nil
}
```

## API Reference

### Strings

| Function | Returns | Notes |
| --- | --- | --- |
| `StringDefault(key, fallback string)` | `(string, error)` | fallback when unset |
| `String(key string)` | `(string, error)` | error when unset |
| `MustString(key string)` | `string` | panics when unset |

### Integers

| Function | Returns | Notes |
| --- | --- | --- |
| `IntDefault(key string, fallback int)` | `(int, error)` | fallback when unset; error when non-integer |
| `Int(key string)` | `(int, error)` | error when unset or non-integer |
| `MustInt(key string)` | `int` | panics when unset or non-integer |

### Booleans

Accepts any value understood by [`strconv.ParseBool`](https://pkg.go.dev/strconv#ParseBool): `1`, `t`, `T`, `TRUE`, `true`, `True`, `0`, `f`, `F`, `FALSE`, `false`, `False`.

| Function | Returns | Notes |
| --- | --- | --- |
| `BoolDefault(key string, fallback bool)` | `(bool, error)` | fallback when unset; error when non-bool |
| `Bool(key string)` | `(bool, error)` | error when unset or non-bool |
| `MustBool(key string)` | `bool` | panics when unset or non-bool |

### Durations

Parsed with [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration). Valid units: `ns`, `us`, `ms`, `s`, `m`, `h`.

| Function | Returns | Notes |
| --- | --- | --- |
| `DurationDefault(key string, fallback time.Duration)` | `(time.Duration, error)` | fallback when unset; error when invalid |
| `Duration(key string)` | `(time.Duration, error)` | error when unset or invalid |
| `MustDuration(key string)` | `time.Duration` | panics when unset or invalid |

### String slices

Split by a caller-supplied separator. No trimming or filtering is applied — the raw `strings.Split` result is returned.

| Function | Returns | Notes |
| --- | --- | --- |
| `StringSpliceDefault(key, sep string, fallback []string)` | `([]string, error)` | fallback when unset |
| `StringSplice(key, sep string)` | `([]string, error)` | error when unset |
| `MustStringSplice(key, sep string)` | `[]string` | panics when unset |

### Custom / generic

Bring your own parser `func(string) (T, error)`.

| Function | Returns | Notes |
| --- | --- | --- |
| `AsDefault[T](key string, fallback T, parser)` | `(T, error)` | fallback when unset; error when parser fails |
| `As[T](key string, parser)` | `(T, error)` | error when unset or parser fails |
| `MustAs[T](key string, parser)` | `T` | panics when unset or parser fails |

## Examples

### Strings

```go
name, err := mangoenv.StringDefault("SERVICE_NAME", "checkout")
token := mangoenv.MustString("BEARER_TOKEN") // panic if unset
```

### Integers

```go
maxConn, err := mangoenv.IntDefault("MAX_CONN", 10)
timeout := mangoenv.MustInt("REQUEST_TIMEOUT_SECONDS")
```

### Booleans

```go
debug, err := mangoenv.BoolDefault("DEBUG", false)
tlsOnly := mangoenv.MustBool("TLS_ONLY")
```

### Durations

```go
ttl, err := mangoenv.DurationDefault("CACHE_TTL", 15*time.Minute)
lease := mangoenv.MustDuration("LEASE_DURATION")
```

### String slices

```go
// comma-separated
hosts, err := mangoenv.StringSpliceDefault("ALLOWED_HOSTS", ",", []string{"localhost"})

// pipe-separated, required
admins := mangoenv.MustStringSplice("ADMIN_EMAILS", "|")
```

### Custom parsing

```go
workers, err := mangoenv.AsDefault("WORKERS", 4, strconv.Atoi)

buildID := mangoenv.MustAs("BUILD_ID", func(v string) (uuid.UUID, error) {
    return uuid.Parse(v)
})
```

## Deprecated helpers

The following names are kept for backwards compatibility and are scheduled for removal in `v1.0.0`:

| Old name | Replacement |
| --- | --- |
| `EnvOrDefault(key, fallback)` | `StringDefault(key, fallback)` |
| `MustEnv(key)` | `MustString(key)` |
| `EnvAsInt(key, fallback)` | `IntDefault(key, fallback)` |
| `MustEnvAsInt(key)` | `MustInt(key)` |
| `EnvAsBool(key, fallback)` | `BoolDefault(key, fallback)` |
| `MustEnvAsBool(key)` | `MustBool(key)` |

## Tips

- Keep all env reads in one place (e.g. a `config.Load()` function) so misconfiguration fails fast and is easy to spot.
- Use `*Default` variants for optional knobs with sane defaults; use `Must*` for production-critical values that have no sensible fallback.
- `Must*` panics propagate to the top of your `main()` — this is intentional. A missing required value should crash the process immediately, not silently misbehave later.

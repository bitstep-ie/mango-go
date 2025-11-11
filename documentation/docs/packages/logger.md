# `pkg/logger`

A structured logging handler built on top of `log/slog`. It fans each record out to CLI, rotating files (via `lumberjack`), and syslog simultaneously, while enforcing context contracts such as `type`, `application`, and `operation`.

## Quick Start

```go
import (
    "context"
    "log/slog"
    mangolog "github.com/bitstep-ie/mango-go/pkg/logger"
)

func newLogger() *slog.Logger {
    cfg := &mangolog.LogConfig{
        MangoConfig: &mangolog.MangoConfig{
            Strict: true,
            CorrelationId: &mangolog.CorrelationIdConfig{
                Strict:       true,
                AutoGenerate: true,
            },
        },
        Out: &mangolog.OutConfig{
            Enabled: true,
            Cli: &mangolog.CliConfig{
                Enabled:  true,
                Friendly: true,
                Verbose:  false,
            },
            File: &mangolog.FileOutputConfig{
                Enabled:  true,
                Debug:    false,
                Path:     "/var/log/service.log",
                MaxSize:  50,
                Compress: true,
            },
        },
    }
    return slog.New(mangolog.NewMangoLogger(cfg))
}

func logExample(logger *slog.Logger) {
    ctx := context.Background()
    ctx = context.WithValue(ctx, mangolog.APPLICATION, "checkout-api")
    ctx = context.WithValue(ctx, mangolog.OPERATION, "cart-create")
    ctx = context.WithValue(ctx, mangolog.TYPE, mangolog.BusinessType)

    logger.InfoContext(ctx, "cart created",
        slog.Int("items", 3),
        slog.String("country", "IE"),
    )
}
```

## Configuration

`LogConfig` is split into:

- `MangoConfig`: strict mode + correlation-id behaviour.
- `Out`: toggles for `File`, `Cli`, and `Syslog`.

```yaml
mango:
  strict: true
  correlation-id:
    strict: true
    auto-generate: true
out:
  enabled: true
  cli:
    enabled: true
    friendly: true
    friendly-format: '"[\(.level)] \(.operation) - \(.message)"'
    verbose: true
    verbose-format: "."
  file:
    enabled: true
    debug: false
    path: /var/log/mango.log
    max-size: 100
    max-backups: 5
    max-age: 7
    compress: true
  syslog:
    facility: local0
```

Friendly/verbose formats consume jq strings (`gojq`) and default to built-in templates when left empty.

## Context Requirements

Strict mode enforces presence (and validity) of:

- `mangolog.TYPE` – must be one of `Business`, `Security`, `Performance`.
- `mangolog.APPLICATION`
- `mangolog.OPERATION`
- `mangolog.CORRELATION_ID` (when `correlation-id.strict` is true; auto-generated if `auto-generate` is true).

On missing or invalid fields, `Handle` logs an error and returns it to the slog caller.

## Outputs

### CLI

- When `friendly` is true, Mango Logger runs the log through the jq template (e.g., `"[INFO] create - success"`).
- Otherwise, it prints raw JSON to stdout/stderr.
- `verbose` gates debug logs on stdout and includes correlation IDs for INFO-level messages.

### File

- Writes newline-delimited JSON (`StructuredLog`) using lumberjack rotation.
- `debug` controls whether `LevelDebug` entries reach the file.

### Syslog

- Enabled by setting `out.syslog.facility` or the corresponding constant (e.g., `mangolog.SyslogFacilityLocal0`).
- Severity is derived from the slog level.
- Not available on Windows (build tags guard the implementation).

## Structured Output

```json
{
  "ts": "2025-01-15T09:53:34.717-0500",
  "type": "Business",
  "application": "checkout-api",
  "operation": "cart-create",
  "correlationid": "a52b0129-9d49-4f29-acbb-3575aa4442f4",
  "logId": "67e36893-0a7e-476c-b799-4a2772e9bd17",
  "level": "INFO",
  "message": "cart created",
  "attributes": {
    "items": 3,
    "country": "IE"
  }
}
```

## Tips

1. Use middleware to stamp context keys (`TYPE`, `APPLICATION`, `OPERATION`, `CORRELATION_ID`) once per request.
2. Toggle `Cli.Verbose` via CLI flags (`--verbose`) to expose debug logs during troubleshooting.
3. When `Strict` is enabled, avoid mutating the global `REQUIRED_FIELDS`; create fresh contexts per request to prevent leaking values across goroutines.

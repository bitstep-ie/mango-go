# Usage Guide
This guide walks you through configuring a Go project that consumes the `github.com/bitstep-ie/mango-go/pkg/...` packages. It covers end-to-end setup plus focused examples for each utility package so you can copy, paste, and start shipping features quickly.

---

## 1. Project Setup

```bash
go env -w GO111MODULE=on
go get github.com/bitstep-ie/mango-go@latest
```

Inside your module, import the packages you need:

```go
import (
    mangoenv "github.com/bitstep-ie/mango-go/pkg/env"
    mangoio "github.com/bitstep-ie/mango-go/pkg/io"
    mangolog "github.com/bitstep-ie/mango-go/pkg/logger"
    mangorand "github.com/bitstep-ie/mango-go/pkg/random"
    mangoslices "github.com/bitstep-ie/mango-go/pkg/slices"
    mangotime "github.com/bitstep-ie/mango-go/pkg/time"
    testutils "github.com/bitstep-ie/mango-go/pkg/testutils"
)
```

> Tip: keep the `mango-` prefixes when aliasing packages to avoid name collisions with the stdlib (`io`, `time`, `slices`, etc.).

---

## 2. Bootstrapping the Mango Logger

The logger wraps `log/slog` and supports friendly CLI output, JSON file output (via lumberjack rotation), and syslog simultaneously.

```go
package main

import (
    "context"
    "log/slog"
    mangolog "github.com/bitstep-ie/mango-go/pkg/logger"
)

func main() {
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
                Enabled:       true,
                Friendly:      true,
                Verbose:       true,
                FriendlyFormat: ``"[\(.level)] \(.application)::\(.operation) - \(.message)"``,
            },
            File: &mangolog.FileOutputConfig{
                Enabled:    true,
                Debug:      false,
                Path:       "/var/log/my-service.log",
                MaxSize:    25,
                MaxBackups: 5,
                MaxAge:     14,
                Compress:   true,
            },
            Syslog: &mangolog.SyslogConfig{
                Facility: mangolog.SyslogFacilityLocal0,
            },
        },
    }

    handler := mangolog.NewMangoLogger(cfg)
    logger := slog.New(handler)

    ctx := context.WithValue(context.Background(), mangolog.OPERATION, "checkout")
    ctx = context.WithValue(ctx, mangolog.APPLICATION, "orders-api")
    ctx = context.WithValue(ctx, mangolog.TYPE, mangolog.BusinessType)

    logger.InfoContext(ctx, "order created",
        slog.Int("orderID", 42),
        slog.String("status", "submitted"),
    )
}
```

Key points:

- Set `MangoConfig.Strict` to enforce `type`, `application`, and `operation` context values.
- Enable `CorrelationId.AutoGenerate` to avoid panics when downstream middleware forgets to set it.
- CLI output uses jq expressions (`DefaultFriendlyFormat`) so you can tailor human logs without touching JSON/file/syslog payloads.

---

## 3. Managing Environment Variables

```go
package config

import mangoenv "github.com/bitstep-ie/mango-go/pkg/env"

type Config struct {
    Port        int
    EnableHTTPS bool
    DBURL       string
}

func Load() Config {
    return Config{
        Port:        mangoenv.EnvAsInt("PORT", 8080),
        EnableHTTPS: mangoenv.EnvAsBool("ENABLE_HTTPS", false),
        DBURL:       mangoenv.MustEnv("DATABASE_URL"),
    }
}
```

- Use `EnvOrDefault`/`EnvAsInt` for soft defaults.
- Prefer `MustEnv*` variants for critical values so misconfiguration fails fast during startup.

---

## 4. File Maintenance Helpers (`pkg/io`)

```go
package backup

import mangoio "github.com/bitstep-ie/mango-go/pkg/io"

func RotateConfigs(dir string) error {
    if err := mangoio.BackupFilesWithExt(dir, []string{".yaml", ".json"}); err != nil {
        return err
    }

    // run your migration/update logic here

    return mangoio.RestoreAllBakFiles(dir)
}
```

`DeleteFileWithExt` works the same way and is handy for cleaning generated artifacts (`.tmp`, `.bak`, etc.). The helpers walk the entire tree, so point them at the narrowest directory possible.

---

## 5. Random Data Recipes (`pkg/random`)

```go
package demo

import (
    "time"
    mangorand "github.com/bitstep-ie/mango-go/pkg/random"
)

func init() {
    rand.Seed(time.Now().UnixNano()) // only needed once per process
}

func buildFixtures() {
    run := mangorand.Number[int](1, 5)
    coinFlip := mangorand.Bool()
    promoCode := mangorand.String(10)
    securePassword := mangorand.Password(24, mangorand.PasswordOptions{
        Letters: true,
        Digits:  true,
        Symbols: true,
        Exclude: "O0l1",
    })
    randomRunAt := mangorand.Date(time.Now().Add(-30*24*time.Hour), time.Now())

    _ = []any{run, coinFlip, promoCode, securePassword, randomRunAt}
}
```

The package mixes math/rand conveniences with crypto-safe building blocks (`Byte`, `Password`). Always seed `math/rand` in your main package if you rely on `Number`, `Bool`, `Choice`, etc.

---

## 6. Slice Utilities

```go
events := []string{"created", "updated", "deleted", "deleted"}
unique := mangoslices.Unique(events) // ["created","updated","deleted"]

if mangoslices.Contains(events, "created") {
    idx := mangoslices.IndexOf(events, "deleted") // first match
    copy := append([]string(nil), events...)
    mangoslices.Reverse(copy)
}

pairs := mangoslices.EqualsIgnoreOrder(
    []int{1, 2, 2, 3},
    []int{2, 1, 3, 2},
) // true

batched := mangoslices.Chunk(events, 2)
```

Generics keep these helpers type-safe and allocation-friendly. `Chunk` panics when `size <= 0`; validate user input before calling it.

---

## 7. Working with Time

```go
now := time.Now()
start := mangotime.StartOfDay(now)
end := mangotime.EndOfDay(now)

if mangotime.IsToday(now) {
    fmt.Println("still today")
}

deadline, err := mangotime.ParseDuration("1w12h30m")
if err != nil { /* handle */ }

fmt.Println(mangotime.TimeAgo(time.Now().Add(-90 * time.Minute))) // "1 hour ago"
```

`ParseDuration` extends `time.ParseDuration` with `d` and `w` suffixes (supports floats like `1.5d`). `TimeAgo` returns human-friendly English strings for UI messaging.

---

## 8. Test Utilities

```go
func TestProcessor(t *testing.T) {
    tmp := testutils.MustMakeTempFile(t, t.TempDir(), "payload-*.json")

    payload := uuid.NewString()
    require.NoError(t, os.WriteFile(tmp.Name(), []byte(payload), 0o600))

    runProcessor(tmp.Name())

    testutils.AssertValidUUID(t, got.ID, "logId")
    testutils.ContainsAllRunes(t, got.Token, "!@#", "token should include special chars")
}
```

- `MustMakeTempFile` simplifies fixture creation; it will close the handle automatically.
- `AssertValidUUID` and `ContainsAllRunes` piggyback on `stretchr/testify` for idiomatic assertions.

---

## 9. Putting It All Together

Here’s a minimal HTTP handler showcasing multiple packages:

```go
type Server struct {
    logger *slog.Logger
    cfg    Config
}

func NewServer(cfg Config) *Server {
    handler := mangolog.NewMangoLogger(buildLogConfig(cfg))
    return &Server{
        logger: slog.New(handler),
        cfg:    cfg,
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := context.WithValue(r.Context(), mangolog.APPLICATION, "billing-api")
    ctx = context.WithValue(ctx, mangolog.OPERATION, "invoice-create")
    ctx = context.WithValue(ctx, mangolog.TYPE, mangolog.BusinessType)

    reqID := r.Header.Get("X-Request-ID")
    if reqID == "" {
        reqID = mangorand.String(12)
    }
    ctx = context.WithValue(ctx, mangolog.CORRELATION_ID, reqID)

    s.logger.InfoContext(ctx, "processing request",
        slog.String("method", r.Method),
        slog.String("path", r.URL.Path),
        slog.String("receivedAt", mangotime.StartOfDay(time.Now()).Format(mangolog.RFC3339NanoMC)),
    )

    w.WriteHeader(http.StatusAccepted)
}
```

This pattern scales: standardize context values at the edge, use Mango Logger for structured telemetry, derive configuration via `pkg/env`, and lean on the helper packages for the rest.

---

## 10. Next Steps

1. Browse `docs/mango-go/packages/*.md` for deeper API details.
2. Wire these helpers into your project’s scaffolding (CLI, HTTP server, worker, etc.).
3. Contribute improvements or new helpers by following `contributing.md`.

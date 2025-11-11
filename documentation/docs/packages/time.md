# `pkg/time`

Utility functions for common calendaring tasks (start/end-of-day boundaries, relative checks), human-friendly “time ago” formatting, and extended duration parsing.

## Quick Start

```go
import mangotime "github.com/bitstep-ie/mango-go/pkg/time"

now := time.Now()
start := mangotime.StartOfDay(now)
end := mangotime.EndOfDay(now)

if mangotime.IsToday(now) {
    fmt.Println("still today")
}

deadline, _ := mangotime.ParseDuration("1w2d12h")
fmt.Println(mangotime.TimeAgo(time.Now().Add(-90 * time.Minute))) // "1 hour ago"
```

## Calendar Helpers

| Function | Description |
| --- | --- |
| `StartOfDay(t)` | midnight in the same location |
| `EndOfDay(t)` | last nanosecond of the day |
| `IsToday(t)` / `IsTomorrow(t)` | check relative to local time |
| `IsTodayLoc(t, loc)` / `IsTomorrowLoc(t, loc)` | same checks using a custom `*time.Location` |

Both “Is” helpers treat the exact boundary values as matches (`StartOfDay(t)` is “today”; `EndOfDay(t)` is still “today”).

## Duration Parsing

`ParseDuration` extends Go’s built-in parser with `d` (days) and `w` (weeks), including fractional values.

```go
dur, err := mangotime.ParseDuration("2w1.5d30m")
// -> 2 weeks + 1.5 days + 30 minutes
```

If the string is empty, the function returns an error. For other invalid formats, it falls back to `time.ParseDuration`’s error.

## Relative Formatting

`TimeAgo(t)` (and the internal `timeAgoWithNow(t, now)`) return concise English strings:

- `<1m` → `just now`
- `=1m` → `1 minute ago`
- `<1h` → `X minutes ago`
- `<24h` → `X hours ago`
- `<48h` → `yesterday`
- otherwise → `X days ago`

Ideal for CLI or UI messages without pulling in a full i18n library.

## Tips

- Use `IsTodayLoc` / `IsTomorrowLoc` when you store timestamps in UTC but need to reason about a customer’s timezone.
- Combine `StartOfDay` + `Add(24*time.Hour)` to build custom windows (e.g., “this business day”).
- When parsing human input (cron-like configs, CLI flags), accept strings such as `90m`, `2d`, `1w12h30m` to keep UX flexible.

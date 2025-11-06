# mango4go - time

The `time` package provides time helper functions

## StartOfDay
Returns the date at `00:00:00.000` in the same location.

## EndOfDay
Returns the last moment of the day `23:59:59.999999999` in the same location.

## IsToday 
It checks if a timestamp is `today` in the current local timezone, compares after StartOfDay and before StartOfDay + 24h.

*Notes:*
  - Exact StartOfDay today is considered today
  - Exact EndOfDay today is considered today

## IsTodayLoc 
It checks if a timestamp is `today` in the specified `time.Location`, compares after StartOfDay and before StartOfDay + 24h.
*Notes:*
  - Exact StartOfDay today is considered today
  - Exact EndOfDay today is considered today

## IsTomorrow 
Checks if a timestamp is `tomorrow` in the current local timezone.
Compares after StartOfDay + 24h and before StartOfDay + 48h.

*Notes:*
  - Exact StartOfDay tomorrow is considered tomorrow
  - Exact EndOfDay tomorrow is considered tomorrow

## IsTomorrowLoc
Checks if a timestamp is `tomorrow` in the specified `time.Location`. Compares after StartOfDay + 24h and before StartOfDay + 48h.
*Notes:*
  - Exact StartOfDay tomorrow is considered tomorrow
  - Exact EndOfDay tomorrow is considered tomorrow

## ParseDuration 
This extends `time.ParseDuration` to support `d` (days) and `w` (weeks).
Replaces `d` with `24h`, `w` with `168h`, then pass to `time.ParseDuration`.

### Example
`2w1d2h30m` would return `362h30m0s` computed as: `2*168*time.Hour + 1*24*time.Hour + 2*time.Hour + 30*time.Minute`

## TimeAgo 
Get a human friendly relative time formatting as follows:
  -  <1m → "just now"
  -  =1m → "1 minute ago"
  -  <1h → "X minutes ago"
  -  =24h → "1 hour ago"
  -  <24h → "X hours ago"
  -  <48h → "yesterday"
  -  otherwise → "X days ago"
Currently this offers fixed english versions, we can explore the options of using i18n in a future release.
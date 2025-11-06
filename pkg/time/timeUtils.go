package time

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// StartOfDay returns the date at 00:00:00.000 in the same location
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the last moment of the day
// 23:59:59.999999999 in the same location
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).AddDate(0, 0, 1).Add(-time.Nanosecond)
}

// IsToday checks if a timestamp is today in the current local timezone
// Compares after StartOfDay and before StartOfDay + 24h
// Exact StartOfDay today is considered today
// Exact EndOfDay today is considered today
func IsToday(t time.Time) bool {
	return IsTodayLoc(t, time.Now().Location())
}

// IsTodayLoc checks if a timestamp is today in the specified time.Location
// Compares after StartOfDay and before StartOfDay + 24h
// Exact StartOfDay today is considered today
// Exact EndOfDay today is considered today
func IsTodayLoc(t time.Time, loc *time.Location) bool {
	now := time.Now()
	now = now.In(loc)
	t = t.In(loc)
	todayStart := StartOfDay(now)
	tomorrowStart := todayStart.Add(24 * time.Hour)
	return t.Equal(todayStart) || (t.After(todayStart) && t.Before(tomorrowStart))
}

// IsTomorrow checks if a timestamp is tomorrow in the current local timezone
// Compares after StartOfDay + 24h and before StartOfDay + 48h
// Exact StartOfDay tomorrow is considered tomorrow
// Exact EndOfDay tomorrow is considered tomorrow
func IsTomorrow(t time.Time) bool {
	return IsTomorrowLoc(t, time.Now().Location())
}

// IsTomorrowLoc checks if a timestamp is tomorrow in the specified time.Location
// Compares after StartOfDay + 24h and before StartOfDay + 48h
// Exact StartOfDay tomorrow is considered tomorrow
// Exact EndOfDay tomorrow is considered tomorrow
func IsTomorrowLoc(t time.Time, loc *time.Location) bool {
	now := time.Now()
	now = now.In(loc)
	t = t.In(loc)
	tomorrowStart := StartOfDay(now.Add(24 * time.Hour))
	dayAfterTomorrow := tomorrowStart.Add(24 * time.Hour)
	return t.Equal(tomorrowStart) || (t.After(tomorrowStart) && t.Before(dayAfterTomorrow))
}

// ParseDuration extends time.ParseDuration to support "d" (days) and "w" (weeks)
// Replaces d with 24h, w with 168h, then pass to time.ParseDuration
// Supports float values too
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("empty duration string")
	}
	dwPattern := regexp.MustCompile(`(\d+(\.\d+)?)([dw])`)

	processed := dwPattern.ReplaceAllStringFunc(s, func(match string) string {
		submatches := dwPattern.FindStringSubmatch(match)

		valueStr := submatches[1]
		unit := submatches[3]

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return match // fallback, will be caught later by ParseDuration
		}

		var hours float64
		switch unit {
		case "d":
			hours = value * 24
		case "w":
			hours = value * 168
		default:
			hours = value // should not happen
		}

		return fmt.Sprintf("%gh", hours)
	})

	// Now parse with time.ParseDuration
	return time.ParseDuration(processed)
}

// TimeAgo returns a human friendly relative time formatting
// <1m → "just now"
// =1m → "1 minute ago"
// <1h → "X minutes ago"
// <24h → "X hours ago"
// <48h → "yesterday"
// otherwise → "X days ago"
// Explore the options of using i18n in a future release
func TimeAgo(t time.Time) string {
	return timeAgoWithNow(t, time.Now())
}

// Internal use method to help with coverage of tests to cover for exact matches of diff
// Not recommended for external use - hence not exposed
func timeAgoWithNow(t time.Time, now time.Time) string {
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return fmt.Sprintf("%d minute ago", minutes)
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return fmt.Sprintf("%d hour ago", hours)
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if diff < 48*time.Hour {
		return "yesterday"
	}
	days := int(diff.Hours() / 24)
	return fmt.Sprintf("%d days ago", days)
}

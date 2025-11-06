package time

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartOfDay(t *testing.T) {
	loc := time.FixedZone("TestZone", 3600)
	input := time.Date(2025, 9, 5, 17, 11, 13, 123456789, loc)
	expected := time.Date(2025, 9, 5, 0, 0, 0, 0, loc)

	result := StartOfDay(input)
	assert.Equal(t, expected, result, "StartOfDay should reset time to midnight")
}

func TestEndOfDay(t *testing.T) {
	loc := time.FixedZone("TestZone", 3600)
	input := time.Date(2025, 9, 5, 17, 11, 13, 123456789, loc)
	expected := time.Date(2025, 9, 5, 23, 59, 59, 999999999, loc)

	result := EndOfDay(input)
	assert.Equal(t, expected, result, "EndOfDay should return the last nanosecond of the day")
}

func TestIsToday(t *testing.T) {
	loc := time.Local
	now := time.Now().In(loc)
	start := StartOfDay(now)
	end := start.Add(24 * time.Hour)

	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{"ExactlyAtStartOfToday", start, true},
		{"StartOfToday", start.Add(time.Second), true},
		{"ExactlyAtEndOfToday", EndOfDay(start), true},
		{"EndOfToday", end.Add(-time.Second), true},
		{"BeforeToday", start.Add(-time.Second), false},
		{"AfterToday", end.Add(time.Second), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsToday(tt.input), fmt.Sprintf("IsToday(%s) failed", tt.name))
		})
	}
}

func TestIsTodayLoc(t *testing.T) {
	loc := time.FixedZone("TestZone", 3600)
	now := time.Now().In(loc)
	start := StartOfDay(now)
	end := start.Add(24 * time.Hour)

	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{"ExactlyAtStartOfToday", start, true},
		{"StartOfTodayLoc", start.Add(time.Second), true},
		{"ExactlyAtEndOfToday", EndOfDay(start), true},
		{"EndOfTodayLoc", end.Add(-time.Second), true},
		{"BeforeTodayLoc", start.Add(-time.Second), false},
		{"AfterTodayLoc", end.Add(time.Second), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsTodayLoc(tt.input, loc), fmt.Sprintf("IsTodayLoc(%s) failed", tt.name))
		})
	}
}

func TestIsTomorrow(t *testing.T) {
	now := time.Now()
	tomorrowStart := StartOfDay(now.Add(24 * time.Hour))
	dayAfterTomorrowStart := tomorrowStart.Add(24 * time.Hour)

	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{"ExactlyAtTomorrowStart", tomorrowStart, true},
		{"StartOfTomorrow", tomorrowStart.Add(time.Second), true},
		{"ExactlyAtDayAfterTomorrowStart", dayAfterTomorrowStart, false},
		{"EndOfTomorrow", dayAfterTomorrowStart.Add(-time.Second), true},
		{"BeforeTomorrow", tomorrowStart.Add(-time.Second), false},
		{"AfterTomorrow", dayAfterTomorrowStart.Add(time.Second), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsTomorrow(tt.input), "IsTomorrow(%s) failed", tt.name)
		})
	}
}

func TestIsTomorrowLoc(t *testing.T) {
	loc := time.FixedZone("TestZone", 3600)
	now := time.Now().In(loc)
	tomorrowStart := StartOfDay(now.Add(24 * time.Hour))
	dayAfterTomorrow := tomorrowStart.Add(24 * time.Hour)

	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{"StartOfTomorrow", tomorrowStart.Add(time.Second), true},
		{"EndOfTomorrow", dayAfterTomorrow.Add(-time.Second), true},
		{"BeforeTomorrow", tomorrowStart.Add(-time.Second), false},
		{"AfterTomorrow", dayAfterTomorrow.Add(time.Second), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsTomorrowLoc(tt.input, loc), fmt.Sprintf("IsTomorrowLoc(%s) failed", tt.name))
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{"EmptyString", "", 0, true},
		{"OnlyDays", "1d", 24 * time.Hour, false},
		{"InvalidDays", "d2h", 0, true},
		{"OnlyWeeks", "1w", 168 * time.Hour, false},
		{"InvalidWeeks", "w7h", 0, true},
		{"DaysAndHours", "2d3h", 51 * time.Hour, false},
		{"WeeksAndMinutes", "1w30m", 168*time.Hour + 30*time.Minute, false},
		{"DecimalDays", "1.5d", 36 * time.Hour, false},
		{"DecimalWeeks", "0.5w", 84 * time.Hour, false},
		{"MixedUnits", "1d2h30m", 24*time.Hour + 2*time.Hour + 30*time.Minute, false},
		{"MixedOrderUnits", "1d2w2h30m", 2*168*time.Hour + 24*time.Hour + 2*time.Hour + 30*time.Minute, false},
		{"InvalidFormatNoNumber", "d", 0, true},
		{"InvalidUnit", "1x", 0, true},
		{"InvalidDecimalPlacement", ".d", 0, true},
		{"InvalidDuration", "2", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)
			if tt.expectError {
				assert.Error(t, err, "Expected error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Unexpected error for input: %s", tt.input)
				assert.Equal(t, tt.expected, result, "Incorrect duration for input: %s", tt.input)
			}
		})
	}
}

func TestTimeAgo(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"JustNow", now.Add(-0 * time.Second), "just now"},
		{"JustNow", now.Add(-30 * time.Second), "just now"},
		{"MinuteAgo", now.Add(-1 * time.Minute), "1 minute ago"},
		{"MinutesAgo", now.Add(-10 * time.Minute), "10 minutes ago"},
		{"HourAgo", now.Add(-1 * time.Hour), "1 hour ago"},
		{"HoursAgo", now.Add(-3 * time.Hour), "3 hours ago"},
		{"Yesterday", now.Add(-23 * time.Hour), "23 hours ago"},
		{"Yesterday", now.Add(-24 * time.Hour), "yesterday"},
		{"Yesterday", now.Add(-25 * time.Hour), "yesterday"},
		{"DaysAgo", now.Add(-48 * time.Hour), "2 days ago"},
		{"DaysAgo", now.Add(-72 * time.Hour), "3 days ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeAgo(tt.input)
			assert.Equal(t, tt.expected, result, fmt.Sprintf("TimeAgo(%s) failed", tt.name))
		})
	}
}

func TestTimeAgoWithNow(t *testing.T) {
	now := time.Date(2025, time.September, 8, 17, 11, 13, 0, time.UTC)

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"JustNow", now.Add(-0 * time.Second), "just now"},
		{"JustNow", now.Add(-30 * time.Second), "just now"},
		{"MinuteAgo", now.Add(-1 * time.Minute), "1 minute ago"},
		{"MinutesAgo", now.Add(-10 * time.Minute), "10 minutes ago"},
		{"HourAgo", now.Add(-1 * time.Hour), "1 hour ago"},
		{"HoursAgo", now.Add(-3 * time.Hour), "3 hours ago"},
		{"Yesterday", now.Add(-23 * time.Hour), "23 hours ago"},
		{"Yesterday", now.Add(-24 * time.Hour), "yesterday"},
		{"Yesterday", now.Add(-25 * time.Hour), "yesterday"},
		{"DaysAgo", now.Add(-48 * time.Hour), "2 days ago"},
		{"DaysAgo", now.Add(-72 * time.Hour), "3 days ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timeAgoWithNow(tt.input, now)
			assert.Equal(t, tt.expected, result, fmt.Sprintf("TimeAgo(%s) failed", tt.name))
		})
	}
}

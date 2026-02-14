package timeutil

import (
	"time"
)

// Now returns the current time in UTC
func Now() time.Time {
	return time.Now().UTC()
}

// Parse parses a time string and returns it in UTC
func Parse(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

// ParseInLocation parses a time string in the given location and converts to UTC
func ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

// Unix converts a Unix timestamp to UTC time
func Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).UTC()
}

// Today returns the start of today in UTC (00:00:00)
func Today() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// FormatDate formats a time in a standard date format
func FormatDate(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}

// FormatDateTime formats a time in a standard datetime format
func FormatDateTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05")
}

// FormatISO8601 formats a time in ISO 8601 format with timezone
func FormatISO8601(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// MustParse is like Parse but panics on error (useful for constants)
func MustParse(layout, value string) time.Time {
	t, err := Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

// StartOfDay returns the start of the day (00:00:00) for the given time in UTC
func StartOfDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// EndOfDay returns the end of the day (23:59:59) for the given time in UTC
func EndOfDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
}

// AddDays adds the specified number of days to the time
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddYears adds the specified number of years to the time
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// IsExpired checks if the given time is in the past
func IsExpired(t time.Time) bool {
	return t.Before(Now())
}

// DaysUntil returns the number of days until the given time
// Returns negative if the time is in the past
func DaysUntil(t time.Time) int {
	duration := t.Sub(Now())
	return int(duration.Hours() / 24)
}

// ToUTC ensures the time is in UTC timezone
func ToUTC(t time.Time) time.Time {
	return t.UTC()
}

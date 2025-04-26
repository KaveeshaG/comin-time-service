package utils

import (
	"fmt"
	"math"
	"time"
)

// CalculateDuration calculates the duration between two time strings in the format "15:04:05"
func CalculateDuration(startTime, endTime string) (float64, error) {
	start, err := time.Parse("15:04:05", startTime)
	if err != nil {
		return 0, fmt.Errorf("invalid start time format: %w", err)
	}

	end, err := time.Parse("15:04:05", endTime)
	if err != nil {
		return 0, fmt.Errorf("invalid end time format: %w", err)
	}

	// Calculate hours
	duration := end.Sub(start).Hours()

	// Handle overnight shifts (when end time is earlier than start time)
	if duration < 0 {
		duration += 24
	}

	// Round to 2 decimal places
	return math.Round(duration*100) / 100, nil
}

// FormatTime formats a time.Time to a string in the format "15:04:05"
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// ParseDate parses a date string in the format "2006-01-02"
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate formats a time.Time to a string in the format "2006-01-02"
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// IsWeekend checks if the given date is a weekend (Saturday or Sunday)
func IsWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// GetStartOfMonth returns the start date of the month for the given date
func GetStartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

// GetEndOfMonth returns the end date of the month for the given date
func GetEndOfMonth(date time.Time) time.Time {
	return GetStartOfMonth(date).AddDate(0, 1, -1)
}

// GetStartOfWeek returns the start date of the week (Sunday) for the given date
func GetStartOfWeek(date time.Time) time.Time {
	offset := int(date.Weekday())
	return date.AddDate(0, 0, -offset)
}

// GetEndOfWeek returns the end date of the week (Saturday) for the given date
func GetEndOfWeek(date time.Time) time.Time {
	return GetStartOfWeek(date).AddDate(0, 0, 6)
}

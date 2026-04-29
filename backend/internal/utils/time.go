// femabras/backend/internal/utils/time.go
package utils

import "time"

// GetTodayDateString returns the current UTC date as YYYY-MM-DD
func GetTodayDateString() string {
	return time.Now().UTC().Format("2006-01-02")
}

// GetTodayTruncated returns the current UTC time truncated to midnight
func GetTodayTruncated() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}

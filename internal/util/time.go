package util

import (
	"time"
	"errors"
)

func ParseDate(input string) (time.Time, error) {
	layout := "2006-01-02"
	parsed, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	// Check if the date is in the past (before today)
	today := time.Now().Truncate(24 * time.Hour) // remove time portion
	if parsed.Before(today) {
		return time.Time{}, errors.New("deadline cannot be in the past")
	}

	return parsed, nil
}

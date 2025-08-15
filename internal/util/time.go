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
	return parsed, nil
}
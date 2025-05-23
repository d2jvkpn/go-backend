package utils

import (
	// "fmt"
	"time"
)

const (
	MonthOnly    = "2006-01"
	YearOnly     = "2006"
	TimeZoneOnly = "-07:00"
	RFC3339Milli = "2006-01-02T15:04:05.000-07:00"
)

func DateRange(start, end string) (from time.Time, to time.Time, err error) {
	if from, err = time.ParseInLocation(time.DateOnly, start, time.Local); err != nil {
		return from, to, err
	}

	if to, err = time.ParseInLocation(time.DateOnly, end, time.Local); err != nil {
		return from, to, err
	}

	to = to.AddDate(0, 0, 1)

	return from, to, nil
}

func DatesInRange(start, end time.Time) (dates []string) {
	dates = make([]string, 0)

	for end.After(start) {
		dates = append(dates, start.Format(time.DateOnly))
		start = start.AddDate(0, 0, 1)
	}

	return dates
}

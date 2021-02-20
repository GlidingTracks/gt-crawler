package crawler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

/*
	Dates in the crawler are represented as UTC-based time.Time instances.
	The persistence layer stores them as UTC-based UNIX timestamps,
	in seconds from January 1, 1970 UTC.

	Crawling typically is done for a predefined dates period, eg. for a year,
	or for a month. There is a utility method that can be used to generate date
	strings for a particular duration range.

 */



// Timestamp2Date returns a Time object from a given timestamp with a specified offset in days.
// Use daysToAdd to either add or subtract days from the current date.
// For subtraction, use negative numbers.
func Timestamp2Date(timestamp int64, daysToAdd int) time.Time {
	when := time.Unix(timestamp, 0)
	return when.Add(time.Hour * 24 * time.Duration(daysToAdd))
}

// Timestamp2DateString returns a string representation from a given timestamp
// with a specified days offset.
// Same as Timestamp2Date but as string.
func Timestamp2DateString(timestamp int64, addedDays int) string {
	return Time2DateString(Timestamp2Date(timestamp, addedDays))
}

func Time2DateString(when time.Time) string {
	return strconv.Itoa(when.Year()) + "-" + checkZero(int(when.Month())) + "-" + checkZero(when.Day())
}

// ParseDate returns the time.Time instance for a text date in a format YYYY-MM-DD.
func ParseDate(date string) (time.Time, error) {
	parts := strings.Split(date, "-")
	if len(parts) < 3 { return time.Now(), errors.New("wrong ?? in date format, should be YYYY-MM-DD") }
	year, err := strconv.Atoi(parts[0])
	if err != nil { return time.Now(), errors.New("wrong YY in format, should be YYYY-MM-DD") }

	month, err := strconv.Atoi(parts[1])
	if err != nil { return time.Now(), errors.New("wrong MM in date format, should be YYYY-MM-DD") }

	day, err := strconv.Atoi(parts[2])
	if err != nil { return time.Now(), errors.New("wrong DD in date format, should be YYYY-MM-DD") }

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// GenerateDateStrings generates string dates starting from a particular day for a given duration.
// The granularity for generated dates are days, in a format YYYY-MM-DD
func GenerateDateStrings(start time.Time, period time.Duration) []string {
	currentOffset := time.Duration(0)
	result := make([]string, 0, int(period.Hours()) / 24)

	for start.Add(currentOffset).Before(start.Add(period)) {
		result = append(result, Time2DateString(start.Add(currentOffset)))
		currentOffset += time.Hour * 24
	}
	return result
}



// checkZero Utility function to pad date items with 0 if single digit.
func checkZero(d int) (v string) {
	if d < 10 {
		v = "0" + strconv.Itoa(d)
		return
	}

	return strconv.Itoa(d)
}

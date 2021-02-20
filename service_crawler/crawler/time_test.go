package crawler

import (
	"testing"
	"time"
)

func TestGetDate(t *testing.T) {
	t.Run("Should handle adding dates", func(t *testing.T) {
		now := time.Now()
		addedTime := Timestamp2Date(now.Unix(), 2)

		if addedTime.Day() != now.Day()+2 {
			t.Error("Wrong day received")
		}
	})

	t.Run("Should handle subtracting dates", func(t *testing.T) {
		now := time.Now()
		addedTime := Timestamp2Date(now.Unix(),-2)

		if addedTime.Day() != now.Day()-2 {
			t.Error("Wrong day received")
		}
	})
}

func TestGenerateDateStrings(t *testing.T) {
	dates := []string {
		"2018-08-01",
		"2018-08-02",
		"2018-08-03",
		"2018-08-04",
		"2018-08-05",
	}
	startDate := time.Date(2018, time.August, 01, 0, 0, 0, 0, time.UTC)
	gen := GenerateDateStrings(startDate, time.Duration(5) * time.Hour * 24)
	for i, d := range gen {
		if d != dates[i] { t.Errorf("Expected: %s, got: %s", dates[i], d) }
	}
}

package main

import (
	"time"
)

type Clock struct {
	Second int
	Minute int
	Hour   int

	Day   int
	Month int
	Year  int
}

func TimeAgo(t time.Time) (string, error) {

	return "", nil
}

func TimeDiff(t1, t2 time.Time) Clock {
	// set same location
	if t1.Location() != t2.Location() {
		t1 = t1.In(t2.Location())
	}

	if t2.After(t1) {
		t1, t2 = t2, t1
	}

	h1, m1, s1 := t1.Clock()
	h2, m2, s2 := t2.Clock()

	y1, M1, d1 := t1.Date()
	y2, M2, d2 := t2.Date()

	seconds := s1 - s2
	minutes := m1 - m2
	hours := h1 - h2
	days := d1 - d2
	months := int(M1 - M2)
	years := y1 - y2

	if seconds < 0 {
		seconds += 60
		minutes--
	}

	if minutes < 0 {
		minutes += 60
		hours--
	}

	if days < 0 {
		days += time.Date(y1, M1, 0, 0, 0, 0, 0, time.UTC).Day()
		months--
	}

	if months < 0 {
		months += 12
		years--
	}

	return Clock{
		Second: seconds,
		Minute: minutes,
		Hour:   hours,
		Day:    days,
		Month:  months,
		Year:   years,
	}
}

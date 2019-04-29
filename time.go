package frain

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrFutureTime is an error that is thrown when the reference time is greater than
	// the supposedly future time to check
	ErrFutureTime = errors.New("reference time argument is in the future")
)

// Clock basically stores the essential time.Time information excluding milliseconds
type Clock struct {
	Second int
	Minute int
	Hour   int

	Day   int
	Month int
	Year  int
}

// TimeAgo takes a time.Time construct from the past and calculates the difference between it (t) and
// the reference time (ref) and then appending 'ago' at the end
func TimeAgo(ref, t time.Time) (string, error) {
	if ref.After(t) {
		return "", ErrFutureTime
	}

	elapsed := TimeDiff(t, ref)
	unit := "years"
	if elapsed.Year == 0 {
		unit = "months"

		if elapsed.Month == 0 {
			unit = "days"

			if elapsed.Day == 0 {
				unit = "hours"

				if elapsed.Hour == 0 {
					unit = "minutes"

					if elapsed.Minute == 0 {
						unit = "seconds"

						if elapsed.Second == 1 {
							unit = "second"
						}

						return fmt.Sprintf("%d %s ago", elapsed.Second, unit), nil
					}
					if elapsed.Minute == 1 {
						unit = "minute"
					}

					return fmt.Sprintf("%d %s ago", elapsed.Minute, unit), nil
				}
				if elapsed.Hour == 1 {
					unit = "hour"
				}

				return fmt.Sprintf("%d %s ago", elapsed.Hour, unit), nil
			}
			if elapsed.Day == 1 {
				unit = "day"
			}

			return fmt.Sprintf("%d %s ago", elapsed.Day, unit), nil
		}
		if elapsed.Month == 1 {
			unit = "month"
		}

		return fmt.Sprintf("%d %s ago", elapsed.Month, unit), nil

	}
	if elapsed.Year == 1 {
		unit = "year"
	}

	return fmt.Sprintf("%d %s ago", elapsed.Year, unit), nil
}

// TimeDiff finds the time elapsed between two time.Time types
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

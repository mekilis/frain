package frain

import (
	"testing"
	"time"
)

func TestTimeDiff(t *testing.T) {
	tests := []struct {
		u    time.Time
		v    time.Time
		want Clock
	}{
		{
			time.Date(2019, 8, 17, 0, 5, 23, 0, time.UTC),
			time.Date(2019, 8, 18, 12, 19, 47, 0, time.UTC),
			Clock{24, 14, 12, 1, 0, 0},
		},
		{
			time.Date(2020, 2, 28, 11, 30, 55, 0, time.UTC),
			time.Date(2019, 4, 7, 9, 46, 2, 0, time.UTC),
			Clock{53, 44, 1, 21, 10, 0},
		},
		{
			time.Date(2019, 4, 7, 9, 46, 2, 0, time.UTC),
			time.Date(2019, 2, 28, 11, 30, 55, 0, time.UTC),
			Clock{7, 15, 22, 9, 1, 0},
		},
	}

	for _, tt := range tests {
		elapsed := TimeDiff(tt.u, tt.v)
		if elapsed != tt.want {
			t.Errorf("expected %v got %v", tt.want, elapsed)
		}
	}
}

func TestTimeAgo(t *testing.T) {
	tests := []struct {
		createdAt   time.Time
		updatedAt   time.Time
		wantElapsed string
		wantErr     error
	}{
		{
			time.Date(2019, 8, 17, 0, 5, 23, 0, time.UTC),
			time.Date(2019, 8, 18, 12, 19, 47, 0, time.UTC),
			"1 day ago",
			nil,
		},
		{
			time.Date(2019, 8, 18, 12, 19, 47, 0, time.UTC),
			time.Date(2019, 8, 17, 0, 5, 23, 0, time.UTC),
			"",
			ErrFutureTime,
		},
		{
			time.Date(2002, 11, 9, 0, 5, 53, 0, time.UTC),
			time.Date(2019, 11, 9, 0, 1, 24, 0, time.UTC),
			"16 years ago",
			nil,
		},
		{
			time.Date(2020, 7, 23, 6, 5, 53, 0, time.UTC),
			time.Date(2020, 7, 23, 6, 7, 23, 0, time.UTC),
			"1 minute ago",
			nil,
		},
	}

	for _, tt := range tests {
		elapsed, err := TimeAgo(tt.createdAt, tt.updatedAt)
		if elapsed != tt.wantElapsed {
			t.Errorf("expected %v got %v", tt.wantElapsed, elapsed)
		} else if err != tt.wantErr {
			t.Errorf("expected %v got %v", tt.wantErr, err)
		}
	}
}

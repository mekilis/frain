package main

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
		{time.Date(2019, 8, 17, 0, 5, 23, 0, time.UTC), time.Date(2019, 8, 18, 12, 19, 47, 0, time.UTC), Clock{24, 14, 12, 1, 0, 0}},
	}

	for _, tt := range tests {
		elapsed := TimeDiff(tt.u, tt.v)
		if elapsed != tt.want {
			t.Errorf("expected %v got %v", tt.want, elapsed)
		}
	}
}

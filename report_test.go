package frain

import (
	"testing"
)

func TestPad(t *testing.T) {
	tests := []struct {
		s         string
		padlength int
		others    bool
		want      string
	}{
		{"Hello World", 12, true, "Hello Wor..."},
		{"Snakes", 10, true, "Snakes..."},
		{"Riptide", 5, true, "Riptide"},
	}

	for _, tt := range tests {
		if got := pad(tt.s, tt.padlength, tt.others); got != tt.want {
			t.Errorf("expected %v, got %v", tt.want, got)
		}
	}
}

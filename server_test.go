package frain

import (
	"testing"
)

func TestServices(t *testing.T) {
	_, err := Services()
	if err != nil {
		t.Errorf("expected <nil> got %v", err)
	}
}

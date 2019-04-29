package frain

import (
	"fmt"
	"testing"
)

func TestServices(t *testing.T) {
	services, err := Services()
	if err != nil {
		t.Errorf("expected <nil> got %v", err)
	}
	fmt.Println(services)
}

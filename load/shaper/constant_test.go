package shaper

import (
	"testing"
	"time"
)

func TestConstantShaper(t *testing.T) {
	// Create a ConstantShaper with a frequency of 100 Hz
	shaper := NewConstantShaper(100)

	expectedInterval := time.Second / 100
	waitTime, send := shaper.GetNextWaitTime()

	if !send {
		t.Fatalf("Expected send to be true")
	}

	if waitTime != expectedInterval {
		t.Fatalf("Expected %d, got %d", expectedInterval, waitTime)
	}
}

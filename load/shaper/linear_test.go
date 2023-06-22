package shaper

import (
	"testing"
	"time"
)

func TestLinearShaper(t *testing.T) {
	// Create a LinearShaper with a start frequency of 1 Hz and increment of 0.5 Hz
	shaper := NewLinearShaper(1.0, 0.5)

	expectedIntervals := []time.Duration{
		time.Duration(float32(time.Second) / 1.0),
		time.Duration(float32(time.Second) / 1.5),
		time.Duration(float32(time.Second) / 2.0),
	}

	for i := 0; i < 10; i++ {
		t.Errorf("Expected %d, got %d", expectedIntervals[i], shaper.GetNextWaitTime())

		time.Sleep(shaper.GetNextWaitTime())
	}
}

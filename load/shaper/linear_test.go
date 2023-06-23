package shaper

import (
	"testing"
	"time"
)

func TestLinearShaper(t *testing.T) {
	// Create a LinearShaper with a start frequency of 5 Hz and increment of 5 Hz
	shaper := NewLinearShaper(5, 5)

	// Starting frequency of 5 Hz means that the interval is 200 ms
	// So we expect to return 200 ms for the first 5 calls
	for i := 0; i < 5; i++ {
		waitTime := shaper.GetNextWaitTime()
		if waitTime != 200*time.Millisecond {
			t.Errorf("Expected 200 ms, got %v", waitTime)
		}
	}

	// The next call should increase the frequency to 10 Hz
	// So we expect to return 100 ms for the next 10 calls
	for i := 0; i < 10; i++ {
		waitTime := shaper.GetNextWaitTime()
		if waitTime != 100*time.Millisecond {
			t.Errorf("Expected 100 ms, got %v", waitTime)
		}
	}

	// The next call should increase the frequency to 15 Hz
	// So we expect to return 66.67 ms for the next 15 calls
	for i := 0; i < 15; i++ {
		waitTime := shaper.GetNextWaitTime()
		if waitTime != 66*time.Millisecond+667*time.Microsecond {
			t.Errorf("Expected 66.67 ms, got %v", waitTime)
		}
	}

	// The next call should increase the frequency to 20 Hz
	// So we expect to return 50 ms for the next 20 calls
	for i := 0; i < 20; i++ {
		waitTime := shaper.GetNextWaitTime()
		if waitTime != 50*time.Millisecond {
			t.Errorf("Expected 50 ms, got %v", waitTime)
		}
	}
}

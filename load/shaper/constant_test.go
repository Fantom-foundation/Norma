package shaper

import (
	"math"
	"testing"
	"time"
)

func TestConstantShaper(t *testing.T) {
	tests := []struct {
		frequency float64
		duration  time.Duration
		messages  float64
	}{
		// 100 Hz constant load
		{100, 0 * time.Second, 0},
		{100, 1 * time.Second, 100},
		{100, 2 * time.Second, 200},
		{100, 1 * time.Millisecond, 0.1},
		{100, 10 * time.Millisecond, 1},
		{100, 100 * time.Millisecond, 10},
		{100, 250 * time.Millisecond, 25},
		{100, 500 * time.Millisecond, 50},

		// Other frequencies.
		{10, 500 * time.Millisecond, 5},
		{20, 500 * time.Millisecond, 10},
		{14, 500 * time.Millisecond, 7},
		{7, 500 * time.Millisecond, 3.5},
		{0, 500 * time.Millisecond, 0},
		{-1, 500 * time.Millisecond, 0},
	}

	for _, test := range tests {
		shaper := NewConstantShaper(test.frequency)
		got := shaper.GetNumMessagesInInterval(time.Now(), test.duration)
		want := test.messages
		if math.Abs(float64(got-want)) > 1e-6 {
			t.Errorf("incorrect number of messages, wanted %f, got %f", want, got)
		}
	}
}

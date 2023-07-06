package shaper

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestSlopeShaper(t *testing.T) {
	tests := []struct {
		// Shaper properties
		startFrequency float64
		increment      float64
		// Query properties
		from     time.Duration
		to       time.Duration
		expected float64
	}{
		// With start-frequency zero and positive increase
		{0, 1, 0 * time.Second, 1 * time.Second, 0.5},
		{0, 1, 0 * time.Second, 2 * time.Second, 2},
		{0, 1, 0 * time.Second, 3 * time.Second, 4.5},

		{0, 1, 1 * time.Second, 2 * time.Second, 1.5},
		{0, 1, 1 * time.Second, 3 * time.Second, 4},
		{0, 1, 2 * time.Second, 3 * time.Second, 2.5},

		// No increase - constant rate
		{1, 0, 0 * time.Second, 1 * time.Second, 1},
		{2, 0, 0 * time.Second, 1 * time.Second, 2},
		{1, 0, 0 * time.Second, 2 * time.Second, 2},
		{1, 0, 1 * time.Second, 2 * time.Second, 1},

		// With initial frequency + increment
		{1, 1, 0 * time.Second, 1 * time.Second, 1.5},
		{1, 1, 0 * time.Second, 2 * time.Second, 4},

		// With negative increment
		{1, -1, 0 * time.Second, 1 * time.Second, 0.5},
		{1, -1, 0 * time.Second, 2 * time.Second, 0.5},
		{1, -1, 2 * time.Second, 3 * time.Second, 0},

		// With negative start frequency
		{-1, 1, 0 * time.Second, 1 * time.Second, 0},
		{-1, 1, 0 * time.Second, 2 * time.Second, 0.5},
		{-1, 1, 2 * time.Second, 3 * time.Second, 1.5},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("start_frequency=%f,inc=%f,offset=%v,duration=%v",
			test.startFrequency, test.increment, test.from, test.to,
		), func(t *testing.T) {
			shaper := NewSlopeShaper(test.startFrequency, test.increment)

			startTime := time.Now()
			if got, want := shaper.GetNumMessagesInInterval(startTime, time.Duration(0)), float64(0); got != want {
				t.Errorf("failed to initialize shaper, wanted %f, got %f", want, got)
			}

			got := shaper.GetNumMessagesInInterval(startTime.Add(test.from), test.to-test.from)
			want := test.expected

			if math.Abs(float64(got-want)) > 1e-6 {
				t.Errorf("expected number of messages %f, got %f", want, got)
			}
		})
	}
}

package shaper

import (
	"testing"
	"time"
)

func TestWaveShaper(t *testing.T) {
	// Create a WaveShaper with a min frequency of 5 Hz, max frequency of 50 Hz
	// and a wave period of 0.5 second.
	shaper, startTime := initializeWaveShaper(5, 50, 2)

	// The equation for the wave shaper is t(n) = A * sin((2 * pi) / p * n) + B,
	// where:
	// 	- `A` is the amplitude, which is half of the difference between the min and max frequency
	// 	- `p` is the wave period
	// 	- `B` is the vertical shift or the mean value of the wave
	// 	- `n` is the time since start

	tests := []struct {
		second   float32
		expected time.Duration
	}{
		// t(0) = 27.5
		{0, time.Duration(float32(time.Second) / 27.5)},
		// t(0.5) = 50
		{0.5, time.Duration(float32(time.Second) / 50)},
		// t(1) = 27.5
		{1, time.Duration(float32(time.Second) / 27.5)},
		// t(1.5) = 5
		{1.5, time.Duration(float32(time.Second) / 5)},
		// t(2) = 27.5
		{2, time.Duration(float32(time.Second) / 27.5)},
	}

	for _, test := range tests {
		waitTime, send := shaper.GetWaitTimeForTimeStamp(startTime.Add(time.Duration(test.second * float32(time.Second))))
		if !send {
			t.Fatalf("Expected send to be true for second %f", test.second)
		}
		if waitTime != test.expected {
			t.Fatalf("Expected %d, got %d for second %f", test.expected, waitTime, test.second)
		}
	}
}

func TestWaveShaperReturnsNoSendForZeroValue(t *testing.T) {
	// Create a WaveShaper with a min frequency of 5 Hz, max frequency of 50 Hz
	// and a wave period of 0.5 second.
	shaper, startTime := initializeWaveShaper(0, 50, 2)

	// For given parameters, the wave shaper will return 0 for t(1.5).
	_, send := shaper.GetWaitTimeForTimeStamp(startTime.Add(time.Duration(1.5 * float32(time.Second))))
	if send {
		t.Fatalf("Expected send to be false for second 1.5")
	}
}

// initializeWaveShaper initializes a WaveShaper with the given parameters and returns the shaper and the start time.
func initializeWaveShaper(minFrequency, maxFrequency, wavePeriod float32) (*WaveShaper, time.Time) {
	shaper := NewWaveShaper(minFrequency, maxFrequency, wavePeriod)
	startTime := time.Now()
	shaper.startTimeStamp = startTime
	return shaper, startTime
}

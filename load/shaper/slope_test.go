package shaper

import (
	"testing"
	"time"
)

func TestSlopeShaper(t *testing.T) {
	// Create a SlopeShaper with a start frequency of 5 Hz and increment of 5 Hz
	shaper, startTime := initializeSlopeShaper(5, 5)

	// The equation for the slope is t(n) = s + n * i, where `s` is the start frequency,
	// `n` is the time since start and `i` is the increment frequency.

	tests := []struct {
		second   float32
		expected time.Duration
	}{
		// t(0) = 5 + 0 * 5 = 5
		{0, time.Duration(float32(time.Second) / 5)},
		// t(0.2) = 5 + 0.2 * 5 = 6
		{0.2, time.Duration(float32(time.Second) / 6)},
		// t(0.5) = 5 + 0.5 * 5 = 7.5
		{0.5, time.Duration(float32(time.Second) / 7.5)},
		// t(1) = 5 + 1 * 5 = 10
		{1, time.Duration(float32(time.Second) / 10)},
		// t(5) = 5 + 5 * 5 = 30
		{5, time.Duration(float32(time.Second) / 30)},
	}

	for _, test := range tests {
		waitTime, send := shaper.GetWaitTimeForTimeStamp(startTime.Add(time.Duration(test.second * float32(time.Second))))
		if !send {
			t.Fatalf("Expected send to be true for second %f", test.second)
		}
		if waitTime != test.expected {
			t.Fatalf("Expected %d, got %d", test.expected, waitTime)
		}
	}
}

func TestSlopeShaperNegativeIncrementFrequency(t *testing.T) {
	// Create a SlopeShaper with a start frequency of 5 Hz and decrement of 5 Hz
	shaper, startTime := initializeSlopeShaper(50, -5)

	// The equation for the slope is t(n) = s - n * i, where `s` is the start frequency,
	// `n` is the time since start and `i` is the decrement frequency.

	tests := []struct {
		second   float32
		expected time.Duration
	}{
		// t(0) = 50 - 0 * 5 = 50
		{0, time.Duration(float32(time.Second) / 50)},
		// t(0.2) = 50 - 0.2 * 5 = 49
		{0.2, time.Duration(float32(time.Second) / 49)},
		// t(0.5) = 50 - 0.5 * 5 = 47.5
		{0.5, time.Duration(float32(time.Second) / 47.5)},
		// t(1) = 50 - 1 * 5 = 45
		{1, time.Duration(float32(time.Second) / 45)},
		// t(5) = 50 - 5 * 5 = 25
		{5, time.Duration(float32(time.Second) / 25)},
	}

	for _, test := range tests {
		waitTime, send := shaper.GetWaitTimeForTimeStamp(startTime.Add(time.Duration(test.second * float32(time.Second))))
		if !send {
			t.Fatalf("Expected send to be true for second %f", test.second)
		}
		if waitTime != test.expected {
			t.Fatalf("Expected %d, got %d", test.expected, waitTime)
		}
	}
}

func TestSlopeShaperSignalsNoSendOnZeroStartFrequency(t *testing.T) {
	// Create a SlopeShaper with a start frequency of 0 Hz and increment of 5 Hz
	shaper, startTime := initializeSlopeShaper(0, 5)

	// Shaper should signalize no send on second 0
	_, send := shaper.GetWaitTimeForTimeStamp(startTime)
	if send {
		t.Fatal("Expected no send at second 0")
	}

	// Shaper should signalize send on second 1
	_, send = shaper.GetWaitTimeForTimeStamp(startTime.Add(time.Second))
	if !send {
		t.Fatal("Expected send at second 1")
	}
}

func TestSlopeShaperSignalsNoSendOnDecrementedZeroFrequency(t *testing.T) {
	// Create a SlopeShaper with a start frequency of 10 Hz and decrement of 5 Hz
	shaper, startTime := initializeSlopeShaper(10, -5)

	// Shaper should signalize send on second 0
	_, send := shaper.GetWaitTimeForTimeStamp(startTime)
	if !send {
		t.Fatal("Expected send at second 0")
	}

	// Shaper should signalize send on second 1
	_, send = shaper.GetWaitTimeForTimeStamp(startTime.Add(time.Second))
	if !send {
		t.Fatal("Expected send at second 1")
	}

	// Shaper should signalize no send on second 2
	_, send = shaper.GetWaitTimeForTimeStamp(startTime.Add(2 * time.Second))
	if send {
		t.Fatal("Expected no send at second 2")
	}
}

// initializeSlopeShaper initializes a SlopeShaper with the given start frequency and increment frequency.
func initializeSlopeShaper(startFrequency, incrementFrequency float32) (*SlopeShaper, time.Time) {
	shaper := NewSlopeShaper(startFrequency, incrementFrequency)
	startTime := time.Now()
	shaper.startTimeStamp = startTime
	return shaper, startTime
}

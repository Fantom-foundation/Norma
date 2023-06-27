package shaper

import (
	"testing"
	"time"
)

func TestSlopeShaper(t *testing.T) {
	// Create a LinearShaper with a start frequency of 5 Hz and increment of 5 Hz
	shaper := NewSlopeShaper(5, 5)

	// Get the wait time we will use as starting point
	startTime := time.Now()

	// The equation for the slope is t(n) = s + n * i, where `s` is the start frequency,
	// `n` is the time since start and `i` is the increment frequency.

	// Get the wait time for second 0
	// t(0) = 5 + 0 * 5 = 5
	sent, waitTime := shaper.GetWaitTimeForTimeStamps(startTime, startTime)
	if !sent {
		t.Fatal("Expected to send at second 0")
	}
	duration := time.Duration(float32(time.Second) / 5)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 0.2
	// t(0.2) = 5 + 0.2 * 5 = 6
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(200*time.Millisecond))
	if !sent {
		t.Fatal("Expected to send at second 0.2")
	}
	duration = time.Duration(float32(time.Second) / 6)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 0.5
	// t(0.5) = 5 + 0.5 * 5 = 7.5
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(500*time.Millisecond))
	if !sent {
		t.Fatal("Expected to send at second 0.5")
	}
	duration = time.Duration(float32(time.Second) / 7.5)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 1
	// t(1) = 5 + 1 * 5 = 10
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(time.Second))
	if !sent {
		t.Fatal("Expected to send at second 1")
	}
	duration = time.Duration(float32(time.Second) / 10)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 5
	// t(5) = 5 + 5 * 5 = 30
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(5*time.Second))
	if !sent {
		t.Fatal("Expected to send at second 5")
	}
	duration = time.Duration(float32(time.Second) / 30)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}
}

func TestSlopeShaperNegativeIncrementFrequency(t *testing.T) {
	// Create a LinearShaper with a start frequency of 5 Hz and decrement of 5 Hz
	shaper := NewSlopeShaper(50, -5)

	// Get the wait time we will use as starting point
	startTime := time.Now()

	// The equation for the slope is t(n) = s - n * i, where `s` is the start frequency,
	// `n` is the time since start and `i` is the decrement frequency.

	// Get the wait time for second 0
	// t(0) = 50 - 0 * 5 = 50
	sent, waitTime := shaper.GetWaitTimeForTimeStamps(startTime, startTime)
	if !sent {
		t.Fatal("Expected to send at second 0")
	}
	duration := time.Duration(float32(time.Second) / 50)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 0.2
	// t(0.2) = 50 - 0.2 * 5 = 49
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(200*time.Millisecond))
	if !sent {
		t.Fatal("Expected to send at second 0.2")
	}
	duration = time.Duration(float32(time.Second) / 49)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 0.5
	// t(0.5) = 50 - 0.5 * 5 = 47.5
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(500*time.Millisecond))
	if !sent {
		t.Fatal("Expected to send at second 0.5")
	}
	duration = time.Duration(float32(time.Second) / 47.5)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 1
	// t(1) = 50 - 1 * 5 = 45
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(time.Second))
	if !sent {
		t.Fatal("Expected to send at second 1")
	}
	duration = time.Duration(float32(time.Second) / 45)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}

	// Get the wait time for second 5
	// t(5) = 50 - 5 * 5 = 25
	sent, waitTime = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(5*time.Second))
	if !sent {
		t.Fatal("Expected to send at second 5")
	}
	duration = time.Duration(float32(time.Second) / 25)
	if waitTime != duration {
		t.Fatalf("Expected %d, got %d", duration, waitTime)
	}
}

func TestSlopeShaperSignalsNoSendOnZeroStartFrequency(t *testing.T) {
	// Create a LinearShaper with a start frequency of 0 Hz and increment of 5 Hz
	shaper := NewSlopeShaper(0, 5)

	// Get the wait time we will use as starting point
	startTime := time.Now()

	// Shaper should signalize no send on second 0
	send, _ := shaper.GetWaitTimeForTimeStamps(startTime, startTime)
	if send {
		t.Fatal("Expected no send at second 0")
	}

	// Shaper should signalize send on second 1
	send, _ = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(time.Second))
	if !send {
		t.Fatal("Expected send at second 1")
	}
}

func TestSlopeShaperSignalsNoSendOnDecrementedZeroFrequency(t *testing.T) {
	// Create a LinearShaper with a start frequency of 10 Hz and decrement of 5 Hz
	shaper := NewSlopeShaper(10, -5)

	// Get the wait time we will use as starting point
	startTime := time.Now()

	// Shaper should signalize send on second 0
	send, _ := shaper.GetWaitTimeForTimeStamps(startTime, startTime)
	if !send {
		t.Fatal("Expected send at second 0")
	}

	// Shaper should signalize send on second 1
	send, _ = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(time.Second))
	if !send {
		t.Fatal("Expected send at second 1")
	}

	// Shaper should signalize no send on second 2
	send, _ = shaper.GetWaitTimeForTimeStamps(startTime, startTime.Add(2*time.Second))
	if send {
		t.Fatal("Expected no send at second 2")
	}
}

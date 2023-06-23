package shaper

import "time"

// LinearShaper is used to send txs with a linearly increasing frequency
type LinearShaper struct {
	interval           time.Duration
	currentFrequency   float32
	incrementFrequency float32
	currentTick        time.Duration
}

func NewLinearShaper(startFrequency, incrementFrequency float32) *LinearShaper {
	if startFrequency <= 0 {
		startFrequency = 1
	}
	return &LinearShaper{
		currentFrequency:   startFrequency,
		interval:           time.Duration(float32(time.Second) / startFrequency),
		incrementFrequency: incrementFrequency,
		currentTick:        0,
	}
}

// GetNextWaitTime returns the next wait time based on the current frequency
// and the increment frequency.
func (s *LinearShaper) GetNextWaitTime() time.Duration {
	// Increase the current frequency if the current tick is greater than or
	// equal to one second. That means that the current frequency is completed.
	if s.currentTick >= time.Second {
		s.currentFrequency += s.incrementFrequency
		// Round the interval to the nearest microsecond.
		s.interval = time.Duration(float32(time.Second) / s.currentFrequency).Round(time.Microsecond)
		s.currentTick = 0
	}
	// Increase the current tick by the interval.
	s.currentTick += s.interval
	return s.interval
}

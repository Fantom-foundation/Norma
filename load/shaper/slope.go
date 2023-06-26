package shaper

import "time"

// SlopeShaper is used to send txs with a linearly increasing frequency
type SlopeShaper struct {
	interval           time.Duration
	currentFrequency   float32
	incrementFrequency float32
	currentTick        time.Duration
}

func NewSlopeShaper(startFrequency, incrementFrequency float32) *SlopeShaper {
	return &SlopeShaper{
		currentFrequency:   startFrequency,
		interval:           time.Duration(float32(time.Second) / startFrequency),
		incrementFrequency: incrementFrequency,
		currentTick:        0,
	}
}

// GetNextWaitTime returns the next wait time based on the current frequency
// and the increment frequency.
func (s *SlopeShaper) GetNextWaitTime() time.Duration {
	// Increase the current frequency if the current tick is greater than or
	// equal to one second. That means that the current frequency is completed.
	if s.currentTick >= time.Second {
		s.currentFrequency += s.incrementFrequency
		s.interval = time.Duration(float32(time.Second) / s.currentFrequency).Round(time.Microsecond)
		s.currentTick = 0
	}
	// Increase the current tick by the interval.
	s.currentTick += s.interval
	return s.interval
}

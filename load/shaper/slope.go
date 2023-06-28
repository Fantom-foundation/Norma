package shaper

import (
	"math"
	"time"
)

// SlopeShaper is used to send txs with a linearly increasing frequency.
// It is defined as follows:
type SlopeShaper struct {
	startFrequency     float32
	incrementFrequency float32
	// startTimeStamp is the time when the wait time was first obtained.
	startTimeStamp time.Time
}

func NewSlopeShaper(startFrequency, incrementFrequency float32) *SlopeShaper {
	return &SlopeShaper{
		startFrequency:     startFrequency,
		incrementFrequency: incrementFrequency,
	}
}

// GetNextWaitTime returns the next wait time based on the current timestamp
// and start and increment frequency
func (s *SlopeShaper) GetNextWaitTime() (time.Duration, bool) {
	now := time.Now()
	// if this is the first call, set the start time stamp
	if s.startTimeStamp.IsZero() {
		s.startTimeStamp = now
	}
	return s.GetWaitTimeForTimeStamp(now)
}

// GetWaitTimeForTimeStamp returns the wait time for the given timestamp
func (s *SlopeShaper) GetWaitTimeForTimeStamp(current time.Time) (time.Duration, bool) {
	timeSinceStart := current.Sub(s.startTimeStamp).Seconds()

	// if the time since start is 0 and the start frequency is 0, then
	// signal to the consumer that he should ask later again
	// this is special case for the first call to GetNextWaitTime
	if timeSinceStart == 0 && s.startFrequency == 0 {
		// calculate the duration from absolute value (might be negative) of the increment frequency
		return time.Duration(float32(time.Second) / float32(math.Abs(float64(s.incrementFrequency)))), false
	}

	// calculate the current frequency as linear function t(n) = s + n * i,
	// where `s` is the start frequency, `n` is the time since start and `i` is the increment frequency
	currentFrequency := s.startFrequency + float32(timeSinceStart)*s.incrementFrequency

	// if the current frequency decreased to 0, it means, that the increment frequency is negative,
	// and we reached from which we won't send anymore txs, so we return infinity as the wait time
	if currentFrequency <= 0 {
		return time.Duration(math.Inf(1)), false
	}

	// return the wait time for the current frequency
	return time.Duration(float32(time.Second) / currentFrequency), true
}

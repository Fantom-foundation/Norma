package shaper

import (
	"math"
	"time"
)

// WaveShaper is a shaper that sends transactions in waves.
type WaveShaper struct {
	minFrequency float32
	maxFrequency float32
	wavePeriod   float32
	// startTimeStamp is the time when the wait time was first obtained.
	startTimeStamp time.Time
}

func NewWaveShaper(minFrequency, maxFrequency, wavePeriod float32) *WaveShaper {
	return &WaveShaper{
		minFrequency: minFrequency,
		maxFrequency: maxFrequency,
		wavePeriod:   wavePeriod,
	}
}

// GetNextWaitTime returns the next wait time based on the current timestamp
// and start and increment frequency
func (s *WaveShaper) GetNextWaitTime() (time.Duration, bool) {
	now := time.Now()
	// if this is the first call, set the start time stamp
	if s.startTimeStamp.IsZero() {
		s.startTimeStamp = now
	}
	return s.GetWaitTimeForTimeStamp(now)
}

// GetWaitTimeForTimeStamp returns the wait time for the given timestamp
func (s *WaveShaper) GetWaitTimeForTimeStamp(current time.Time) (time.Duration, bool) {
	timeSinceStart := current.Sub(s.startTimeStamp).Seconds()

	// calculate the current frequency as function t(n) = A * sin((2 * pi) / p * n) + B,
	// where:
	// 	- `A` is the amplitude, which is half of the difference between the min and max frequency
	// 	- `p` is the wave period
	// 	- `B` is the vertical shift or the mean value of the wave
	// 	- `n` is the time since start
	currentFrequency := (s.maxFrequency-s.minFrequency)/2*float32(math.Sin((2*math.Pi)/
		float64(s.wavePeriod)*timeSinceStart)) + (s.maxFrequency+s.minFrequency)/2

	// if the current frequency is 0, it means, that min frequency is set to 0, and we reached
	// the bottom of the wave. signal to the consumer that he should ask later again
	// and wait for one second divided by the max frequency (no special reason for that)
	if currentFrequency <= 0 {
		return time.Duration(float32(time.Second) / s.maxFrequency), false
	}

	// return the wait time for the current frequency
	return time.Duration(float32(time.Second) / currentFrequency), true
}

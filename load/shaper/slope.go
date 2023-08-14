package shaper

import (
	"math"
	"time"
)

// SlopeShaper is used to send txs with a linearly increasing frequency.
// It is defined as follows:
type SlopeShaper struct {
	startFrequency     float64
	incrementFrequency float64
	// startTimeStamp is the time when the wait time was first obtained.
	startTimeStamp time.Time
}

func NewSlopeShaper(startFrequency, incrementFrequency float64) *SlopeShaper {
	return &SlopeShaper{
		startFrequency:     startFrequency,
		incrementFrequency: incrementFrequency,
	}
}

func (s *SlopeShaper) Start(start time.Time, info LoadInfoSource) {
	s.startTimeStamp = start
}

// GetNumMessagesInInterval provides the number of messages to be produced
// in the given time interval.
func (s *SlopeShaper) GetNumMessagesInInterval(start time.Time, duration time.Duration) float64 {
	// The number of messages to be sent in the interval is equal to the area
	// under the frequency-time curve given by
	//
	//   f(t) := startFrequency + increment * t = d + k * t
	//
	// The area is equal to the definit integral over the given time range [a,b]
	// which can be computed as
	//
	//   m(a,b) := \int_{a}^{b} k * t + d
	//           = (k*b^2/2 + d*b + c) - (k*a^2/2 + d*a + c)
	//           = k*b^2/2 + d*b - k*a^2/2 - d*a
	//           = k/2 * (b^2-a^2) + d * (b-a)
	//           = incrment/2 * (b^2-a^2) + startFrequency * (b-a)
	//
	// assuming the frequency is always positive in the range [a,b]. Otherwise,
	// the interval needs to be limited to the positive part.

	// Relative begin and end time of interval [a,b].
	a := start.Sub(s.startTimeStamp).Seconds()
	b := a + duration.Seconds()

	// Special case: if the slope is constant (there may be no zero point).
	if s.incrementFrequency == 0 {
		return (b - a) * s.startFrequency
	}

	// the zero point (time at which the frequency is zero)
	z := -s.startFrequency / s.incrementFrequency

	// Restrict range if the zero point is in the interval.
	if s.incrementFrequency > 0 {
		a = math.Max(z, a)
		b = math.Max(z, b)
	} else {
		a = math.Min(z, a)
		b = math.Min(z, b)
	}

	return (s.incrementFrequency/2)*(b*b-a*a) + s.startFrequency*(b-a)
}

package shaper

import (
	"math"
	"time"
)

// ConstantShaper is used to send txs with a constant frequency
type ConstantShaper struct {
	frequency float64
}

func NewConstantShaper(frequency float64) *ConstantShaper {
	return &ConstantShaper{
		frequency: frequency,
	}
}

func (s *ConstantShaper) GetNumMessagesInInterval(start time.Time, duration time.Duration) float64 {
	return math.Max(duration.Seconds()*float64(s.frequency), 0)
}

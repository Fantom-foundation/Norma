package shaper

import "time"

// ConstantShaper is used to send txs with a constant frequency
type ConstantShaper struct {
	interval time.Duration
}

func NewConstantShaper(frequency float32) *ConstantShaper {
	return &ConstantShaper{
		interval: time.Duration(float32(time.Second) / frequency),
	}
}

func (s *ConstantShaper) GetNextWaitTime() (time.Duration, bool) {
	return s.interval, true
}

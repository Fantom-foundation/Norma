package shaper

import "time"

// LinearShaper is used to send txs with a linearly increasing frequency
type LinearShaper struct {
	interval  time.Duration
	increment time.Duration
	lastTime  time.Time
}

func NewLinearShaper(startFrequency, incrementFrequency float32) *LinearShaper {
	return &LinearShaper{
		interval:  time.Duration(float32(time.Second) / startFrequency),
		increment: time.Duration(float32(time.Second) / incrementFrequency),
		lastTime:  time.Now(),
	}
}

func (s *LinearShaper) GetNextWaitTime() time.Duration {
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(s.lastTime)

	waitTime := s.interval - elapsedTime

	s.interval += s.increment
	s.lastTime = currentTime

	return waitTime
}

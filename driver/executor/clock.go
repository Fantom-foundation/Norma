package executor

import (
	"fmt"
	"time"
)

// Clock models time for the execution of a scenario.
type Clock interface {
	// Now obtains the current time, which must progress monotonous. Note, this
	// might not be related to any real world time, but is based on the definition
	// of time as implemented by a given clock.
	Now() Time

	// Suspends execution until the given time (+/- a few milliseconds).
	SleepUntil(Time) error
}

// Time is used to model time in a scenario, relative to the start time. Thus,
// at the start of a simulation the time t=0 is by definition. From there on,
// a clock implementation defines the progress of time (see the Clock interface
// above).
type Time time.Duration

func (t Time) String() string {
	return fmt.Sprintf("%.2f", float64(time.Duration(t).Nanoseconds())/1e9)
}

func Nanoseconds(nanos int64) Time {
	return Time(time.Duration(nanos))
}

func Microseconds(micros int64) Time {
	return Nanoseconds(micros * 1000)
}

func Milliseconds(millis int64) Time {
	return Microseconds(millis * 1000)
}

func Seconds(seconds float32) Time {
	return Nanoseconds(int64(float64(seconds) * 1000 * 1000 * 1000))
}

// SimClock is a simple simulated clock never suspending execution. It can be
// used as a stand-in for other clocks in test cases or dry-run setups.
type SimClock struct {
	now Time
}

func NewSimClock() Clock {
	return &SimClock{}
}

func (c *SimClock) Now() Time {
	return c.now
}

func (c *SimClock) SleepUntil(time Time) error {
	if c.now < time {
		c.now = time
	}
	return nil
}

// WallTimeClock is a clock aiming to follow real wall-clock time. It is
// intended to be used when running scenarios for actual evaluations.
type WallTimeClock struct {
	startTime time.Time
}

func NewWallTimeClock() Clock {
	return &WallTimeClock{time.Now()}
}

func (c *WallTimeClock) Now() Time {
	return Time(time.Since(c.startTime))
}

func (c *WallTimeClock) SleepUntil(deadline Time) error {
	if time.Since(c.startTime) < time.Duration(deadline) {
		time.Sleep(time.Until(c.startTime.Add(time.Duration(deadline))))
	}
	return nil
}

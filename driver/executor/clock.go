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

	// Restart starts ticking of this clock, if the clock was already ticking, it is reset.
	Restart()

	// Suspends execution until the given time (+/- a few milliseconds).
	SleepUntil(Time) error

	// NotifyAt creates a channel sending a message when the given time is reached.
	NotifyAt(Time) <-chan Time
}

// Time is used to model time in a scenario, relative to the start time. Thus,
// at the start of a simulation the time t=0 is by definition. From there on,
// a clock implementation defines the progress of time (see the Clock interface
// above).
type Time time.Duration

func (t Time) String() string {
	return fmt.Sprintf("%.1f", float64(time.Duration(t).Nanoseconds())/1e9)
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

func (c *SimClock) Restart() {
	c.now = 0
}

func (c *SimClock) SleepUntil(time Time) error {
	if c.now < time {
		c.now = time
	}
	return nil
}

func (c *SimClock) NotifyAt(time Time) <-chan Time {
	if c.now < time {
		c.now = time
	}
	ch := make(chan Time, 1)
	ch <- c.now
	return ch
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

func (c *WallTimeClock) Restart() {
	c.startTime = time.Now()
}

func (c *WallTimeClock) SleepUntil(deadline Time) error {
	<-c.NotifyAt(deadline)
	return nil
}

func (c *WallTimeClock) NotifyAt(deadline Time) <-chan Time {
	res := make(chan Time, 1)
	if time.Since(c.startTime) < time.Duration(deadline) {
		// Wait for the required amount of time asynchroniously and then send
		// the notification.
		go func() {
			<-time.After(time.Until(c.startTime.Add(time.Duration(deadline))))
			res <- c.Now()
		}()
	} else {
		// The deadline has passed, no need to wait.
		res <- c.Now()
	}
	return res
}

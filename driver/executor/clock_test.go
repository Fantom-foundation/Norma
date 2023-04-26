package executor

import (
	"testing"
)

type namedClock struct {
	name  string
	clock Clock
}

func getClocks() []namedClock {
	return []namedClock{
		{"SimTime", NewSimClock()},
		{"WallTime", NewWallTimeClock()},
	}
}

func TestClock_NowIsMonotone(t *testing.T) {
	for _, test := range getClocks() {
		t.Run(test.name, func(t *testing.T) {
			clock := test.clock
			now1 := clock.Now()
			now2 := clock.Now()
			if now2 < now1 {
				t.Errorf("time not progressing")
			}

		})
	}
}

func TestClock_SleepSkipsTimeAccurately(t *testing.T) {
	for _, test := range getClocks() {
		t.Run(test.name, func(t *testing.T) {
			clock := test.clock
			start := clock.Now()
			duration := Seconds(0.5)
			clock.SleepUntil(duration)
			end := clock.Now()

			offset := end - (start + duration)
			if offset < Milliseconds(-10) {
				t.Errorf("sleep did not suspend execution long enough, offset: %v", offset)
			}

			if offset > Milliseconds(10) {
				t.Errorf("slept too long, offset: %v", offset)
			}
		})
	}
}

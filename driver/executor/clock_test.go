// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package executor

import (
	"testing"
	"time"
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
			time := Seconds(0.5)
			clock.SleepUntil(time)
			end := clock.Now()

			offset := end - (start + time)
			if offset < Milliseconds(-10) {
				t.Errorf("sleep did not suspend execution long enough, offset: %v", offset)
			}

			if offset > Milliseconds(10) {
				t.Errorf("slept too long, offset: %v", offset)
			}
		})
	}
}

func TestClock_NotifyAtSkipsTimeAccurately(t *testing.T) {
	for _, test := range getClocks() {
		t.Run(test.name, func(t *testing.T) {
			clock := test.clock
			start := clock.Now()
			time := Seconds(0.5)
			weakUpTime := <-clock.NotifyAt(time)
			end := clock.Now()

			offset := end - (start + time)
			if offset < Milliseconds(-10) {
				t.Errorf("sleep did not suspend execution long enough, offset: %v", offset)
			}

			if offset > Milliseconds(10) {
				t.Errorf("slept too long, offset: %v", offset)
			}

			offset = end - weakUpTime
			if offset < 0 {
				t.Errorf("weak-up time received through notification is not consistent with clock")
			}
			if offset > Milliseconds(1) {
				t.Errorf("delay between weak-up time and clock time is to high: %v", offset)
			}
		})
	}
}

func TestClock_Start(t *testing.T) {
	for _, test := range getClocks() {
		t.Run(test.name, func(t *testing.T) {
			clock := test.clock
			t1 := clock.Now()
			// wait and reset - the time diff must be below or zero
			time.Sleep(100 * time.Millisecond)
			clock.Restart()
			t2 := clock.Now()
			diff := time.Duration(t2 - t1)
			if diff > 0 {
				t.Errorf("time has not reset")
			}

		})
	}
}

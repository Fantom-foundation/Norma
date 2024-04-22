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

package shaper

import (
	"math"
	"testing"
	"time"
)

func TestConstantShaper(t *testing.T) {
	tests := []struct {
		frequency float64
		duration  time.Duration
		messages  float64
	}{
		// 100 Hz constant load
		{100, 0 * time.Second, 0},
		{100, 1 * time.Second, 100},
		{100, 2 * time.Second, 200},
		{100, 1 * time.Millisecond, 0.1},
		{100, 10 * time.Millisecond, 1},
		{100, 100 * time.Millisecond, 10},
		{100, 250 * time.Millisecond, 25},
		{100, 500 * time.Millisecond, 50},

		// Other frequencies.
		{10, 500 * time.Millisecond, 5},
		{20, 500 * time.Millisecond, 10},
		{14, 500 * time.Millisecond, 7},
		{7, 500 * time.Millisecond, 3.5},
		{0, 500 * time.Millisecond, 0},
		{-1, 500 * time.Millisecond, 0},
	}

	for _, test := range tests {
		shaper := NewConstantShaper(test.frequency)
		got := shaper.GetNumMessagesInInterval(time.Now(), test.duration)
		want := test.messages
		if math.Abs(float64(got-want)) > 1e-6 {
			t.Errorf("incorrect number of messages, wanted %f, got %f", want, got)
		}
	}
}

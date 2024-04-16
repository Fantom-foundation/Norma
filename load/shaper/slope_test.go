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
	"fmt"
	"math"
	"testing"
	"time"
)

func TestSlopeShaper(t *testing.T) {
	tests := []struct {
		// Shaper properties
		startFrequency float64
		increment      float64
		// Query properties
		from     time.Duration
		to       time.Duration
		expected float64
	}{
		// With start-frequency zero and positive increase
		{0, 1, 0 * time.Second, 1 * time.Second, 0.5},
		{0, 1, 0 * time.Second, 2 * time.Second, 2},
		{0, 1, 0 * time.Second, 3 * time.Second, 4.5},

		{0, 1, 1 * time.Second, 2 * time.Second, 1.5},
		{0, 1, 1 * time.Second, 3 * time.Second, 4},
		{0, 1, 2 * time.Second, 3 * time.Second, 2.5},

		// No increase - constant rate
		{1, 0, 0 * time.Second, 1 * time.Second, 1},
		{2, 0, 0 * time.Second, 1 * time.Second, 2},
		{1, 0, 0 * time.Second, 2 * time.Second, 2},
		{1, 0, 1 * time.Second, 2 * time.Second, 1},

		// With initial frequency + increment
		{1, 1, 0 * time.Second, 1 * time.Second, 1.5},
		{1, 1, 0 * time.Second, 2 * time.Second, 4},

		// With negative increment
		{1, -1, 0 * time.Second, 1 * time.Second, 0.5},
		{1, -1, 0 * time.Second, 2 * time.Second, 0.5},
		{1, -1, 2 * time.Second, 3 * time.Second, 0},

		// With negative start frequency
		{-1, 1, 0 * time.Second, 1 * time.Second, 0},
		{-1, 1, 0 * time.Second, 2 * time.Second, 0.5},
		{-1, 1, 2 * time.Second, 3 * time.Second, 1.5},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("start_frequency=%f,inc=%f,offset=%v,duration=%v",
			test.startFrequency, test.increment, test.from, test.to,
		), func(t *testing.T) {
			shaper := NewSlopeShaper(test.startFrequency, test.increment)

			startTime := time.Now()
			shaper.Start(startTime, nil)

			got := shaper.GetNumMessagesInInterval(startTime.Add(test.from), test.to-test.from)
			want := test.expected

			if math.Abs(float64(got-want)) > 1e-6 {
				t.Errorf("expected number of messages %f, got %f", want, got)
			}
		})
	}
}

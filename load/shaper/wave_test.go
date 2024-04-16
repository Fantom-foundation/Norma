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

func TestWaveShaper(t *testing.T) {
	tests := []struct {
		// Shaper properties
		minFrequency float64
		maxFrequency float64
		wavePeriod   float64
		// Query properties
		from     time.Duration
		to       time.Duration
		expected float64
	}{
		// Half-wave cycle
		{0, 5, 2, time.Duration(0), 1 * time.Second, 2.5},
		// Full-wave cycle
		{0, 5, 2, time.Duration(0), 2 * time.Second, 5},
		// Two full-wave cycles
		{0, 10, 5, time.Duration(0), 10 * time.Second, 50},
		// Start in the middle of a wave cycle with min frequency of 5
		{5, 10, 4, time.Duration(2) * time.Second, 6 * time.Second, 30},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("minFrequency=%f,maxFrequency=%f,wavePeriod=%f,from=%v,to=%v",
			test.minFrequency, test.maxFrequency, test.wavePeriod, test.from, test.to,
		), func(t *testing.T) {
			shaper := &WaveShaper{
				minFrequency:   float32(test.minFrequency),
				maxFrequency:   float32(test.maxFrequency),
				wavePeriod:     float32(test.wavePeriod),
				startTimeStamp: time.Now(),
			}

			shaper.startTimeStamp = time.Now()
			got := shaper.GetNumMessagesInInterval(shaper.startTimeStamp.Add(test.from), test.to-test.from)
			want := test.expected

			if math.Abs(got-want) > 1e-6 {
				t.Errorf("expected number of messages %f, got %f", want, got)
			}
		})
	}
}

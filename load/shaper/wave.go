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
	"time"
)

// WaveShaper is used to send txs with a wave shaped frequency.
// It is defined as follows:
type WaveShaper struct {
	minFrequency float32
	maxFrequency float32
	wavePeriod   float32
	// startTimeStamp is the wall-time when the wave generator got started. At that time, the frequency
	// is the minFrequency and the first period starts.
	startTimeStamp time.Time
}

func NewWaveShaper(minFrequency, maxFrequency, wavePeriod float32) *WaveShaper {
	return &WaveShaper{
		minFrequency: minFrequency,
		maxFrequency: maxFrequency,
		wavePeriod:   wavePeriod,
	}
}

func (w *WaveShaper) Start(start time.Time, _ LoadInfoSource) {
	w.startTimeStamp = start
}

// GetNumMessagesInInterval provides the number of messages to be produced
// in the given time interval.
func (w *WaveShaper) GetNumMessagesInInterval(start time.Time, duration time.Duration) float64 {
	a := float64(w.minFrequency)
	b := float64(w.maxFrequency)
	p := float64(w.wavePeriod)

	// Calculate the relative begin and end time of the interval [x,y].
	x := start.Sub(w.startTimeStamp).Seconds()
	y := x + duration.Seconds()

	// Our function is defined as follows:
	// f(x) = a + (1-cos((2*pi)/p*x))/2 * (b-a)
	// where a is the minimum frequency, b is the maximum frequency and p is the wave period.
	// Integral of our function is: (x*(a+b))/2 + (p*(a-b)*sin((2*pi*x)/p))/(4*pi)
	// We can calculate the integral at the beginning and end of the interval and return the difference.
	integralAtY := (y*(a+b))/2 + (p*(a-b)*math.Sin((2*math.Pi*y)/p))/(4*math.Pi)
	integralAtX := (x*(a+b))/2 + (p*(a-b)*math.Sin((2*math.Pi*x)/p))/(4*math.Pi)

	// Return the difference between the two integrals to get the number of messages in the interval.
	return integralAtY - integralAtX
}

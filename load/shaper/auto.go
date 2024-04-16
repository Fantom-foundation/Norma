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
	"log"
	"time"
)

// autoShaper implements an additive-increase/multiplicative-decrease load control
// algorithm using the gap between submitted and received transactions of an application
// as an overload signal.
//
// See: https://en.wikipedia.org/wiki/Additive_increase/multiplicative_decrease
type autoShaper struct {
	increase          float64 // the additive increase in a non-overload case
	decrease          float64 // the multiplicative decrease in a overload case
	rate              float64 // < the current rate
	lastOverflowCheck time.Time
	loadInfo          LoadInfoSource
}

func NewAutoShaper(increase, decrease float64) Shaper {
	return &autoShaper{
		increase: increase,
		decrease: decrease,
	}
}

func (s *autoShaper) Start(start time.Time, info LoadInfoSource) {
	s.lastOverflowCheck = start
	s.loadInfo = info
}

func (s *autoShaper) GetNumMessagesInInterval(start time.Time, duration time.Duration) float64 {

	// The goal of this shaper is to maximize throughput without creating an overload scenario.
	// To detect overloads, the gap between the submitted and received transactions is tracked.
	// If the gap becomes > twice the current rate per second, the transaction rate is reduced
	// by the configurable `decrease` factor. Otherwise, the transaction rate is increased a
	// configurable `increase` constant.

	// Periodically adjust the transfer rate.
	if start.Sub(s.lastOverflowCheck) >= time.Second {
		s.lastOverflowCheck = start

		// Fetch the latest refresh rates.
		gap := getProcessingGap(s.loadInfo)

		if float64(gap) > 2*s.rate {
			s.rate *= 1 - s.decrease
		} else {
			s.rate += s.increase
		}
	}

	// Use the current rate as the targeted rate.
	return s.rate * duration.Seconds()
}

func getProcessingGap(info LoadInfoSource) uint64 {
	sent, err := info.GetSentTransactions()
	if err != nil {
		log.Printf("autoShaper: failed to fetch number of sent transactions: %v", err)
		return 0
	}
	received, err := info.GetReceivedTransactions()
	if err != nil {
		log.Printf("autoShaper: failed to fetch number of received transactions: %v", err)
		return 0
	}
	return sent - received
}

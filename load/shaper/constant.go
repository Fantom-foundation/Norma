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

// ConstantShaper is used to send txs with a constant frequency
type ConstantShaper struct {
	frequency float64
}

func NewConstantShaper(frequency float64) *ConstantShaper {
	return &ConstantShaper{
		frequency: frequency,
	}
}

func (s *ConstantShaper) Start(time.Time, LoadInfoSource) {}

func (s *ConstantShaper) GetNumMessagesInInterval(start time.Time, duration time.Duration) float64 {
	return math.Max(duration.Seconds()*float64(s.frequency), 0)
}

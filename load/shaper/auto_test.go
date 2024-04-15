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

	"github.com/golang/mock/gomock"
)

func TestAutoShaper_GrowsAdditive(t *testing.T) {
	ctrl := gomock.NewController(t)
	info := NewMockLoadInfoSource(ctrl)

	info.EXPECT().GetSentTransactions().AnyTimes().Return(uint64(120), nil)
	info.EXPECT().GetReceivedTransactions().AnyTimes().Return(uint64(120), nil)

	shaper := NewAutoShaper(10, 0.2)

	start := time.Now()
	shaper.Start(start, info)

	for i := 0; i < 10; i++ {
		want := float64(i * 10)
		if got := shaper.GetNumMessagesInInterval(start, time.Second); math.Abs(got-want) > 1e-6 {
			t.Errorf("invalid number of messages, wanted %f, got %f", want, got)
		}
		start = start.Add(time.Second)
	}
}

func TestAutoShaper_ShrinksMultiplicative(t *testing.T) {
	ctrl := gomock.NewController(t)
	info := NewMockLoadInfoSource(ctrl)

	info.EXPECT().GetSentTransactions().AnyTimes().Return(uint64(100000), nil)
	info.EXPECT().GetReceivedTransactions().AnyTimes().Return(uint64(0), nil)

	rate := 1000.0
	shaper := NewAutoShaper(10, 0.2)
	shaper.(*autoShaper).rate = rate

	start := time.Now()
	shaper.Start(start, info)

	for i := 0; i < 10; i++ {
		want := float64(rate)
		if got := shaper.GetNumMessagesInInterval(start, time.Second); math.Abs(got-want) > 1e-6 {
			t.Errorf("invalid number of messages, wanted %f, got %f", want, got)
		}
		rate *= 0.8
		start = start.Add(time.Second)
	}
}

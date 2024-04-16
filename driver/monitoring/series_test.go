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

package monitoring

import (
	"testing"

	"golang.org/x/exp/slices"
)

type TestBlockSeries struct {
	data []int
}

func (s *TestBlockSeries) GetRange(from, to BlockNumber) []DataPoint[BlockNumber, int] {
	if int(to) > len(s.data) {
		to = BlockNumber(len(s.data))
	}
	if to <= from {
		return nil
	}
	res := make([]DataPoint[BlockNumber, int], 0, to-from)
	for i := from; i < to; i++ {
		res = append(res, DataPoint[BlockNumber, int]{BlockNumber(i), s.data[i]})
	}
	return res
}

func (s *TestBlockSeries) GetLatest() *DataPoint[BlockNumber, int] {
	if len(s.data) == 0 {
		return nil
	}
	pos := len(s.data) - 1
	return &DataPoint[BlockNumber, int]{BlockNumber(pos), s.data[pos]}
}

func (s *TestBlockSeries) SetData(data []int) {
	s.data = make([]int, len(data))
	copy(s.data[:], data[:])
}

func TestTestSeries_IsASeries(t *testing.T) {
	var s TestBlockSeries
	var _ Series[BlockNumber, int] = &s
}

func TestTestSeries_GetRange(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	tests := []struct {
		from, to BlockNumber
		result   []DataPoint[BlockNumber, int]
	}{
		{
			from: 0,
			to:   5,
			result: []DataPoint[BlockNumber, int]{
				{BlockNumber(0), 1},
				{BlockNumber(1), 2},
				{BlockNumber(2), 3},
				{BlockNumber(3), 4},
				{BlockNumber(4), 5},
			},
		},
		{
			from: 3,
			to:   5,
			result: []DataPoint[BlockNumber, int]{
				{BlockNumber(3), 4},
				{BlockNumber(4), 5},
			},
		},
		{
			from:   3,
			to:     2,
			result: nil,
		},
		{
			from:   7,
			to:     10,
			result: nil,
		},
	}

	series := TestBlockSeries{}
	series.SetData(data)
	for _, test := range tests {
		res := series.GetRange(test.from, test.to)
		if !slices.Equal(res, test.result) {
			t.Errorf("invalid result, expected %v, got %v", series.data, res)
		}
	}
}

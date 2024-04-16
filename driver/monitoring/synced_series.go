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
	"fmt"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// SyncedSeries implements a generic series retaining all data in memory and
// offering synchronized access to its content.
type SyncedSeries[K constraints.Ordered, T any] struct {
	data  []DataPoint[K, T]
	mutex sync.Mutex
}

// GetRange extracts a snapshot of a value range of the maintained data.
func (s *SyncedSeries[K, T]) GetRange(from, to K) []DataPoint[K, T] {
	if to < from {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	start := sort.Search(len(s.data), func(i int) bool {
		return from <= s.data[i].Position
	})
	end := sort.Search(len(s.data), func(i int) bool {
		return to <= s.data[i].Position
	})
	if start == end {
		return nil
	}
	res := make([]DataPoint[K, T], end-start)
	copy(res[:], s.data[start:end])
	return res
}

// GetLatest returns the latest collected data point in this series.
func (s *SyncedSeries[K, T]) GetLatest() *DataPoint[K, T] {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.data) == 0 {
		return nil
	}
	res := s.data[len(s.data)-1]
	return &res
}

// Append adds new data to the end of the series. The operation fails if the
// provided point is <= the last added point.
func (s *SyncedSeries[K, T]) Append(point K, value T) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.data) > 0 && s.data[len(s.data)-1].Position >= point {
		return fmt.Errorf("cannot append data out-of-order")
	}
	s.data = append(s.data, DataPoint[K, T]{point, value})
	return nil
}

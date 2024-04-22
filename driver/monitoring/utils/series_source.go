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

package utils

import (
	"fmt"
	"sync"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"golang.org/x/exp/constraints"
)

// SyncedSeriesSource is a base type for all data sources using synced series
// data as their data storage solution.
type SyncedSeriesSource[S comparable, K constraints.Ordered, T any] struct {
	metric   monitoring.Metric[S, monitoring.Series[K, T]]
	data     map[S]*monitoring.SyncedSeries[K, T]
	dataLock sync.Mutex
}

// NewSyncedSeriesSource creates a new data source managing per-subject synced
// series data. Instances are intended to be the base of source implementations.
func NewSyncedSeriesSource[S comparable, K constraints.Ordered, T any](
	metric monitoring.Metric[S, monitoring.Series[K, T]],
) *SyncedSeriesSource[S, K, T] {
	res := &SyncedSeriesSource[S, K, T]{
		metric: metric,
		data:   map[S]*monitoring.SyncedSeries[K, T]{},
	}
	return res
}

func (s *SyncedSeriesSource[S, K, T]) GetMetric() monitoring.Metric[S, monitoring.Series[K, T]] {
	return s.metric
}

func (s *SyncedSeriesSource[S, K, T]) GetSubjects() []S {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	res := make([]S, 0, len(s.data))
	for subject := range s.data {
		res = append(res, subject)
	}
	return res
}

func (s *SyncedSeriesSource[S, K, T]) GetData(subject S) (monitoring.Series[K, T], bool) {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	res, exists := s.data[subject]
	return res, exists
}

func (s *SyncedSeriesSource[S, K, T]) Shutdown() error {
	return nil
}

func (s *SyncedSeriesSource[S, K, T]) ForEachRecord(consumer func(r monitoring.Record)) {
	for subject, series := range s.data {
		r := monitoring.Record{}
		r.SetSubject(subject)

		var first K
		latest := series.GetLatest()
		if latest == nil {
			continue
		}
		allData := series.GetRange(first, latest.Position)
		for _, point := range allData {
			r.SetPosition(point.Position).SetValue(point.Value)
			consumer(r)
		}
		r.SetPosition(latest.Position).SetValue(latest.Value)
		consumer(r)
	}
}

// NewSubject registers a new subject and initiates a data Series for it. The operation fails if the
// same subject is already present. Use this method when you want to make sure that there are no duplicates.
func (s *SyncedSeriesSource[S, K, T]) NewSubject(subject S) (*monitoring.SyncedSeries[K, T], error) {
	// Register a new data series if the subject is new.
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	if _, exist := s.data[subject]; exist {
		return nil, fmt.Errorf("subject %v already present", subject)
	}
	data := &monitoring.SyncedSeries[K, T]{}
	s.data[subject] = data
	return data, nil
}

// GetOrAddSubject looks up a registered subject and creates a new series if the subject has not been
// encountered before. Use this method when it is irrelevant whether the subject has been seen before.
func (s *SyncedSeriesSource[S, K, T]) GetOrAddSubject(subject S) *monitoring.SyncedSeries[K, T] {
	// Register a new data series if the subject is new.
	s.dataLock.Lock()
	if res, exist := s.data[subject]; exist {
		s.dataLock.Unlock()
		return res
	}
	data := &monitoring.SyncedSeries[K, T]{}
	s.data[subject] = data
	s.dataLock.Unlock()
	return data
}

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

func (s *SyncedSeries[K, T]) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.data)
}

func (s *SyncedSeries[K, T]) GetAt(index int) DataPoint[K, T] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.data[index]
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

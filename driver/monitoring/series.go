package monitoring

import (
	"golang.org/x/exp/constraints"
)

// DataPoint is one entry of a data series.
type DataPoint[K constraints.Ordered, T any] struct {
	Position K
	Value    T
}

// Series is a generic interface for arbitrarily indexed sequences of values.
// The type K is the index type, the type T the value associated to the keys.
type Series[K constraints.Ordered, T any] interface {
	// GetRange captures a snapshot of all points collected for the half-open
	// interval [from,to).
	GetRange(from, to K) []DataPoint[K, T]
	// GetLatest retrieves the latest collected data point or nil if no data
	// was collected.
	GetLatest() *DataPoint[K, T]
}

// TimeSeries is a data Series using time-stamps as the index type.
type TimeSeries[T any] interface {
	Series[Time, T]
}

// BlockSeries is a data Series using block numbers as the index type.
type BlockSeries[T any] interface {
	Series[BlockNumber, T]
}

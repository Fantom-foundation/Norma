package monitoring

import "time"

// DataPoint is one entry of a data series.
type DataPoint[K any, T any] struct {
	Position K
	Value    T
}

// Series is a generic interface for arbitrarily indexed sequences of values.
// The type K is the index type, the type T the value associated to the keys.
type Series[K any, T any] interface {
	GetRange(from, to K) []DataPoint[K, T]
}

// TimeSeries is a data Series using time-stamps as the index type.
type TimeSeries[T any] interface {
	Series[time.Time, T]
}

// BlockSeries is a data Series using block numbers as the index type.
type BlockSeries[T any] interface {
	Series[BlockNumber, T]
}

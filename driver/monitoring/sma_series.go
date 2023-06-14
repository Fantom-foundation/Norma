package monitoring

import (
	"golang.org/x/exp/constraints"
	"log"
)

// Number is a constraint type for float and integer types.
type Number interface {
	constraints.Float | constraints.Integer
}

// SmaSeries is computes a Simple Moving Average from the input series and stores it to this series.
// Everytime methods to get data from this series are called, calculation of the average is triggered.
// It is always checked if the input series has grown since and potential missing averages are computed and added to this series.
// The moving average is computed fow a moving window sized to the configured period.
// If the input series is shorter, the window used is smaller and the average computed for this smaller window.
type SmaSeries[K constraints.Ordered, T Number] struct {
	input            Series[K, T]
	output           *SyncedSeries[K, T]
	startKey, endKey K // keep positions of first end last element of the window
	period           int
}

func NewSMASeries[K constraints.Ordered, T Number](input Series[K, T], period int) *SmaSeries[K, T] {
	return &SmaSeries[K, T]{
		input:  input,
		output: &SyncedSeries[K, T]{},
		period: period,
	}
}

func (s *SmaSeries[K, T]) calculate() {
	latest := s.input.GetLatest()
	if latest != nil && (latest.Position > s.endKey || s.endKey == s.startKey) {
		var sum T
		points := append(s.input.GetRange(s.startKey, latest.Position), *latest)

		// This loop moves the window from the start key to the end of the input series.
		// SMA is computed and added to the output series for the moving window.
		// The key marking the start of the rolling window is stored globally to restart
		// rolling from the beginning of the window when this method is repeatedly called.
		for i, point := range points {
			sum += point.Value

			// the window started to move
			// - move key of the start of the window
			// - reduce the sum of the value before the start of the window
			if i >= s.period {
				s.startKey = points[i-s.period].Position
				sum -= points[i-s.period].Value // remove the value before this window start
			}

			// compute and store averages for not already included keys
			// or when the series is empty (i.e. last and end keys are equal)
			if point.Position > s.endKey || s.endKey == s.startKey {
				count := T(s.period)
				if i < s.period {
					count = T(i + 1)
				}
				avr := sum / count
				s.endKey = point.Position
				if err := s.output.Append(s.endKey, avr); err != nil {
					log.Printf("err: %v", err)
				}
			}
		}

	}
}

// GetRange extracts a snapshot of a value range of the maintained output.
func (s *SmaSeries[K, T]) GetRange(from, to K) []DataPoint[K, T] {
	s.calculate()
	return s.output.GetRange(from, to)
}

// GetLatest returns the latest collected output point in this series.
func (s *SmaSeries[K, T]) GetLatest() *DataPoint[K, T] {
	s.calculate()
	return s.output.GetLatest()
}

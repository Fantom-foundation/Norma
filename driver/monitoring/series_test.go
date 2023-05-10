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

func (s *TestBlockSeries) SetData(data []int) {
	s.data = make([]int, len(data))
	copy(s.data[:], data[:])
}

func TestTestSeries_IsABlockSeries(t *testing.T) {
	var s TestBlockSeries
	var _ BlockSeries[int] = &s
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

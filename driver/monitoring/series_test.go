package monitoring

import (
	"testing"

	"golang.org/x/exp/slices"
)

type TestBlockSeries struct {
	data []int
}

func (s *TestBlockSeries) GetRange(from, to Block) []DataPoint[Block, int] {
	if int(to) > len(s.data) {
		to = Block(len(s.data))
	}
	if to <= from {
		return nil
	}
	res := make([]DataPoint[Block, int], 0, to-from)
	for i := from; i < to; i++ {
		res = append(res, DataPoint[Block, int]{Block(i), s.data[i]})
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
		from, to Block
		result   []DataPoint[Block, int]
	}{
		{
			from: 0,
			to:   5,
			result: []DataPoint[Block, int]{
				{Block(0), 1},
				{Block(1), 2},
				{Block(2), 3},
				{Block(3), 4},
				{Block(4), 5},
			},
		},
		{
			from: 3,
			to:   5,
			result: []DataPoint[Block, int]{
				{Block(3), 4},
				{Block(4), 5},
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

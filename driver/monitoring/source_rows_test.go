package monitoring

import "testing"

func TestSourceRows_ForEachRow(t *testing.T) {
	seriesA := &TestBlockSeries{[]int{1, 2}}
	seriesB := &TestBlockSeries{[]int{3, 4, 5}}

	source := TestSource{}
	source.setData("A", seriesA)
	source.setData("B", seriesB)

	expectedRows := []Row[Node, BlockNumber, int, Series[BlockNumber, int]]{
		{TestNodeMetric, "A", 0, 1},
		{TestNodeMetric, "A", 1, 2},
		{TestNodeMetric, "B", 0, 3},
		{TestNodeMetric, "B", 1, 4},
		{TestNodeMetric, "B", 2, 5},
	}

	sr := SourceRowsForEacher[Node, BlockNumber, int, Series[BlockNumber, int]]{&source}
	var i int
	sr.ForEachRow(func(row Row[Node, BlockNumber, int, Series[BlockNumber, int]]) {
		if row != expectedRows[i] {
			t.Errorf("rows do not match: %v != %v", row, expectedRows[i])
		}
		i++
	})

}

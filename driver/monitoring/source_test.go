package monitoring

import (
	"sort"
	"testing"

	"golang.org/x/exp/slices"
)

var (
	TestNodeMetric = Metric[Node, BlockSeries[int]]{
		Name:        "TestNodeMetric",
		Description: "A test metric for unit tests.",
	}
)

// TestSource is a data source providing a stand-in for actual sources in
// tests. This is required since gomock is (yet) not supporting the generation
// of generic mocks.
type TestSource struct {
	data map[Node]BlockSeries[int]
}

func (s *TestSource) GetMetric() Metric[Node, BlockSeries[int]] {
	return TestNodeMetric
}

func (s *TestSource) GetSubjects() []Node {
	res := make([]Node, 0, len(s.data))
	for node := range s.data {
		res = append(res, node)
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return res
}

func (s *TestSource) GetData(node Node) *BlockSeries[int] {
	res := s.data[node]
	if res == nil {
		return nil
	}
	return &res
}
func (s *TestSource) Start() error {
	// Nothing to do.
	return nil
}

func (s *TestSource) Shutdown() error {
	// Nothing to do.
	return nil
}

func (s *TestSource) setData(node Node, data BlockSeries[int]) {
	if s.data == nil {
		s.data = map[Node]BlockSeries[int]{}
	}
	s.data[node] = data
}

func TestTestSourceIsSource(t *testing.T) {
	var source TestSource
	var _ Source[Node, BlockSeries[int]] = &source
}

func TestTestSource_ListsCorrectSubjects(t *testing.T) {
	source := TestSource{}
	want := []Node{}
	if got := source.GetSubjects(); !slices.Equal(got, want) {
		t.Errorf("invalid subject list, wanted %v, got %v", want, got)
	}
	source.setData(Node("A"), &TestBlockSeries{[]int{1, 2, 3}})
	want = []Node{Node("A")}
	if got := source.GetSubjects(); !slices.Equal(got, want) {
		t.Errorf("invalid subject list, wanted %v, got %v", want, got)
	}
	source.setData(Node("B"), &TestBlockSeries{[]int{1}})
	want = []Node{Node("A"), Node("B")}
	if got := source.GetSubjects(); !slices.Equal(got, want) {
		t.Errorf("invalid subject list, wanted %v, got %v", want, got)
	}
}

func TestTestSource_RetrievesCorrectDataSeries(t *testing.T) {
	seriesA := &TestBlockSeries{[]int{1, 2}}
	seriesB := &TestBlockSeries{[]int{3, 4, 5}}

	source := TestSource{}
	source.setData(Node("A"), seriesA)
	source.setData(Node("B"), seriesB)

	if *source.GetData(Node("A")) != seriesA {
		t.Errorf("test source returned wrong series")
	}
	if *source.GetData(Node("B")) != seriesB {
		t.Errorf("test source returned wrong series")
	}
	if source.GetData(Node("C")) != nil {
		t.Errorf("test source returned wrong series")
	}
}

package nodemon

import (
	"math"
	"sort"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
	"golang.org/x/exp/slices"
)

// Unfortunatley, gomock can not (yet) create mocks for generic interfaces.
// So we need to write our own fake sensors for this test.

var (
	testNodeMetric = mon.Metric[mon.Node, mon.TimeSeries[int]]{
		Name:        "TestNodeMetric",
		Description: "A test metric for this unit test.",
	}
)

type testSensor struct {
	next int
}

func (s *testSensor) ReadValue() (int, error) {
	s.next++
	return s.next, nil
}

type testSensorFactory struct{}

func (f *testSensorFactory) CreateSensor(driver.Node) (Sensor[int], error) {
	return &testSensor{}, nil
}

func TestNodeSourceRetrievesSensorData(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1Id := driver.NodeID("A")
	node2Id := driver.NodeID("B")
	node3Id := driver.NodeID("C")

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetNodeID().AnyTimes().Return(node1Id, nil)
	node2.EXPECT().GetNodeID().AnyTimes().Return(node2Id, nil)
	node3.EXPECT().GetNodeID().AnyTimes().Return(node3Id, nil)

	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().UnregisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().Return([]driver.Node{node1, node2})

	source := newPeriodicNodeDataSource[int](testNodeMetric, net, 50*time.Millisecond, &testSensorFactory{})

	// Check that existing nodes are tracked.
	subjects := source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	want := []mon.Node{mon.Node(node1Id), mon.Node(node2Id)}
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	// Simulate the creation of a node after source initialization.
	source.(driver.NetworkListener).AfterNodeCreation(node3)

	// Check that subject list has updated.
	subjects = source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	want = append(want, mon.Node(node3Id))
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	time.Sleep(200 * time.Millisecond)
	if err := source.Shutdown(); err != nil {
		t.Errorf("erros encountered during shutdown: %v", err)
	}

	// Check that subject are still all there.
	subjects = source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	for _, subject := range subjects {
		data := source.GetData(subject)
		if data == nil {
			t.Errorf("no data found for node %s", subject)
			continue
		}
		subrange := (*data).GetRange(mon.Time(0), mon.Time(math.MaxInt64))
		if len(subrange) == 0 {
			t.Errorf("no data collected for node %s", subject)
		}
		for i, point := range subrange {
			if got, want := point.Value, i+1; got != want {
				t.Errorf("unexpected value collected for subject %s: wanted %d, got %d", subject, want, got)
			}
		}
	}
}
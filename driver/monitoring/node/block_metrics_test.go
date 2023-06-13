package nodemon

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
	"io"
	"strings"
	"testing"
	"time"
)

func TestCaptureSeriesFromNodeBlocksNodeMetrics(t *testing.T) {

	ctrl := gomock.NewController(t)
	producer := monitoring.NewMockNodeLogProvider(ctrl)
	producer.EXPECT().RegisterLogListener(gomock.Any()).AnyTimes()

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	writer := monitoring.NewMockWriterChain(ctrl)
	writer.EXPECT().Add(gomock.Any()).AnyTimes()

	source1 := NewBlockTimeSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))
	source2 := NewBlockProcessingTimeSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))

	// simulate data received to metric
	testNodeSource(t, source1)
	testNodeSource(t, source2)
}

func TestIntegrateRegistryWithShutdownNodeMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetNodeID().AnyTimes().Return(driver.NodeID(monitoring.Node1TestId), nil)
	node2.EXPECT().GetNodeID().AnyTimes().Return(driver.NodeID(monitoring.Node2TestId), nil)
	node3.EXPECT().GetNodeID().AnyTimes().Return(driver.NodeID(monitoring.Node3TestId), nil)

	node1.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node1TestLog)), nil)
	node2.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node2TestLog)), nil)
	node3.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node2TestLog)), nil)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	writer := monitoring.NewMockWriterChain(ctrl)
	writer.EXPECT().Add(gomock.Any()).AnyTimes()

	reg := monitoring.NewNodeLogDispatcher(net)
	source := NewBlockTimeSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))
	reg.RegisterLogListener(source)

	// pre-existing node with some blocks
	testNodeSubjects(t, []monitoring.Node{monitoring.Node1TestId}, source)
	testNodeSeriesData(t, monitoring.Node1TestId, monitoring.NodeBlockTestData[monitoring.Node1TestId], source)

	// add second node
	reg.AfterNodeCreation(node2)
	testNodeSubjects(t, []monitoring.Node{monitoring.Node1TestId, monitoring.Node2TestId}, source)
	testNodeSeriesData(t, monitoring.Node2TestId, monitoring.NodeBlockTestData[monitoring.Node2TestId], source)

	// next node will NOT be registered, since the metric is shutdown
	_ = source.Shutdown()
	reg.AfterNodeCreation(node3)
	testNodeSubjects(t, []monitoring.Node{monitoring.Node1TestId, monitoring.Node2TestId}, source)
	// series not created at all
	if _, exists := source.GetData(monitoring.Node3TestId); exists {
		t.Errorf("series shold not exist")
	}
}

// testNodeSubjects tests subjects are present in the source
func testNodeSubjects[T any](t *testing.T, expected []monitoring.Node, source *BlockNodeMetricSource[T]) {
	for _, want := range expected {
		var found bool
		// query the data for some time due to async goroutine
		for i := 0; i < 1000; i++ {
			for _, got := range source.GetSubjects() {
				if got == want {
					found = true
					break
				}
			}
			if found {
				break
			}
			time.Sleep(2 * 10 * time.Millisecond)
		}
		if !found {
			t.Errorf("Node %v not found in: %v", want, source.GetSubjects())
		}
	}

	if len(expected) != len(source.GetSubjects()) {
		t.Errorf("sizes do not match: %v != %v", expected, source.GetSubjects())
	}
}

// testNodeSeriesData verifies if series contains the expected blocks, it queries a few times as data arrive from the goroutine
func testNodeSeriesData[T comparable](t *testing.T, node monitoring.Node, expectedBlocks []monitoring.Block, source *BlockNodeMetricSource[T]) {
	// wait for data for some time due to async goroutine
	// match the series contains expected blocks, loop a few times to let the goroutine provide the data
	for _, want := range expectedBlocks {
		var found bool
		for i := 0; i < 1000; i++ {
			series, exists := source.GetData(node)
			if exists {
				for _, got := range series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000)) {
					if source.getBlockProperty(want) == got.Value {
						found = true
						break
					}
				}
			}
			if found {
				break
			}
			time.Sleep(2 * 10 * time.Millisecond)
		}

		if !found {
			t.Errorf("value: %v not found for block: %v", source.getBlockProperty(want), want)
		}
	}

	// check the size of the series matches the expected blocks
	series, exists := source.GetData(node)
	if !exists {
		t.Fatalf("series should exist")
	}
	if got, want := len(series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000))), len(expectedBlocks); got != want {
		t.Errorf("block series do not match")
	}
}

// testNodeSource tests subjects and series data matches expected constants, defined as globals. The source is filled with expected blocks first
func testNodeSource[T comparable](t *testing.T, source *BlockNodeMetricSource[T]) {

	// insert data into metric
	for node, blocks := range monitoring.NodeBlockTestData {
		for _, block := range blocks {
			source.OnBlock(node, block)
		}
	}

	// check no subject is missing
	for _, node := range source.GetSubjects() {
		if _, exists := monitoring.NodeBlockTestData[node]; !exists {
			t.Errorf("node does not exist: %v", node)
		}
	}

	// check subject length is correct
	if got, want := len(source.GetSubjects()), len(monitoring.NodeBlockTestData); got != want {
		t.Errorf("wrong number of nodes received: %d != %d", got, want)
	}

	// table check results in each series for every node
	for _, node := range source.GetSubjects() {
		series, exists := source.GetData(node)
		if !exists {
			t.Fatalf("series should exist")
		}
		for _, block := range series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000)) {
			if got, want := block.Value, source.getBlockProperty(monitoring.NodeBlockTestData[node][block.Position-1]); got != want {
				t.Errorf("data series contain unexpected value: %v != %v", got, want)
			}
		}
	}
}

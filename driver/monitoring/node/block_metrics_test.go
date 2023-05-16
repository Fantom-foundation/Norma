package nodemon

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

var (
	nodeId1 = monitoring.Node("A")
	nodeId2 = monitoring.Node("B")
	nodeId3 = monitoring.Node("C")

	time1, _ = time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:15.080]")
	time2, _ = time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:15.537]")
	time3, _ = time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:16.027]")
	time4, _ = time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:16.512]")
	time5, _ = time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:17.003]")

	// TODO now the data matches the hardcoded string in the Reader, mock it in next PR
	block1 = monitoring.Block{Height: 2, Time: time1, Txs: 2, GasUsed: 417_928}
	block2 = monitoring.Block{Height: 3, Time: time2, Txs: 1, GasUsed: 117_867}
	block3 = monitoring.Block{Height: 4, Time: time3, Txs: 1, GasUsed: 43426}
	block4 = monitoring.Block{Height: 5, Time: time4, Txs: 5, GasUsed: 138_470}
	block5 = monitoring.Block{Height: 6, Time: time5, Txs: 5, GasUsed: 105_304}

	expected = map[monitoring.Node][]monitoring.Block{
		nodeId1: {block1, block2, block3, block4, block5},
		nodeId2: {block1, block2, block3, block4, block5},
		nodeId3: {block1, block2, block3, block4, block5},
	}
)

func TestCaptureSeriesFromNodeBlocksNodeMetrics(t *testing.T) {

	ctrl := gomock.NewController(t)
	producer := monitoring.NewMockNodeLogProvider(ctrl)
	producer.EXPECT().RegisterLogListener(gomock.Any()).AnyTimes()

	source1 := NewBlockTimeSource(producer)

	// simulate data received to metric
	testNodeSource(t, source1)
}

func TestIntegrateRegistryWithShutdownNodeMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetNodeID().AnyTimes().Return(driver.NodeID(nodeId1), nil)
	node2.EXPECT().GetNodeID().AnyTimes().Return(driver.NodeID(nodeId2), nil)
	node3.EXPECT().GetNodeID().AnyTimes().Return(driver.NodeID(nodeId3), nil)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	reg := monitoring.NewNodeLogDispatcher(net)
	source := NewBlockTimeSource(reg)
	reg.RegisterLogListener(source)

	// pre-existing node with some blocks
	testNodeSubjects(t, []monitoring.Node{nodeId1}, source)
	testNodeSeriesData(t, nodeId1, expected[nodeId1], source)

	// add second node
	reg.AfterNodeCreation(node2)
	testNodeSubjects(t, []monitoring.Node{nodeId1, nodeId2}, source)
	testNodeSeriesData(t, nodeId2, expected[nodeId2], source)

	// next node will NOT be registered, since the metric is shutdown
	_ = source.Shutdown()
	reg.AfterNodeCreation(node3)
	testNodeSubjects(t, []monitoring.Node{nodeId1, nodeId2}, source)
	// series not created at all
	if series := source.GetData(nodeId3); series != nil {
		t.Errorf("series shold not exist")
	}
}

// testNodeSubjects tests subjects are present in the source
func testNodeSubjects[T any](t *testing.T, expected []monitoring.Node, source *BlockNodeMetricSource[T]) {
	for _, want := range expected {
		var found bool
		// query the data for some time due to async goroutine
		for i := 0; i < 100; i++ {
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
		for i := 0; i < 100; i++ {
			series := source.GetData(node)
			if series != nil {
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
	if got, want := len(source.GetData(node).GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000))), len(expectedBlocks); got != want {
		t.Errorf("block series do not match")
	}
}

// testNodeSource tests subjects and series data matches expected constants, defined as globals. The source is filled with expected blocks first
func testNodeSource[T comparable](t *testing.T, source *BlockNodeMetricSource[T]) {

	// insert data into metric
	for node, blocks := range expected {
		for _, block := range blocks {
			source.OnBlock(node, block)
		}
	}

	// check no subject is missing
	for _, node := range source.GetSubjects() {
		if _, exists := expected[node]; !exists {
			t.Errorf("node does not exist: %v", node)
		}
	}

	// check subject length is correct
	if got, want := len(source.GetSubjects()), len(expected); got != want {
		t.Errorf("wrong number of nodes received: %d != %d", got, want)
	}

	// table check results in each series for every node
	for _, node := range source.GetSubjects() {
		for _, block := range source.GetData(node).GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000)) {
			if got, want := block.Value, source.getBlockProperty(expected[node][block.Position-2]); got != want {
				t.Errorf("data series contain unexpected value: %v != %v", got, want)
			}
		}
	}
}

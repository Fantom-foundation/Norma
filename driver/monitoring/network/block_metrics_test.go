package netmon

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
)

func TestCaptureSeriesFromNodeBlocks(t *testing.T) {

	ctrl := gomock.NewController(t)
	producer := monitoring.NewMockNodeLogProvider(ctrl)
	producer.EXPECT().RegisterLogListener(gomock.Any()).AnyTimes()

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}

	source1 := NewNumberOfTransactionsSource(monitor)
	source2 := NewGasUsedSource(monitor)

	// simulate data received to metric
	testNetworkSource(t, source1)
	testNetworkSource(t, source2)
}

func TestIntegrateRegistryWithShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetLabel().AnyTimes().Return(string(monitoring.Node1TestId))
	node2.EXPECT().GetLabel().AnyTimes().Return(string(monitoring.Node2TestId))
	node3.EXPECT().GetLabel().AnyTimes().Return(string(monitoring.Node3TestId))

	node1.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node1TestLog)), nil)
	node2.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node2TestLog)), nil)
	node3.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node3TestLog)), nil)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}

	reg := monitoring.NewNodeLogDispatcher(net)
	source := NewNumberOfTransactionsSource(monitor)
	reg.RegisterLogListener(source)

	// pre-existing node with some blocks
	testNetworkSeriesData(t, monitoring.BlockchainTestData, source)

	// add second node - but we got still the same blockchain
	reg.AfterNodeCreation(node2)
	testNetworkSeriesData(t, monitoring.BlockchainTestData, source)

	// next node will NOT be registered, since the metric is shutdown,
	// but we got the same blockchain as before, i.e. no new blocks
	_ = source.Shutdown()
	reg.AfterNodeCreation(node3)
	testNetworkSeriesData(t, monitoring.BlockchainTestData, source)

	testNetworkSubjects(t, source)
}

// testNetworkSubjects tests subjects are present in the source
func testNetworkSubjects[T any](t *testing.T, source *BlockNetworkMetricSource[T]) {
	// The subject is a constant one network, i.e. no need to wait any incoming subjects
	want := monitoring.Network{}
	var found bool
	for _, got := range source.GetSubjects() {
		if got == want {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Subject %v not found in: %v", want, source.GetSubjects())
	}

	if len(source.GetSubjects()) != 1 {
		t.Errorf("sizes do not match: %v != %v", want, source.GetSubjects())
	}
}

// testNetworkSeriesData verifies series contains expected blocks.
func testNetworkSeriesData[T comparable](t *testing.T, expectedBlocks []monitoring.Block, source *BlockNetworkMetricSource[T]) {
	// wait for data for some time due to async goroutine
	// match the series contains expected blocks, loop a few times to let the goroutine provide the data
	var network monitoring.Network
	for _, want := range expectedBlocks {
		var found bool
		for i := 0; i < 100; i++ {
			series, exists := source.GetData(network)
			if exists {
				for _, got := range series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000)) {
					if source.getBlockProperty(want) == got.Value {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			time.Sleep(2 * 10 * time.Millisecond)
		}

		if !found {
			t.Errorf("value: %v not found for block: %v", source.getBlockProperty(want), want)
		}
	}

	// check the size of the series matches the expected blocks
	series, exists := source.GetData(network)
	if !exists {
		t.Fatalf("series should exist")
	}
	if got, want := len(series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000))), len(expectedBlocks); got != want {
		t.Errorf("block series lengths do not match")
	}
}

// testNetworkSource tests subjects and series data matches expected constants, defined as globals. The source is filled with expected blocks first
func testNetworkSource[T comparable](t *testing.T, source *BlockNetworkMetricSource[T]) {

	// insert data into metric
	for node, blocks := range monitoring.NodeBlockTestData {
		for _, block := range blocks {
			source.OnBlock(node, block)
		}
	}

	// The subject is always one Network
	if got, want := len(source.GetSubjects()), 1; got != want {
		t.Errorf("wrong number of nodes received: %d != %d", got, want)
	}

	// The subject is always one Network
	var mon monitoring.Network
	if got, want := source.GetSubjects()[0], mon; got != want {
		t.Errorf("subject is not a network")
	}

	// table check results
	for _, network := range source.GetSubjects() {
		series, exists := source.GetData(network)
		if !exists {
			t.Fatalf("series should exist")
		}
		for _, block := range series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000)) {
			if got, want := block.Value, source.getBlockProperty(monitoring.BlockchainTestData[block.Position-1]); got != want {
				t.Errorf("data series contain unexpected value: %v != %v", got, want)
			}
		}
	}
}

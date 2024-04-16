// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package nodemon

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
	"io"
	"strings"
	"testing"
)

func TestCaptureSeriesFromNodeBlocksNodeMetrics(t *testing.T) {

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
	source1 := NewBlockTimeSource(monitor)
	source2 := NewBlockProcessingTimeSource(monitor)

	// simulate data received to metric
	testNodeSource(t, source1)
	testNodeSource(t, source2)
}

func TestIntegrateRegistryWithShutdownNodeMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetLabel().AnyTimes().Return(string(monitoring.Node1TestId))
	node2.EXPECT().GetLabel().AnyTimes().Return(string(monitoring.Node2TestId))
	node3.EXPECT().GetLabel().AnyTimes().Return(string(monitoring.Node3TestId))

	node1.EXPECT().StreamLog().AnyTimes().DoAndReturn(func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(monitoring.Node1TestLog)), nil
	})
	node2.EXPECT().StreamLog().AnyTimes().DoAndReturn(func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(monitoring.Node2TestLog)), nil
	})
	node3.EXPECT().StreamLog().AnyTimes().DoAndReturn(func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(monitoring.Node3TestLog)), nil
	})

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}

	source := NewBlockTimeSource(monitor)
	reg := monitor.NodeLogProvider().(*monitoring.NodeLogDispatcher)

	// add first node
	reg.AfterNodeCreation(node1)
	reg.WaitForLogsToBeConsumed()

	// add first node
	testNodeSubjects(t, []monitoring.Node{monitoring.Node1TestId}, source)
	testNodeSeriesData(t, monitoring.Node1TestId, monitoring.NodeBlockTestData[monitoring.Node1TestId], source)

	// add second node
	reg.AfterNodeCreation(node2)
	reg.WaitForLogsToBeConsumed()
	testNodeSubjects(t, []monitoring.Node{monitoring.Node1TestId, monitoring.Node2TestId}, source)
	testNodeSeriesData(t, monitoring.Node2TestId, monitoring.NodeBlockTestData[monitoring.Node2TestId], source)

	// next node will NOT be registered, since the metric is shutdown
	_ = source.Shutdown()
	reg.AfterNodeCreation(node3)
	reg.WaitForLogsToBeConsumed()
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
		for _, got := range source.GetSubjects() {
			if got == want {
				found = true
				break
			}
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
	// match the series contains expected blocks, loop a few times to let the goroutine provide the data
	for _, want := range expectedBlocks {
		var found bool
		series, exists := source.GetData(node)
		if exists {
			for _, got := range series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000)) {
				if source.getBlockProperty(want) == got.Value {
					found = true
					break
				}
			}
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

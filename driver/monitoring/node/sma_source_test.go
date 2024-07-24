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
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"go.uber.org/mock/gomock"
)

func TestSMASource(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().Return([]driver.Node{}).AnyTimes()

	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}

	source := NewTransactionsThroughputSource(monitor)
	sf := sourceFactory[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]]{TransactionsThroughput, source}
	if err := monitoring.InstallSource[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]](monitor, &sf); err != nil {
		t.Fatalf("failed to install source: %v", err)
	}

	period := 2
	smaFactory := func(input monitoring.Series[monitoring.BlockNumber, float32]) monitoring.Series[monitoring.BlockNumber, float32] {
		return monitoring.NewSMASeries[monitoring.BlockNumber, float32](input, period)
	}
	TransactionThroughputSMA := monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]]{
		Name:        fmt.Sprintf("TransactionThroughputSMA_%d", period),
		Description: "Transaction throughput standard moving average",
	}

	smaSource := newNodeBlockSeriesTransformation(monitor, TransactionThroughputSMA, TransactionsThroughput, smaFactory)

	// fill in original source with data
	for node, blocks := range monitoring.NodeBlockTestData {
		for _, block := range blocks {
			source.OnBlock(node, block)
		}
	}

	expected := map[monitoring.Node][]float32{
		monitoring.Node1TestId: {43.763676, 52.494083},
		monitoring.Node2TestId: {40.733196},
		monitoring.Node3TestId: {},
	}

	// test SMA computed for each node
	for node := range monitoring.NodeBlockTestData {
		series, exists := smaSource.GetData(node)
		if !exists {
			t.Fatalf("series does not exist for subject: %v", node)
		}

		points := series.GetRange(monitoring.BlockNumber(0), monitoring.BlockNumber(1000))
		for i, block := range points {
			if got, want := block.Value, expected[node][i]; got != want {
				t.Errorf("data series contain unexpected value for: %v: %v != %v", node, got, want)
			}
		}

		if got, want := len(points), len(expected[node]); got != want {
			t.Errorf("number of points does not mathc: %d != %d", got, want)
		}
	}

	// test subjects present
	for _, node := range smaSource.GetSubjects() {
		if _, exists := monitoring.NodeBlockTestData[node]; !exists {
			t.Errorf("subject %v is not present", node)
		}
	}

	if got, want := len(smaSource.GetSubjects()), len(monitoring.NodeBlockTestData); got != want {
		t.Errorf("number of subjects does not mathc: %d != %d", got, want)
	}
}

func TestSMACsvExport(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	config := monitoring.MonitorConfig{OutputDir: t.TempDir()}
	monitor, err := monitoring.NewMonitor(net, config)
	if err != nil {
		t.Fatalf("failed to start monitor instance: %v", err)
	}
	source := NewTransactionsThroughputSource(monitor)
	sf := sourceFactory[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]]{TransactionsThroughput, source}
	if err := monitoring.InstallSource[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]](monitor, &sf); err != nil {
		t.Fatalf("failed to install source: %v", err)
	}

	period := 2
	smaFactory := func(input monitoring.Series[monitoring.BlockNumber, float32]) monitoring.Series[monitoring.BlockNumber, float32] {
		return monitoring.NewSMASeries[monitoring.BlockNumber, float32](input, period)
	}
	TransactionThroughputSMA := monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]]{
		Name:        fmt.Sprintf("TransactionThroughputSMA_%d", period),
		Description: "Transaction throughput standard moving average",
	}

	smaSource := newNodeBlockSeriesTransformation(monitor, TransactionThroughputSMA, TransactionsThroughput, smaFactory)
	sf = sourceFactory[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]]{TransactionThroughputSMA, smaSource}
	if err := monitoring.InstallSource[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]](monitor, &sf); err != nil {
		t.Fatalf("failed to install source: %v", err)
	}

	seconds := time.Now().Unix()
	// time diff only 50ns
	source.OnBlock("A", monitoring.Block{Height: 10, Time: time.Unix(seconds, 0), Txs: 10})
	source.OnBlock("A", monitoring.Block{Height: 11, Time: time.Unix(seconds+1, 0), Txs: 10})

	if err := monitor.Shutdown(); err != nil {
		t.Fatalf("failed to shut down monitoring: %v", err)
	}

	content, _ := os.ReadFile(monitor.GetMeasurementFileName())

	lines := []string{
		"TransactionsThroughput, network, A, , , 11, , 10\n",
		"TransactionThroughputSMA_2, network, A, , , 11, , 10\n",
	}

	csv := string(content)
	for _, line := range lines {
		if !strings.Contains(csv, line) {
			t.Errorf("missing line %v, got %v", line, csv)
		}
	}
}

type sourceFactory[S any, T any] struct {
	metric monitoring.Metric[S, T]
	source monitoring.Source[S, T]
}

func (f *sourceFactory[S, T]) GetMetric() monitoring.Metric[S, T] {
	return f.metric
}

func (f *sourceFactory[S, T]) CreateSource(*monitoring.Monitor) monitoring.Source[S, T] {
	return f.source
}

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
	"sort"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	opera "github.com/Fantom-foundation/Norma/driver/node"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slices"
)

func TestLogsAddedToSeries(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}

	source := NewPromLogSource(monitor, monitoring.NewPrometheusNameKey("metric_name_a"))

	requestedItems := 1000
	expectedTimes := make([]monitoring.Time, 0, requestedItems)
	expectedValues := make([]float64, 0, requestedItems)

	for i := 0; i < requestedItems; i++ {
		expectedTimes = append(expectedTimes, monitoring.Time(i))
		expectedValues = append(expectedValues, float64(i))
		source.OnLog("A", monitoring.Time(i), float64(i))
		source.OnLog("B", monitoring.Time(2*i), float64(2*i))
	}

	// test subjects
	want := []monitoring.Node{"A", "B"}
	subjects := source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	for i, multiplier := range []int{1, 2} {
		t.Run(fmt.Sprintf("subject %s", want[i]), func(t *testing.T) {
			// test series content
			series, exists := source.GetData(want[i])
			if !exists {
				t.Errorf("series should exist")
			}
			seriesRange := series.GetRange(monitoring.Time(0), monitoring.Time(2*requestedItems))
			if got, want := len(seriesRange), requestedItems; got != want {
				t.Errorf("unexpected series size: %d != %d", got, want)
			}
			for i, datapoint := range seriesRange {
				if got, want := datapoint.Position, monitoring.Time(i*multiplier); got != want {
					t.Errorf("expected position does not match: %d != %d", got, want)
				}
				if got, want := int(datapoint.Value), i*multiplier; got != want {
					t.Errorf("expected position does not match: %d != %d", got, want)
				}
			}
		})
	}
}

func TestLogsIntegrationGetRealMetric(t *testing.T) {
	t.Cleanup(SuppressVerboseLog())

	t.Parallel()
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	client, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})
	node, err := opera.StartOperaDockerNode(client, nil, &opera.OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	t.Cleanup(func() {
		_ = node.Cleanup()
	})

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node})

	outDir := t.TempDir()
	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: outDir})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}
	t.Cleanup(func() {
		_ = monitor.Shutdown()
	})

	for _, metricKeys := range []monitoring.PrometheusLogKey{
		monitoring.NewPrometheusKey("chain_execution", monitoring.Quantile05),
		monitoring.NewPrometheusNameKey("chain_execution_count")} {

		// wait for the metric to arrive for some time
		source := NewPromLogSource(monitor, metricKeys)
		var found bool
		for i := 0; i < 100; i++ {
			series, exists := source.GetData("test")
			if exists {
				datapoint := series.GetLatest()
				if datapoint != nil {
					found = datapoint.Position > 0
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
		}

		if !found {
			t.Errorf("Metric data not arrived within give time for: %s", metricKeys.Name)
		}
	}
}

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

package monitoring

import (
	"bytes"
	"os"
	"sync"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	opera "github.com/Fantom-foundation/Norma/driver/node"
	"github.com/golang/mock/gomock"
	"golang.org/x/exp/slices"
)

func TestMonitor_CreateAndShutdown(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	monitor, err := NewMonitor(net, MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}
	if err := monitor.Shutdown(); err != nil {
		t.Errorf("shutdown of empty monitor failed: %v", err)
	}
}

func TestMonitor_RegisterAndRetrievalOfDataWorks(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	seriesA := &TestBlockSeries{[]int{1, 2}}
	seriesB := &TestBlockSeries{[]int{3, 4, 5}}

	source := TestSource{}
	source.setData("A", seriesA)
	source.setData("B", seriesB)

	metric := source.GetMetric()

	monitor, err := NewMonitor(net, MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}
	if IsSupported(monitor, metric) {
		t.Errorf("empty monitor should not support any metric")
	}

	if subjects := GetSubjects(monitor, metric); len(subjects) != 0 {
		t.Errorf("empty monitor should not report available subjects")
	}

	factory := &genericSourceFactory[Node, Series[BlockNumber, int]]{
		TestNodeMetric,
		func(*Monitor) Source[Node, Series[BlockNumber, int]] { return &source },
	}
	InstallSource[Node, Series[BlockNumber, int]](monitor, factory)

	if !IsSupported(monitor, metric) {
		t.Errorf("registered metric is not supported")
	}

	want := []Node{Node("A"), Node("B")}
	if got := GetSubjects(monitor, metric); !slices.Equal(want, got) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, got)
	}

	if series, exists := GetData(monitor, Node("A"), metric); !exists || series != seriesA {
		t.Errorf("obtained wrong data for node 1")
	}

	if series, exists := GetData(monitor, Node("B"), metric); !exists || series != seriesB {
		t.Errorf("obtained wrong data for node 2")
	}

	if series, exists := GetData(monitor, Node("C"), metric); exists || series != nil {
		t.Errorf("should not have obtained any data for node 3")
	}

	if err := monitor.Shutdown(); err != nil {
		t.Errorf("failed to shutdown monitor: %v", err)
	}
}

func TestMonitor_CsvExport(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	outDir := t.TempDir()
	monitor, err := NewMonitor(net, MonitorConfig{OutputDir: outDir})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}

	_ = monitor.Shutdown()

	content, _ := os.ReadFile(monitor.GetMeasurementFileName())

	buffer := new(bytes.Buffer)
	WriteCsvHeader(buffer)
	if got, want := string(content), buffer.String(); got != want {
		t.Errorf("unexpected export: %v != %v", got, want)
	}
}

func TestMonitorPrometheusLogProviderConfigured(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	outDir := t.TempDir()
	monitor, err := NewMonitor(net, MonitorConfig{OutputDir: outDir})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}
	defer func() {
		_ = monitor.Shutdown()
	}()

	if monitor.PrometheusLogProvider() == nil {
		t.Errorf("prometheus log provider not configured")
	}
}

func TestMonitorIntegrationPrometheusLogReceived(t *testing.T) {
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
	monitor, err := NewMonitor(net, MonitorConfig{OutputDir: outDir})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}
	defer func() {
		_ = monitor.Shutdown()
	}()

	// check metric arrived
	once := sync.Once{}
	done := make(chan bool)
	listener := NewMockTimeLogListener(ctrl)
	// expected to be called
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Do(func(Node, Time, float64) {
		once.Do(func() { close(done) })
	})
	monitor.PrometheusLogProvider().RegisterLogListener(NewPrometheusNameKey("chain_block_age"), listener)

	<-done
}

func TestMonitor_GetPrometheusLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	outDir := t.TempDir()
	monitor, err := NewMonitor(net, MonitorConfig{OutputDir: outDir})
	if err != nil {
		t.Fatalf("failed to create monitor instance: %v", err)
	}
	t.Cleanup(func() {
		_ = monitor.Shutdown()
	})

	if monitor.PrometheusLogProvider() == nil {
		t.Errorf("prometheus log provider not configured")
	}
}

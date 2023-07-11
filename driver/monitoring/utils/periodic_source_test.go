package utils

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
	"sync/atomic"
	"testing"
	"time"
)

func TestPeriodicSourceShutdownBeforeAnyAction(t *testing.T) {
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

	testMetric := monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.Time, int]]{
		Name:        "TestMetric",
		Description: "Test Metric",
	}

	source := NewPeriodicDataSource[monitoring.Node, int](testMetric, monitor)
	if err := source.Shutdown(); err != nil {
		t.Errorf("error to shutdown: %s", err)
	}
}

func TestPeriodicSourceShutdownSourcesAdded(t *testing.T) {
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

	testMetric := monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.Time, int]]{
		Name:        "TestMetric",
		Description: "Test Metric",
	}

	source := NewPeriodicDataSource[monitoring.Node, int](testMetric, monitor)

	var node monitoring.Node
	if err := source.AddSubject(node, &testSensor{}); err != nil {
		t.Errorf("error to add subject: %s", err)
	}

	series, exists := source.GetData(node)
	if !exists {
		t.Fatalf("series should exist")
	}
	// wait for data
	var found bool
	for !found {
		if series.GetLatest() != nil {
			found = true
		}
		time.Sleep(100 * time.Millisecond)
	}

	_ = source.Shutdown()
}

func TestPeriodicSourceErrors(t *testing.T) {
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

	testMetric := monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.Time, int]]{
		Name:        "TestMetric",
		Description: "Test Metric",
	}

	sensor := &buggySensor{}
	source := NewPeriodicDataSourceWithPeriod[monitoring.Node, int](testMetric, monitor, 1*time.Nanosecond)

	var node monitoring.Node
	if err := source.AddSubject(node, sensor); err != nil {
		t.Errorf("error to add subject: %s", err)
	}

	// wait for sensor called many times
	for sensor.count() < 5 {
		time.Sleep(1 * time.Millisecond)
	}

	if err := source.Shutdown(); err == nil {
		t.Errorf("shutdown should return an error")
	}
}

type testSensor struct {
	counts atomic.Int32
}

func (s *testSensor) ReadValue() (int, error) {
	s.counts.Add(1)
	return 123, nil
}

func (s *testSensor) count() int {
	return int(s.counts.Load())
}

type buggySensor struct {
	testSensor
}

func (s *buggySensor) ReadValue() (int, error) {
	s.counts.Add(1)
	return 123, fmt.Errorf("buggy senzor")
}

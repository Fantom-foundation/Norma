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

func TestPeriodicSourceSourceRemoved(t *testing.T) {
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

	sensor1 := &testSensor{}
	node1 := monitoring.Node("A")
	if err := source.AddSubject(node1, sensor1); err != nil {
		t.Errorf("error to add subject: %s", err)
	}

	sensor2 := &testSensor{}
	node2 := monitoring.Node("B")
	if err := source.AddSubject(node2, sensor2); err != nil {
		t.Errorf("error to add subject: %s", err)
	}

	// wait for sensor called a few times
	for sensor1.count() < 5 {
		time.Sleep(100 * time.Millisecond)
	}

	if err := source.RemoveSubject(node1); err != nil {
		t.Errorf("error to remove subject: %s", err)
	}

	// sensor2 should keep increasing while the other one not
	start1 := sensor1.count()
	start2 := sensor2.count()
	for sensor2.count() < 5+start2 {
		time.Sleep(100 * time.Millisecond)
		if sensor1.count() != start1 {
			t.Errorf("subject which was removes keeps being updated: %d != %d", sensor1.count(), start1)
		}
	}
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

package monitoring

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/golang/mock/gomock"
	"golang.org/x/exp/slices"
)

func TestMonitor_CreateAndShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	monitor := NewMonitor(net, MonitorConfig{})
	if err := monitor.Shutdown(); err != nil {
		t.Errorf("shutdown of empty monitor failed: %v", err)
	}
}

func TestMonitor_RegisterAndRetrievalOfDataWorks(t *testing.T) {
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

	monitor := NewMonitor(net, MonitorConfig{})
	if IsSupported(monitor, metric) {
		t.Errorf("empty monitor should not support any metric")
	}

	if subjects := GetSubjects(monitor, metric); len(subjects) != 0 {
		t.Errorf("empty monitor should not report available subjects")
	}

	factory := &genericSourceFactory[Node, BlockSeries[int]]{
		TestNodeMetric,
		func(*Monitor) Source[Node, BlockSeries[int]] { return &source },
	}
	InstallSource[Node, BlockSeries[int]](monitor, factory)

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

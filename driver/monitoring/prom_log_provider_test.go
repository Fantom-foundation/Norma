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
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/uber-go/mock/gomock"
)

func TestLogsDispatchedImplements(t *testing.T) {
	var inst PrometheusLogDispatcher
	var _ PrometheusLogProvider = &inst
}

func TestLogsDispatched(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	aUrl := driver.URL("A")
	bUrl := driver.URL("B")
	cUrl := driver.URL("C")

	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)
	node2.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&bUrl)
	node3.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&cUrl)

	node1.EXPECT().GetLabel().AnyTimes().Return("A")
	node2.EXPECT().GetLabel().AnyTimes().Return("B")
	node3.EXPECT().GetLabel().AnyTimes().Return("C")

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1, node2})

	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		return testData[url], nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Second, testFunc)
	defer dispatcher.Shutdown()

	dispatcher.AfterNodeCreation(node3)

	var wg sync.WaitGroup
	wg.Add(6)
	done := func(node Node, timestamp Time, value float64) {
		wg.Done()
	}
	listener := NewMockTimeLogListener(ctrl)
	// expected logs distinguished by unique values for simplicity
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(1)).Do(done)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(2)).Do(done)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(3)).Do(done)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(8)).Do(done)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(9)).Do(done)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(10)).Do(done)

	// register only for metrics
	// A - summary 0.99 quantile, counter and gauge
	// C
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_summary", "0.99"}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_counter", "0"}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_gauge", "0"}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"C", "0"}, listener)

	wg.Wait()
}

func TestLogsDispatchedLogsOrdered(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)
	aUrl := driver.URL("A")
	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)
	node1.EXPECT().GetLabel().AnyTimes().Return("A")

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	requestedItems := 100

	var prevVal atomic.Int64
	var wg sync.WaitGroup
	wg.Add(1)
	listener := NewMockTimeLogListener(ctrl)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Do(func(node Node, timestamp Time, value float64) {
		if int64(value) <= prevVal.Load() {
			t.Errorf("values not ordered: %d <= %d", int64(value), prevVal.Load())
		}
		if int(value) >= requestedItems {
			wg.Done()
		}
		prevVal.Store(int64(value))
	})

	var counter atomic.Int64
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		res := PrometheusLogValue{PrometheusLogKey{"A", "0"}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}

	dispatcher := newPrometheusLogDispatcher(net, 10*time.Millisecond, testFunc)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", "0"}, listener)
	defer dispatcher.Shutdown()

	wg.Wait()
}

func TestLogsDispatchedShutdown(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)
	aUrl := driver.URL("A")
	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)
	node1.EXPECT().GetLabel().AnyTimes().Return("A")

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	requestedItems := 100

	var counter atomic.Int64
	done := make(chan bool, 1)
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		if int(next) == requestedItems {
			done <- true
		}
		res := PrometheusLogValue{PrometheusLogKey{"A", "0"}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 10*time.Millisecond, testFunc)
	listener := NewMockTimeLogListener(ctrl)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(requestedItems)

	// listen for metrics
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", "0"}, listener)

	// wait for first 1000 values
	<-done
	dispatcher.Shutdown()
	// fetch real counts after shutdown as there could be execution between reading from the channel and shutdown.
	sizeAfterShutdown := counter.Load()

	// make sure there is some delay after Shutdown
	// the mocked listener cannot be called within this time
	time.Sleep(2 * time.Second)

	if got, want := counter.Load(), sizeAfterShutdown; got > want {
		t.Errorf("more data has accunulated after shutdown: %d != %d", got, want)
	}
}

func TestLogsDispatchedUnregisterListener(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)
	aUrl := driver.URL("A")
	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)
	node1.EXPECT().GetLabel().AnyTimes().Return("A")

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	requestedItems := 1000

	listener := NewMockTimeLogListener(ctrl)
	// altering A, B metrics denoted by continuously growing 1 to N
	calls := make([]*gomock.Call, 0, 2*requestedItems)
	for i := 1; i <= requestedItems; i++ {
		calls = append(calls, listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(i)))
	}
	// only B metric denoted by multiplies of 2
	for i := requestedItems; i < 2*requestedItems; i = i + 2 {
		calls = append(calls, listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(i+1)))
	}

	// check all called in order
	gomock.InOrder(calls...)

	metrics := []string{"A", "B"}
	var mu sync.Mutex // use lock to either run logs processing loop or run checking the processing is at the end
	var counter atomic.Int64
	current := make(chan bool, 1)
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		mu.Lock()
		next := counter.Add(1)
		current <- true
		res := PrometheusLogValue{PrometheusLogKey{metrics[next%2], "0"}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)
	defer dispatcher.Shutdown()

	// listen for two metrics A and B
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", "0"}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"B", "0"}, listener)

	// wait till the requested amount of calls received
	for range current {
		if int(counter.Load()) == requestedItems {
			// unregister the listener A, only B will remain
			dispatcher.UnregisterLogListener(PrometheusLogKey{"A", "0"}, listener)
		}
		if int(counter.Load()) == 2*requestedItems {
			// unregister the listener, so we have no more execution of the listener before the end of this test
			dispatcher.UnregisterLogListener(PrometheusLogKey{"B", "0"}, listener)
			mu.Unlock()
			break
		}
		mu.Unlock()
	}
}

func TestLogsDispatchedShutdownTwice(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)

	aUrl := driver.URL("A")

	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)

	node1.EXPECT().GetLabel().AnyTimes().Return("A")

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	dispatcher := NewPrometheusLogDispatcher(net)
	listener := NewMockTimeLogListener(ctrl)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	dispatcher.RegisterLogListener(PrometheusLogKey{"C", "0"}, listener)

	dispatcher.Shutdown()
	dispatcher.Shutdown() /// should not block, panic, etc
}

func TestLogsCannotAddListenerAfterShutdown(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)

	aUrl := driver.URL("A")

	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)

	node1.EXPECT().GetLabel().AnyTimes().Return("A")

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})

	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		return testData[url], nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)

	listener := NewMockTimeLogListener(ctrl)
	// expected to be called
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(1)).AnyTimes()
	// registered after shutdown - cannot be called
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(2)).Times(0)

	dispatcher.RegisterLogListener(PrometheusLogKey{"A_summary", "0.99"}, listener)

	dispatcher.Shutdown()
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_counter", "0"}, listener)

	// wait some time after shutdown to make sure no data has arrived
	time.Sleep(2 * time.Second)
}

func TestParseQuantileEnum(t *testing.T) {
	if Quantile099 != Quantile("0.99") {
		t.Errorf("quantiles do not match")
	}
	if Quantile05 != Quantile("0.5") {
		t.Errorf("quantiles do not match")
	}
}

var (
	testData = map[driver.URL][]PrometheusLogValue{
		"A": {
			{PrometheusLogKey{"A_summary", "0.99"}, summaryPrometheusMetricType, 1},
			{PrometheusLogKey{"A_counter", "0"}, counterPrometheusMetricType, 2},
			{PrometheusLogKey{"A_gauge", "0"}, gaugePrometheusMetricType, 3},

			{PrometheusLogKey{"B", "0.99"}, summaryPrometheusMetricType, 4},
			{PrometheusLogKey{"B", "0"}, counterPrometheusMetricType, 5},
			{PrometheusLogKey{"B", "0"}, gaugePrometheusMetricType, 6},
		},
		"B": {
			{PrometheusLogKey{"A_summary", "0.999"}, summaryPrometheusMetricType, 7},
			{PrometheusLogKey{"A_gauge", "0"}, gaugePrometheusMetricType, 8},
			{PrometheusLogKey{"A_counter", "0"}, counterPrometheusMetricType, 9},
		},
		"C": {
			{PrometheusLogKey{"C", "0"}, counterPrometheusMetricType, 10},
		},
	}
)

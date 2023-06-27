package monitoring

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/golang/mock/gomock"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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
	listener := newTestLogListener()

	dispatcher.AfterNodeCreation(node3)

	// register only for metrics
	// A - summary 0.99 quantile, counter and gauge
	// C
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_summary", 0.99}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_counter", 0}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_gauge", 0}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"C", 0}, listener)

	// wait for data to arrive
	time.Sleep(2 * time.Second)

	// encoded expected values - selected metrics for nodes A, B, C must exist
	expected := []string{"A_1", "A_2", "A_3", "B_8", "B_9", "C_10"}

	for _, want := range expected {
		var found bool
		for _, got := range listener.encoded {
			if want == got {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("value not found: %s, got: %v", want, listener.encoded)
		}
	}
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

	var counter atomic.Int64
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		res := PrometheusLogValue{PrometheusLogKey{"A", 0}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)
	listener := newTestLogListener()

	// listen for metrics
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", 0}, listener)

	// wait for data to arrive
	time.Sleep(2 * time.Second)

	// stop adding more elements
	listener.mut.Lock()
	defer listener.mut.Unlock()

	if len(listener.values) == 0 {
		t.Errorf("no values collected")
	}

	for i, val := range listener.values {
		if i > 0 && val <= listener.values[i-1] {
			t.Errorf("values not ascending: %v > %v", val, listener.values[i])
		}
	}

	if len(listener.times) == 0 {
		t.Errorf("no timestamps collected")
	}

	for i, val := range listener.times {
		if i > 0 && val <= listener.times[i-1] {
			t.Errorf("timestamps not ascending: %v > %v", val, listener.times[i])
		}
	}
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

	var counter atomic.Int64
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		res := PrometheusLogValue{PrometheusLogKey{"A", 0}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)
	listener := newTestLogListener()

	// listen for metrics
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", 0}, listener)

	// wait for some data to arrive
	time.Sleep(2 * time.Second)

	if counter.Load() == 0 {
		t.Errorf("no data loaded before shutdown")
	}

	// counter will not be updated after shutdown
	dispatcher.Shutdown()

	// make sure there is some delay after Shutdown
	sizeAfterShutdown := counter.Load()
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

	metrics := []string{"A", "B"}
	var counter atomic.Int64
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		res := PrometheusLogValue{PrometheusLogKey{metrics[next%2], 0}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)
	listener := newTestLogListener()

	// listen for two metrics A and B
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", 0}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"B", 0}, listener)

	// wait for data to arrive
	time.Sleep(2 * time.Second)

	listener.mut.Lock()
	// altering A, B received -> values mod 2 represents 0 = A or  1 = B,
	// we have to get just continuous increasing values as both metrics are listened to
	for i, val := range listener.values {
		if got, want := int(val), i+1; got != want {
			t.Errorf("unexpected metric: %v != %v", got, want)
		}
	}
	listener.mut.Unlock()

	// unregister from metric A -> will start to receive only B
	dispatcher.UnregisterLogListener(PrometheusLogKey{"A", 0}, listener)

	time.Sleep(1 * time.Second)
	endIndex := int(counter.Load())

	time.Sleep(2 * time.Second)
	listener.mut.Lock()
	defer listener.mut.Unlock()

	// check elements were still arriving after one listener was removed.
	if got, want := len(listener.encoded), endIndex; got <= want {
		t.Errorf("no more elements were added: %d != %d", got, want)
	}

	// must get only B's - check by modulo 2 of value that we are getting only 1 = B
	for i := endIndex; i < len(listener.encoded); i++ {
		if got, want := int(listener.values[i])%2, 1; got != want {
			t.Errorf("unexpected metric encountered: %s != %s", metrics[got], metrics[want])
		}
	}
}

var (
	testData = map[driver.URL][]PrometheusLogValue{
		"A": {
			{PrometheusLogKey{"A_summary", 0.99}, summaryPrometheusMetricType, 1},
			{PrometheusLogKey{"A_counter", 0}, counterPrometheusMetricType, 2},
			{PrometheusLogKey{"A_gauge", 0}, gaugePrometheusMetricType, 3},

			{PrometheusLogKey{"B", 0.99}, summaryPrometheusMetricType, 4},
			{PrometheusLogKey{"B", 0}, counterPrometheusMetricType, 5},
			{PrometheusLogKey{"B", 0}, gaugePrometheusMetricType, 6},
		},
		"B": {
			{PrometheusLogKey{"A_summary", 0.999}, summaryPrometheusMetricType, 7},
			{PrometheusLogKey{"A_gauge", 0}, gaugePrometheusMetricType, 8},
			{PrometheusLogKey{"A_counter", 0}, counterPrometheusMetricType, 9},
		},
		"C": {
			{PrometheusLogKey{"C", 0}, counterPrometheusMetricType, 10},
		},
	}
)

type testLogListener struct {
	encoded []string
	times   []Time
	values  []float64
	mut     sync.Mutex
}

func newTestLogListener() *testLogListener {
	return &testLogListener{
		encoded: make([]string, 0, 1000),
		times:   make([]Time, 0, 1000),
		values:  make([]float64, 0, 1000),
	}
}

func (l *testLogListener) OnLog(node Node, t Time, value float64) {
	l.mut.Lock()
	defer l.mut.Unlock()

	l.encoded = append(l.encoded, fmt.Sprintf("%s_%.f", node, value))
	l.times = append(l.times, t)
	l.values = append(l.values, value)
}

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
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_summary", 0.99}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_counter", 0}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"A_gauge", 0}, listener)
	dispatcher.RegisterLogListener(PrometheusLogKey{"C", 0}, listener)

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

	var counter atomic.Int64
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		res := PrometheusLogValue{PrometheusLogKey{"A", 0}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)

	listener := NewMockTimeLogListener(ctrl)

	var mu sync.Mutex
	max := make(chan bool, 1)
	// check the value is only growing
	var prevVal, prevTime int64
	done := func(node Node, timestamp Time, value float64) {
		mu.Lock()
		defer mu.Unlock()
		if int64(value) < prevVal {
			t.Errorf("values not growing: %d >= %d", int64(value), prevVal)
		}
		if int64(timestamp) < prevTime {
			t.Errorf("times not growing: %d >= %d", int64(timestamp), prevTime)
		}
		prevVal = int64(value)
		prevTime = int64(timestamp)

		// wait for a thousand values
		if int(value) >= 1000 {
			max <- true
		}
	}

	//calls := make([]*gomock.Call, 0, 1000)
	//for i := 1; i <= 1000; i++ {
	//	calls = append(calls, listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), float64(i)).AnyTimes())
	//}

	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Do(done)

	// listen for metrics
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", 0}, listener)

	<-max
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
	max := make(chan bool, 1)
	testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
		next := counter.Add(1)
		if next == 1000 {
			max <- true
		}
		res := PrometheusLogValue{PrometheusLogKey{"A", 0}, counterPrometheusMetricType, float64(next)}
		return []PrometheusLogValue{res}, nil
	}
	dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)
	listener := NewMockTimeLogListener(ctrl)
	listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1000)

	// listen for metrics
	dispatcher.RegisterLogListener(PrometheusLogKey{"A", 0}, listener)

	// wait for first 1000 values
	<-max
	dispatcher.Shutdown()
	// fetch real value after shutdown as there could be execution between reading from the channel and shutdown.
	sizeAfterShutdown := counter.Load()

	// make sure there is some delay after Shutdown
	// the mocked listener cannot be called within this time
	time.Sleep(2 * time.Second)

	if got, want := counter.Load(), sizeAfterShutdown; got > want {
		t.Errorf("more data has accunulated after shutdown: %d != %d", got, want)
	}
}

func TestLogsDispatchedUnregisterListener(t *testing.T) {
	//t.Parallel()
	//
	//ctrl := gomock.NewController(t)
	//net := driver.NewMockNetwork(ctrl)
	//
	//node1 := driver.NewMockNode(ctrl)
	//aUrl := driver.URL("A")
	//node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&aUrl)
	//node1.EXPECT().GetLabel().AnyTimes().Return("A")
	//
	//// simulate existing nodes
	//net.EXPECT().RegisterListener(gomock.Any())
	//net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1})
	//
	//metrics := []string{"A", "B"}
	//var counter atomic.Int64
	//testFunc := func(url driver.URL) ([]PrometheusLogValue, error) {
	//	next := counter.Add(1)
	//	res := PrometheusLogValue{PrometheusLogKey{metrics[next%2], 0}, counterPrometheusMetricType, float64(next)}
	//	return []PrometheusLogValue{res}, nil
	//}
	//dispatcher := newPrometheusLogDispatcher(net, 1*time.Millisecond, testFunc)
	//
	//done := func(node Node, timestamp Time, value float64) {
	//	wg.Done()
	//}
	//
	//listener := NewMockTimeLogListener(ctrl)
	//listener.EXPECT().OnLog(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Do(done)
	//
	//// listen for two metrics A and B
	//dispatcher.RegisterLogListener(PrometheusLogKey{"A", 0}, listener)
	//dispatcher.RegisterLogListener(PrometheusLogKey{"B", 0}, listener)
	//
	//// wait for data to arrive
	//time.Sleep(2 * time.Second)
	//
	//listener.mut.Lock()
	//// altering A, B received -> values mod 2 represents 0 = A or  1 = B,
	//// we have to get just continuous increasing values as both metrics are listened to
	//for i, val := range listener.values {
	//	if got, want := int(val), i+1; got != want {
	//		t.Errorf("unexpected metric: %v != %v", got, want)
	//	}
	//}
	//listener.mut.Unlock()
	//
	//// unregister from metric A -> will start to receive only B
	//dispatcher.UnregisterLogListener(PrometheusLogKey{"A", 0}, listener)
	//
	//time.Sleep(1 * time.Second)
	//endIndex := int(counter.Load())
	//
	//time.Sleep(2 * time.Second)
	//listener.mut.Lock()
	//defer listener.mut.Unlock()
	//
	//// check elements were still arriving after one listener was removed.
	//if got, want := len(listener.encoded), endIndex; got <= want {
	//	t.Errorf("no more elements were added: %d != %d", got, want)
	//}
	//
	//// must get only B's - check by modulo 2 of value that we are getting only 1 = B
	//for i := endIndex; i < len(listener.encoded); i++ {
	//	if got, want := int(listener.values[i])%2, 1; got != want {
	//		t.Errorf("unexpected metric encountered: %s != %s", metrics[got], metrics[want])
	//	}
	//}
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

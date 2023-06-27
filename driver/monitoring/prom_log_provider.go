package monitoring

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/node"
	"log"
	"net/http"
	"sync"
	"time"
)

// TimeLogListener is an interface implemented by any subject that wants to receive
// a value for the specific time and node. Its method is triggered every time a new value
// is available for the timestamp.
type TimeLogListener interface {

	// OnLog is executed every time a new value occurs for the given time and node.
	OnLog(node Node, timestamp Time, value float64)
}

// PrometheusLogKey is a unique identifier of the log entry obtained from prometheus.
// The prometheus log entry contains the name and type.
// The type can be one of the Summary, Counter, or Gauge. If the type is Summary, quantile is provided
// additionally.
// The log entry must have a unique name, i.e. an entry with the same name cannot be defined e.g. for Summary and Counter.
// For this reason, the key is composed of the name and quantile only. For non-Summary types, the quantile is set to zero.
type PrometheusLogKey struct {
	name     string
	quantile float32
}

// NewPrometheusNameKey composes the key using the name of the metric only. The quantile is set to zero.
func NewPrometheusNameKey(name string) *PrometheusLogKey {
	return &PrometheusLogKey{name: name}
}

// PrometheusLogProvider is an interface for registering listeners that will be notified about incoming
// prometheus logs. The listeners subscribe themselves for a concrete log type they want to receive,
// while other logs are ignored.
// The implementation should assure that the logs are sent in the same order they have occurred
// on the target system (i.e. the Opera Node).
type PrometheusLogProvider interface {

	// RegisterLogListener registers the input listener to receive new log entries
	// for the given log type.
	RegisterLogListener(key PrometheusLogKey, listener TimeLogListener)

	// UnregisterLogListener removes the input listener from receiving new logs
	UnregisterLogListener(key PrometheusLogKey, listener TimeLogListener)
}

// PrometheusLogDispatcher allows for registering objects to receive Prometheus log messages
// from nodes. Listeners register themselves for particular log type, which they start to receive.
// This object maintains a list of active nodes in the network and periodically fetches their
// Prometheus logs. The logs are distributed to registered listeners together with the timestamp,
// source node and log value information. The logs are sent ordered according to the timestamp.
// The timestamp is equivalent to the time the log was fetched, not necessarily the time the log
// was produced at the node.
type PrometheusLogDispatcher struct {
	nodes     map[Node]chan Time
	nodesLock sync.Mutex

	listeners     map[PrometheusLogKey]map[TimeLogListener]bool
	listenersLock sync.Mutex

	network driver.Network

	ticker *time.Ticker
	period time.Duration
	wg     sync.WaitGroup
	done   chan bool

	logReader func(driver.URL) ([]PrometheusLogValue, error)
}

// NewPrometheusLogDispatcher creates a new object that periodically parses Prometheus logs
// for all nodes active in the network and dispatches the entries to registered listeners.
// The logs are parsed from all nodes every 1s and distributed to registered listeners.
func NewPrometheusLogDispatcher(network driver.Network) *PrometheusLogDispatcher {

	logReader := func(url driver.URL) ([]PrometheusLogValue, error) {
		resp, err := http.Get(fmt.Sprintf("%s/debug/metrics/prometheus", url))
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		return ParsePrometheusLogReader(resp.Body)
	}
	return newPrometheusLogDispatcher(network, 1*time.Second, logReader)
}

// newPrometheusLogDispatcher is the same as its public counterpart, but allows for setting the period of fetching logs from the nodes,
// and allows for customising the method to fetch logs from nodes.
func newPrometheusLogDispatcher(network driver.Network, period time.Duration, logReader func(driver.URL) ([]PrometheusLogValue, error)) *PrometheusLogDispatcher {
	res := &PrometheusLogDispatcher{
		network:   network,
		nodes:     make(map[Node]chan Time, 50),
		listeners: make(map[PrometheusLogKey]map[TimeLogListener]bool, 50),
		logReader: logReader,
		period:    period,
		ticker:    time.NewTicker(period),
		done:      make(chan bool),
	}

	res.startPeriodicDispatch()

	// listen for new Nodes
	network.RegisterListener(res)

	// get nodes that have been started before this instance creation
	for _, n := range res.network.GetActiveNodes() {
		res.AfterNodeCreation(n)
	}

	return res
}

// Shutdown terminates periodic parsing of the logs. No more new logs will be provided
// after this method is called.
func (n *PrometheusLogDispatcher) Shutdown() {
	n.ticker.Stop()
	n.done <- true
	n.wg.Wait()
}

func (n *PrometheusLogDispatcher) RegisterLogListener(key PrometheusLogKey, listener TimeLogListener) {
	n.listenersLock.Lock()
	defer n.listenersLock.Unlock()

	listeners, exist := n.listeners[key]
	if !exist {
		listeners = make(map[TimeLogListener]bool, 50)
		n.listeners[key] = listeners
	}

	listeners[listener] = true
}

func (n *PrometheusLogDispatcher) UnregisterLogListener(key PrometheusLogKey, listener TimeLogListener) {
	n.listenersLock.Lock()
	defer n.listenersLock.Unlock()
	listeners, exist := n.listeners[key]
	if exist {
		delete(listeners, listener)
	}
}

func (n *PrometheusLogDispatcher) AfterNodeCreation(driverNode driver.Node) {
	n.nodesLock.Lock()
	defer n.nodesLock.Unlock()

	// register the node
	nodeId := Node(driverNode.GetLabel())
	_, exists := n.nodes[nodeId]
	if !exists {
		// each node has its own channel, which inform that the log must be parsed
		// and distributed to listeners.
		// It is done so to assure the logs are provided in the right order,
		// not to swap more planned go routines.
		url := driverNode.GetServiceUrl(&node.OperaDebugService)
		ch := make(chan Time, 100)
		n.nodes[nodeId] = ch
		n.startNodeLogsDispatch(nodeId, url, ch)
	}
}

func (n *PrometheusLogDispatcher) AfterApplicationCreation(driver.Application) {
	// ignored
}

// startPeriodicDispatch starts a go-routine that periodically triggers the
// fetching, parsing, and distribution of node logs to registered listeners.
// This method only sends the signal to trigger node parsing, and waits
// to send next signal every period.
// Each node maintains its own channel, which triggers parsing of its log
// to assure the logs are always parsed in the same order.
func (n *PrometheusLogDispatcher) startPeriodicDispatch() {
	go func() {
		for {
			select {
			case <-n.done:
				n.nodesLock.Lock()
				for nodeId, ch := range n.nodes {
					close(ch)
					delete(n.nodes, nodeId)
				}
				n.nodesLock.Unlock()
				return
			case t := <-n.ticker.C:
				n.nodesLock.Lock()
				for _, ch := range n.nodes {
					ch <- NewTime(t)
				}
				n.nodesLock.Unlock()
			}
		}
	}()
}

// startNodeLogsDispatch starts a go-routine, which parses logs of the input node and distributes the logs
// to registered listeners.
// The input channel is read and every time it contains data, parsing of log is triggered.
// When the log is parsed, the go-routine is blocked on the channel until next signal arrives.
func (n *PrometheusLogDispatcher) startNodeLogsDispatch(nodeId Node, url *driver.URL, ch chan Time) {
	n.wg.Add(1)
	go func(node Node, url driver.URL) {
		defer n.wg.Done()
		for range ch {
			if logs, err := n.logReader(url); err == nil {
				n.distributeLog(NewTime(time.Now()), node, logs)
			} else {
				log.Printf("failed to parse log: %s", err)
			}
		}
	}(nodeId, *url)
}

// distributeLog sends the input log into all listeners, using the input timestamp and the nodeID.
// This method locks all listeners until the logs are distributed, i.e. it is coarse grained
// at the moment, as assumption is that the slowest part of logs processing is actually I/O to retrieve
// the logs, but not their distribution to receivers.
func (n *PrometheusLogDispatcher) distributeLog(timestamp Time, nodeId Node, logs []PrometheusLogValue) {
	n.listenersLock.Lock()
	localCopy := make(map[PrometheusLogKey]map[TimeLogListener]bool, len(n.listeners))
	for _, value := range logs {
		receivers := n.listeners[value.PrometheusLogKey]
		for receiver := range receivers {
			localCopy[value.PrometheusLogKey][receiver] = true
		}
	}
	defer n.listenersLock.Unlock()

	for _, value := range logs {
		receivers := localCopy[value.PrometheusLogKey]
		for receiver := range receivers {
			receiver.OnLog(nodeId, timestamp, value.value)
		}
	}
}

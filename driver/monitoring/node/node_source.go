package nodemon

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
)

// Sensor is an abstraction of some input device capable of probing a node
// for some metric of type T.
type Sensor[T any] interface {
	ReadValue() (T, error)
}

// SensorFactory is a factory for sensors targeting selected nodes.
type SensorFactory[T any] interface {
	CreateSensor(driver.Node) (Sensor[T], error)
}

// periodicNodeDataSource is a generic data source periodically querying
// node-associated sensors for data.
type periodicNodeDataSource[T any] struct {
	metric   mon.Metric[mon.Node, mon.TimeSeries[T]]
	network  driver.Network
	period   time.Duration
	factory  SensorFactory[T]
	data     map[mon.Node]*mon.SyncedSeries[mon.Time, T]
	dataLock sync.Mutex
	stop     chan bool  // used to signal per-node collectors about the shutdown
	done     chan error // used to signal collector shutdown to source
}

// NewPeriodicNodeDataSource creates a new data source managing per-node sensor
// instances for a given metric and periodically collecting data from those.
func NewPeriodicNodeDataSource[T any](
	metric mon.Metric[mon.Node, mon.TimeSeries[T]],
	network driver.Network,
	factory SensorFactory[T],
) mon.Source[mon.Node, mon.TimeSeries[T]] {
	return newPeriodicNodeDataSource(metric, network, time.Second, factory)
}

// newPeriodicNodeDataSource is the same as NewPeriodicNodeDataSource but with
// a customizable sampling periode.
func newPeriodicNodeDataSource[T any](
	metric mon.Metric[mon.Node, mon.TimeSeries[T]],
	network driver.Network,
	period time.Duration,
	factory SensorFactory[T],
) mon.Source[mon.Node, mon.TimeSeries[T]] {
	stop := make(chan bool)
	done := make(chan error)

	res := &periodicNodeDataSource[T]{
		metric:  metric,
		network: network,
		period:  period,
		factory: factory,
		data:    map[mon.Node]*mon.SyncedSeries[mon.Time, T]{},
		stop:    stop,
		done:    done,
	}

	network.RegisterListener(res)

	for _, node := range network.GetActiveNodes() {
		res.AfterNodeCreation(node)
	}

	return res
}

func startCollector[T any](
	node driver.Node,
	period time.Duration,
	factory SensorFactory[T],
	stop <-chan bool,
	done chan<- error,
) *mon.SyncedSeries[mon.Time, T] {
	res := &mon.SyncedSeries[mon.Time, T]{}

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		sensor, err := factory.CreateSensor(node)
		if err != nil {
			return
		}

		var errs []error
		ticker := time.NewTicker(period)
		for {
			select {
			case now := <-ticker.C:
				value, err := sensor.ReadValue()
				if err != nil {
					errs = append(errs, err)
				} else {
					res.Append(mon.NewTime(now), value)
				}
			case <-stop:
				err = errors.Join(errs...)
				return
			}
		}
	}()

	return res
}

func (s *periodicNodeDataSource[T]) GetMetric() mon.Metric[mon.Node, mon.TimeSeries[T]] {
	return s.metric
}

func (s *periodicNodeDataSource[T]) GetSubjects() []mon.Node {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	res := make([]mon.Node, 0, len(s.data))
	for node := range s.data {
		res = append(res, node)
	}
	return res
}

func (s *periodicNodeDataSource[T]) GetData(node mon.Node) (mon.TimeSeries[T], bool) {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	res, exists := s.data[node]
	return res, exists
}

func (s *periodicNodeDataSource[T]) Shutdown() error {
	if s.network == nil {
		return nil
	}
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	s.network.UnregisterListener(s)
	s.network = nil
	close(s.stop)

	<-s.done
	return nil
}

func (s *periodicNodeDataSource[T]) AfterNodeCreation(node driver.Node) {
	nodeId, err := node.GetNodeID()
	if err != nil {
		log.Printf("failed to obtain node ID of node, will not be able to track block height: %v", err)
		return
	}
	s.dataLock.Lock()
	defer s.dataLock.Unlock()

	_, present := s.data[mon.Node(nodeId)]
	if present {
		// TODO: improve logging by tracking source (see Aida as a reference)
		log.Printf("received notification of already known node")
		return
	}

	s.data[mon.Node(nodeId)] = startCollector(node, s.period, s.factory, s.stop, s.done)
}

func (s *periodicNodeDataSource[T]) AfterApplicationCreation(driver.Application) {
	// ignored
}

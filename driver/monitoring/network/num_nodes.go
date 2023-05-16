package netmon

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
)

// NumberOfNodes retains a time-series for the number of nodes in the network
// run by Norma. This includes all types of nodes.
var NumberOfNodes = mon.Metric[mon.Network, mon.TimeSeries[int]]{
	Name:        "NumberOfNodes",
	Description: "The number of connected nodes at various times.",
}

func init() {
	if err := mon.RegisterSource(NumberOfNodes, NewNumNodesSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// numNodesSource is a monitoring data source tracking the number of active
// nodes in a network environment.
type numNodesSource struct {
	network driver.Network
	data    mon.SyncedSeries[mon.Time, int]
	stop    chan<- bool
	done    <-chan bool
}

// NewNumNodesSource creates a new data source periodically collecting data on
// the number of nodes in the network.
func NewNumNodesSource(network driver.Network) mon.Source[mon.Network, mon.TimeSeries[int]] {
	return newNumNodesSource(network, time.Second)
}

// newNumNodesSource creates a new data source periodically collecting data on
// the number of nodes in the network.
func newNumNodesSource(network driver.Network, period time.Duration) mon.Source[mon.Network, mon.TimeSeries[int]] {
	stop := make(chan bool)
	done := make(chan bool)

	res := &numNodesSource{
		network: network,
		stop:    stop,
		done:    done,
	}

	go func() {
		defer close(done)
		ticker := time.NewTicker(period)
		for {
			select {
			case now := <-ticker.C:
				numNodes := len(network.GetActiveNodes())
				res.data.Append(mon.NewTime(now), numNodes)
			case <-stop:
				return
			}
		}
	}()

	return res
}

func (s *numNodesSource) GetMetric() mon.Metric[mon.Network, mon.TimeSeries[int]] {
	return NumberOfNodes
}

func (s *numNodesSource) GetSubjects() []mon.Network {
	return []mon.Network{{}}
}

func (s *numNodesSource) GetData(mon.Network) *mon.TimeSeries[int] {
	var res mon.TimeSeries[int] = &s.data
	return &res
}

func (s *numNodesSource) Shutdown() error {
	close(s.stop)
	<-s.done
	return nil
}

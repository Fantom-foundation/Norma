package nodemon

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/ethereum/go-ethereum/rpc"
)

// NodeBlockHeight collects a per-node time series of its current block height.
var NodeBlockHeight = mon.Metric[mon.Node, mon.TimeSeries[int]]{
	Name:        "NodeBlockHeight",
	Description: "The block height of nodes at various times.",
}

// nodeBlockHeightSource is a data source for tracking the block height of individuel
// nodes over time.
type nodeBlockHeightSource struct {
	network  driver.Network
	period   time.Duration
	data     map[mon.Node]*mon.SyncedSeries[mon.Time, int]
	dataLock sync.Mutex
	stop     chan bool  // used to signal per-node collectors about the shutdown
	done     chan error // used to signal collector shutdown to source
}

// NewNumNodesSource creates a new data source periodically collecting data on
// the number of nodes in the network.
func NewNodeBlockHeightSource(network driver.Network) mon.Source[mon.Node, mon.TimeSeries[int]] {
	return newNodeBlockHeightSource(network, time.Second)
}

func newNodeBlockHeightSource(network driver.Network, period time.Duration) mon.Source[mon.Node, mon.TimeSeries[int]] {
	stop := make(chan bool)
	done := make(chan error)

	res := &nodeBlockHeightSource{
		network: network,
		period:  period,
		data:    map[mon.Node]*mon.SyncedSeries[mon.Time, int]{},
		stop:    stop,
		done:    done,
	}

	network.RegisterListener(res)

	for _, node := range network.GetActiveNodes() {
		res.AfterNodeCreation(node)
	}

	return res
}

func startCollector(node driver.Node, period time.Duration, stop <-chan bool, done chan<- error) *mon.SyncedSeries[mon.Time, int] {
	res := &mon.SyncedSeries[mon.Time, int]{}

	go func() {
		var err error
		defer func() {
			done <- err
		}()

		url := node.GetRpcServiceUrl()
		if url == nil {
			err = fmt.Errorf("node does not export an RPC server")
			return
		}
		rpcClient, err := rpc.DialContext(context.Background(), string(*url))
		if err != nil {
			return
		}

		var errs []error
		ticker := time.NewTicker(period)
		for {
			select {
			case now := <-ticker.C:
				var blockNumber string
				err = rpcClient.Call(&blockNumber, "eth_blockNumber")
				if err != nil {
					errs = append(errs, err)
					continue
				}
				blockNumber = strings.TrimPrefix(blockNumber, "0x")
				value, err := strconv.ParseInt(blockNumber, 16, 32)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				res.Append(mon.NewTime(now), int(value))
			case <-stop:
				err = errors.Join(errs...)
				return
			}
		}
	}()

	return res
}

func (s *nodeBlockHeightSource) GetMetric() mon.Metric[mon.Node, mon.TimeSeries[int]] {
	return NodeBlockHeight
}

func (s *nodeBlockHeightSource) GetSubjects() []mon.Node {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	res := make([]mon.Node, 0, len(s.data))
	for node := range s.data {
		res = append(res, node)
	}
	return res
}

func (s *nodeBlockHeightSource) GetData(node mon.Node) *mon.TimeSeries[int] {
	s.dataLock.Lock()
	defer s.dataLock.Unlock()
	var res mon.TimeSeries[int] = s.data[node]
	return &res
}

func (s *nodeBlockHeightSource) Shutdown() error {
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

func (s *nodeBlockHeightSource) AfterNodeCreation(node driver.Node) {
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

	s.data[mon.Node(nodeId)] = startCollector(node, s.period, s.stop, s.done)
}

func (s *nodeBlockHeightSource) AfterApplicationCreation(driver.Application) {
	// ignored
}

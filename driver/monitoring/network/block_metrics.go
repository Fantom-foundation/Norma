package netmon

import (
	"fmt"
	"log"
	"sync"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/export"
)

var (
	// BlockNumberOfTransactions is a metric capturing number of transactions for each block of a node.
	BlockNumberOfTransactions = monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]]{
		Name:        "BlockNumberOfTransactions",
		Description: "The number of transactions processed in a block",
	}

	// BlockGasUsed is a metric capturing Gas used for each block of a node.
	BlockGasUsed = monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]]{
		Name:        "BlockGasUsed",
		Description: "The gas used in a block",
	}
)

func init() {
	if err := monitoring.RegisterSource(BlockNumberOfTransactions, newNumberOfTransactionsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}

	if err := monitoring.RegisterSource(BlockGasUsed, newGasUsedSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// BlockNetworkMetricSource is a metric source that captures block properties where the Metric is the subject
type BlockNetworkMetricSource[T any] struct {
	metric           monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, T]]
	getBlockProperty func(b monitoring.Block) T
	monitor          *monitoring.Monitor
	series           *monitoring.SyncedSeries[monitoring.BlockNumber, T]
	seriesLock       *sync.Mutex
	lastBlock        int // track last block added in the series not to add duplicated block heights
}

// NewNumberOfTransactionsSource creates a metric capturing number of transactions for each block of a network
func NewNumberOfTransactionsSource(monitor *monitoring.Monitor) *BlockNetworkMetricSource[int] {
	f := func(b monitoring.Block) int {
		return b.Txs
	}
	return newBlockNetworkMetricsSource[int](monitor, f, BlockNumberOfTransactions)
}

// NewGasUsedSource creates a metric capturing Gas used for each block of a network.
func NewGasUsedSource(monitor *monitoring.Monitor) *BlockNetworkMetricSource[int] {
	f := func(b monitoring.Block) int {
		return b.GasUsed
	}
	return newBlockNetworkMetricsSource[int](monitor, f, BlockGasUsed)
}

// newNumberOfTransactionsSource is the same as its public counterpart, it only returns the Source interface instead of the struct to be used in factories
func newNumberOfTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]] {
	return NewNumberOfTransactionsSource(monitor)
}

// newGasUsedSource is the same as its public counterpart, it only returns the Source interface instead of the struct to be used in factories
func newGasUsedSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]] {
	return NewGasUsedSource(monitor)
}

// newBlockNodeMetricsSource creates a new data source periodically collecting data from the Node log
func newBlockNetworkMetricsSource[T any](
	monitor *monitoring.Monitor,
	getBlockProperty func(b monitoring.Block) T,
	metric monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, T]]) *BlockNetworkMetricSource[T] {

	m := &BlockNetworkMetricSource[T]{
		metric:           metric,
		getBlockProperty: getBlockProperty,
		monitor:          monitor,
		series:           &monitoring.SyncedSeries[monitoring.BlockNumber, T]{},
		seriesLock:       &sync.Mutex{},
		lastBlock:        -1,
	}

	monitor.NodeLogProvider().RegisterLogListener(m)
	monitor.Writer().Add(func() error {
		source := (monitoring.Source[monitoring.Network, monitoring.Series[monitoring.BlockNumber, T]])(m)
		return export.AddSeriesData(monitor.Writer(), source)
	})

	return m
}

func (s *BlockNetworkMetricSource[T]) GetMetric() monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, T]] {
	return s.metric
}

func (s *BlockNetworkMetricSource[T]) GetSubjects() []monitoring.Network {
	var item monitoring.Network
	return []monitoring.Network{item}
}

func (s *BlockNetworkMetricSource[T]) GetData(monitoring.Network) (monitoring.Series[monitoring.BlockNumber, T], bool) {
	return s.series, true
}

func (s *BlockNetworkMetricSource[T]) Shutdown() error {
	s.monitor.NodeLogProvider().UnregisterLogListener(s)
	return nil
}

func (s *BlockNetworkMetricSource[T]) OnBlock(_ monitoring.Node, block monitoring.Block) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	if block.Height > s.lastBlock {
		if err := s.series.Append(monitoring.BlockNumber(block.Height), s.getBlockProperty(block)); err != nil {
			log.Printf("error to add to the series: %s", err)
		}
		s.lastBlock = block.Height
	}
}

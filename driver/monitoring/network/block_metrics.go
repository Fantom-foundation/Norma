package netmon

import (
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/ethereum/go-ethereum/log"
	"sync"
)

var (
	// BlockNumberOfTransactions is a metric capturing number of transactions for each block of a node.
	BlockNumberOfTransactions = monitoring.Metric[monitoring.Network, monitoring.BlockSeries[int]]{
		Name:        "BlockNumberOfTransactions",
		Description: "The number of transactions processed in a block",
	}

	// BlockGasUsed is a metric capturing Gas used for each block of a node.
	BlockGasUsed = monitoring.Metric[monitoring.Network, monitoring.BlockSeries[int]]{
		Name:        "BlockGasUsed",
		Description: "The gas used in a block",
	}
)

// BlockNetworkMetricSource is a metric source that captures block properties where the Metric is the subject
type BlockNetworkMetricSource[T any] struct {
	metric           monitoring.Metric[monitoring.Network, monitoring.BlockSeries[T]]
	getBlockProperty func(b monitoring.Block) T
	registry         monitoring.NodeLogProvider
	series           *monitoring.SyncedSeries[monitoring.BlockNumber, T]
	seriesLock       sync.Mutex
	lastBlock        int // track last block added in the series not to add duplicated block heights
}

// NewNumberOfTransactionsSource creates a metric capturing number of transactions for each block of a node.
func NewNumberOfTransactionsSource(reg monitoring.NodeLogProvider) *BlockNetworkMetricSource[int] {
	f := func(b monitoring.Block) int {
		return b.Txs
	}
	return newBlockNetworkMetricsSource[int](reg, f, BlockNumberOfTransactions)
}

// NewGasUsedSource creates a metric capturing Gas used for each block of a node.
func NewGasUsedSource(reg monitoring.NodeLogProvider) *BlockNetworkMetricSource[int] {
	f := func(b monitoring.Block) int {
		return b.GasUsed
	}
	return newBlockNetworkMetricsSource[int](reg, f, BlockGasUsed)
}

// newBlockNodeMetricsSource creates a new data source periodically collecting data from the Node log
func newBlockNetworkMetricsSource[T any](
	reg monitoring.NodeLogProvider,
	getBlockProperty func(b monitoring.Block) T,
	metric monitoring.Metric[monitoring.Network, monitoring.BlockSeries[T]]) *BlockNetworkMetricSource[T] {

	m := &BlockNetworkMetricSource[T]{
		metric:           metric,
		getBlockProperty: getBlockProperty,
		registry:         reg,
		series:           &monitoring.SyncedSeries[monitoring.BlockNumber, T]{},
		lastBlock:        -1,
	}

	reg.RegisterLogListener(m)

	return m
}

func (s *BlockNetworkMetricSource[T]) GetMetric() monitoring.Metric[monitoring.Network, monitoring.BlockSeries[T]] {
	return s.metric
}

func (s *BlockNetworkMetricSource[T]) GetSubjects() []monitoring.Network {
	var item monitoring.Network
	return []monitoring.Network{item}
}

func (s *BlockNetworkMetricSource[T]) GetData(monitoring.Network) monitoring.BlockSeries[T] {
	return s.series
}

func (s *BlockNetworkMetricSource[T]) Shutdown() error {
	s.registry.UnregisterLogListener(s)
	return nil
}

func (s *BlockNetworkMetricSource[T]) OnBlock(_ monitoring.Node, block monitoring.Block) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	if block.Height > s.lastBlock {
		if err := s.series.Append(monitoring.BlockNumber(block.Height), s.getBlockProperty(block)); err != nil {
			log.Error("error to add to the series: %s", err)
		}
		s.lastBlock = block.Height
	}
}
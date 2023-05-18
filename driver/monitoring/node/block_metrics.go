package nodemon

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/ethereum/go-ethereum/log"
	"sync"
	"time"
)

var (
	// BlockCompletionTime is a metric capturing time of the block finalisation.
	BlockCompletionTime = monitoring.Metric[monitoring.Node, monitoring.BlockSeries[time.Time]]{
		Name:        "BlockCompletionTime",
		Description: "Time the block was completed",
	}

	BlockEventAndTxsProcessingTime = monitoring.Metric[monitoring.Node, monitoring.BlockSeries[time.Duration]]{
		Name:        "BlockEventAndTxsProcessingTime",
		Description: "Time to process a block, it applies all lachesis events, applies all transactions, and commits stateDB",
	}
)

func init() {
	if err := monitoring.RegisterSource(BlockCompletionTime, monitoring.AdaptLogProviderToMonitorFactory(newBlockTimeSource)); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
	if err := monitoring.RegisterSource(BlockEventAndTxsProcessingTime, monitoring.AdaptLogProviderToMonitorFactory(newBlockProcessingTimeSource)); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// BlockNodeMetricSource is a metric source that captures block properties where the Node is the subject
type BlockNodeMetricSource[T any] struct {
	metric           monitoring.Metric[monitoring.Node, monitoring.BlockSeries[T]]
	getBlockProperty func(b monitoring.Block) T
	registry         monitoring.NodeLogProvider
	series           map[monitoring.Node]*monitoring.SyncedSeries[monitoring.BlockNumber, T]
	seriesLock       sync.Mutex
}

// NewBlockTimeSource creates a metric capturing time of the block finalisation for each Node.
func NewBlockTimeSource(reg monitoring.NodeLogProvider) *BlockNodeMetricSource[time.Time] {
	f := func(b monitoring.Block) time.Time {
		return b.Time
	}
	return newBlockNodeMetricsSource[time.Time](reg, f, BlockCompletionTime)
}

// newBlockTimeSource is the same as its public counterpart, it only returns the struct instead of the Source interface
func newBlockTimeSource(reg monitoring.NodeLogProvider) monitoring.Source[monitoring.Node, monitoring.BlockSeries[time.Time]] {
	return NewBlockTimeSource(reg)
}

// NewBlockProcessingTimeSource creates a metric capturing time of the block finalisation for each Node.
func NewBlockProcessingTimeSource(reg monitoring.NodeLogProvider) *BlockNodeMetricSource[time.Duration] {
	f := func(b monitoring.Block) time.Duration {
		return b.ProcessingTime
	}
	return newBlockNodeMetricsSource[time.Duration](reg, f, BlockEventAndTxsProcessingTime)
}

// newBlockProcessingTimeSource is the same as its public counterpart, it only returns the struct instead of the Source interface
func newBlockProcessingTimeSource(reg monitoring.NodeLogProvider) monitoring.Source[monitoring.Node, monitoring.BlockSeries[time.Duration]] {
	return NewBlockProcessingTimeSource(reg)
}

// newBlockNodeMetricsSource creates a new data source periodically collecting data from the Node log
func newBlockNodeMetricsSource[T any](
	reg monitoring.NodeLogProvider,
	getBlockProperty func(b monitoring.Block) T,
	metric monitoring.Metric[monitoring.Node, monitoring.BlockSeries[T]]) *BlockNodeMetricSource[T] {

	m := &BlockNodeMetricSource[T]{
		metric:           metric,
		getBlockProperty: getBlockProperty,
		registry:         reg,
		series:           make(map[monitoring.Node]*monitoring.SyncedSeries[monitoring.BlockNumber, T], 50),
	}

	reg.RegisterLogListener(m)

	return m
}

func (s *BlockNodeMetricSource[T]) GetMetric() monitoring.Metric[monitoring.Node, monitoring.BlockSeries[T]] {
	return s.metric
}

func (s *BlockNodeMetricSource[T]) GetSubjects() []monitoring.Node {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	res := make([]monitoring.Node, 0, len(s.series))
	for k := range s.series {
		res = append(res, k)
	}
	return res
}

func (s *BlockNodeMetricSource[T]) GetData(node monitoring.Node) (monitoring.BlockSeries[T], bool) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	res, exists := s.series[node]
	return res, exists
}

func (s *BlockNodeMetricSource[T]) Shutdown() error {
	s.registry.UnregisterLogListener(s)
	return nil
}

func (s *BlockNodeMetricSource[T]) OnBlock(node monitoring.Node, block monitoring.Block) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	series, exists := s.series[node]
	if !exists {
		series = &monitoring.SyncedSeries[monitoring.BlockNumber, T]{}
		s.series[node] = series
	}

	if err := series.Append(monitoring.BlockNumber(block.Height), s.getBlockProperty(block)); err != nil {
		log.Error("error to add to the series: %s", err)
	}
}

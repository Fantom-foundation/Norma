package nodemon

import (
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
)

// BlockNodeMetricSource is a metric source that captures block properties where the Node is the subject
type BlockNodeMetricSource[T any] struct {
	metric           monitoring.Metric[monitoring.Node, monitoring.BlockSeries[T]]
	getBlockProperty func(b monitoring.Block) T
	registry         monitoring.NodeLogProvider
	series           map[monitoring.Node]*monitoring.SyncedSeries[monitoring.BlockNumber, T]
	seriesLock       sync.Mutex
}

// NewBlockTimeSource creates a metric capturing time of the block creation.
func NewBlockTimeSource(reg monitoring.NodeLogProvider) *BlockNodeMetricSource[time.Time] {
	f := func(b monitoring.Block) time.Time {
		return b.Time
	}
	return newBlockNodeMetricsSource[time.Time](reg, f, BlockCompletionTime)
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

func (s *BlockNodeMetricSource[T]) GetData(node monitoring.Node) monitoring.BlockSeries[T] {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	var res monitoring.BlockSeries[T]
	if val, exists := s.series[node]; exists {
		res = val
	}
	return res
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

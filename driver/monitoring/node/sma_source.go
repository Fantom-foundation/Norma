package nodemon

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/export"
	"golang.org/x/exp/constraints"
	"sync"
)

func init() {
	smaPeriods := []int{10, 100, 1000}
	for _, period := range smaPeriods {

		// TransactionThroughputSMA is a metric that aggregates output from another series and computes
		// Simple Moving Average.
		TransactionThroughputSMA := monitoring.Metric[monitoring.Node, monitoring.BlockSeries[float32]]{
			Name:        fmt.Sprintf("TransactionsThroughputSMA_%d", period),
			Description: "Transaction throughput standard moving average",
		}

		smaFactory := func(input monitoring.BlockSeries[float32]) monitoring.BlockSeries[float32] {
			return monitoring.NewSMASeries[monitoring.BlockNumber, float32](input, period)
		}
		metricsFactory := func(monitor *monitoring.Monitor) monitoring.Source[monitoring.Node, monitoring.BlockSeries[float32]] {
			return newNodeBlockSeriesTransformation(monitor, TransactionThroughputSMA, TransactionsThroughput, smaFactory)
		}

		if err := monitoring.RegisterSource(TransactionThroughputSMA, metricsFactory); err != nil {
			panic(fmt.Sprintf("failed to register metric source: %v", err))
		}
	}
}

// NodeBlockSeriesTransformation is a source that captures an input source, and computes certain transformation on top of it.
// The input source for this type must have the node as a subject and the Series as a value.
// This type produces the same Nodes as the subjects and the series with the required transformation.
type NodeBlockSeriesTransformation[K constraints.Ordered, T any, X monitoring.Series[K, T]] struct {
	metric        monitoring.Metric[monitoring.Node, X]
	source        monitoring.Metric[monitoring.Node, X] // source metrics to transform
	monitor       *monitoring.Monitor
	series        map[monitoring.Node]X
	seriesFactory func(X) X // transform input series to the output series
	seriesLock    *sync.Mutex
}

// NewNodeSeriesTransformation creates a new source that can transform input source to the output source.
// This transformation is limited to the source where the Node is the subject and values are series.
// The output source is transformed to contain the same subjects, which addresses new, transformed, series
func NewNodeSeriesTransformation[K constraints.Ordered, T any, X monitoring.Series[K, T]](
	monitor *monitoring.Monitor,
	metric monitoring.Metric[monitoring.Node, X],
	source monitoring.Metric[monitoring.Node, X],
	seriesFactory func(X) X) *NodeBlockSeriesTransformation[K, T, X] {

	m := &NodeBlockSeriesTransformation[K, T, X]{
		metric:        metric,
		source:        source,
		seriesFactory: seriesFactory,
		monitor:       monitor,
		series:        make(map[monitoring.Node]X, 50),
		seriesLock:    &sync.Mutex{},
	}

	return m
}

// newNodeBlockSeriesTransformation creates the same instance as public NewNodeSeriesTransformation but typed to the BlockSeries as a Series.
func newNodeBlockSeriesTransformation[T any](
	monitor *monitoring.Monitor,
	metric monitoring.Metric[monitoring.Node, monitoring.BlockSeries[T]],
	source monitoring.Metric[monitoring.Node, monitoring.BlockSeries[T]],
	seriesFactory func(monitoring.BlockSeries[T]) monitoring.BlockSeries[T]) monitoring.Source[monitoring.Node, monitoring.BlockSeries[T]] {

	res := NewNodeSeriesTransformation[monitoring.BlockNumber, T, monitoring.BlockSeries[T]](monitor, metric, source, seriesFactory)
	monitor.Writer().Add(func() error {
		return export.AddNodeBlockSeriesSource[T](monitor.Writer(), res, export.DirectConverter[T]{})
	})

	return res
}

func (s *NodeBlockSeriesTransformation[K, T, X]) GetMetric() monitoring.Metric[monitoring.Node, X] {
	return s.metric
}

func (s *NodeBlockSeriesTransformation[K, T, X]) GetSubjects() []monitoring.Node {
	return monitoring.GetSubjects(s.monitor, s.source)
}

func (s *NodeBlockSeriesTransformation[K, T, X]) GetData(node monitoring.Node) (X, bool) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	res, exists := s.series[node]
	if !exists {
		source, existsSource := monitoring.GetData(s.monitor, node, s.source)
		if existsSource {
			newSeries := s.seriesFactory(source)
			s.series[node] = newSeries
			return newSeries, true
		}
	}

	return res, exists
}

func (s *NodeBlockSeriesTransformation[K, T, X]) Shutdown() error {
	return nil
}

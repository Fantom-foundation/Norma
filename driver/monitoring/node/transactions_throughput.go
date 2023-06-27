package nodemon

import (
	"fmt"
	"log"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
)

var (
	// TransactionsThroughput is a metric capturing number of transactions per certain time period, i.e. the throughput
	TransactionsThroughput = monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]]{
		Name:        "TransactionsThroughput",
		Description: "The number of transactions processed per certain time period by each node",
	}
)

func init() {
	if err := monitoring.RegisterSource(TransactionsThroughput, newTransactionsThroughputSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// TransactionsThroughputSource is a metric source that captures transaction throughput.
type TransactionsThroughputSource struct {
	BlockNodeMetricSource[float32]
	lastTimes map[monitoring.Node]time.Time // timestamps of the latest received blocks
}

// NewTransactionsThroughputSource creates a metric capturing transaction throughput.
func NewTransactionsThroughputSource(monitor *monitoring.Monitor) *TransactionsThroughputSource {
	blockMetrics := BlockNodeMetricSource[float32]{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(TransactionsThroughput),
		monitor:            monitor,
	}

	m := &TransactionsThroughputSource{
		BlockNodeMetricSource: blockMetrics,
		lastTimes:             make(map[monitoring.Node]time.Time, 50),
	}
	monitor.NodeLogProvider().RegisterLogListener(m)

	return m
}

// newTransactionsThroughputSource is the same as its public counterpart, it only returns the Source interface instead of the struct to be used in factories
func newTransactionsThroughputSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Node, monitoring.Series[monitoring.BlockNumber, float32]] {
	return NewTransactionsThroughputSource(monitor)
}

func (s *TransactionsThroughputSource) OnBlock(node monitoring.Node, block monitoring.Block) {

	prevTime, exists := s.lastTimes[node]
	s.lastTimes[node] = block.Time
	if !exists {
		// very first node received - no difference can be computed, but the data series is expected to be created
		s.GetOrAddSubject(node)
		return
	}

	timeDiff := block.Time.Sub(prevTime).Nanoseconds()
	// prevent NaN or Inf: when the time difference is bellow measured value, skip the block.
	if timeDiff != 0 {
		txs := float64(block.Txs) * 1e9 / float64(timeDiff)
		series := s.GetOrAddSubject(node)
		if err := series.Append(monitoring.BlockNumber(block.Height), float32(txs)); err != nil {
			log.Printf("error to add to the series: %s", err)
		}
	}
}

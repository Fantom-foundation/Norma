package monitoring

import (
	"time"
)

// Metric defines a metric in the monitoring system. The type S is the type of
// subject this metric is collected for (e.g., a node or the full network), and
// the type T is the type of value produced by this metric.
//
// The subject of a metric is the object the metric's property is to be
// associated to. There are two major subjects:
//
//   - monitoring.Network ... to be used for network-wide properties like the
//     number of transactions in a block or the utilized gas. Those metrics are
//     consistent throughout the network and do not require any finer
//     granularity.
//
//   - monitoring.Node ... to be used for node-level properties like the time a
//     block was completed, or the CPU usage at a givne time.
//
// Metric data is typically organized in data series, of which there are two
// main types: monitoring.TimeSeries and monitoring.BlockSeries. The former is
// associating a value to various points in time (using absolute time-stamps).
// The latter associates a value to various block-numbers.
type Metric[S any, T any] struct {
	Name        string // used for unique identification of a metric
	Description string // a description documenting the details of the metric
}

// Per-Node metrics:
var (
	NodeCpuUsage = Metric[Node, TimeSeries[Percent]]{
		Name:        "NodeCpuUsage",
		Description: "The relative CPU usage of a node at various times.",
	}

	BlockReadyForProcessingTime = Metric[Node, BlockSeries[time.Time]]{
		Name:        "BlockReadyForProcessingTime",
		Description: "The time at which a block was ready to be processed on a node.",
	}

	BlockTransactionProcessingFinishTime = Metric[Node, BlockSeries[time.Time]]{
		Name:        "BlockTransactionProcessingFinishTime",
		Description: "The time the processing of the transactions of a block has finished on a node.",
	}
)

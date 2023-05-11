package monitoring

import (
	"time"
)

// Metric defines a metric in the monitoring system. The type S is the type of
// subject this metric is collected for (e.g., a node or the full network), and
// the type T is the type of value produced by this metric.
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

package nodemon

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
	"log"
)

var (
	// A list of Prometheus metrics that will be registered and obtained for processing.
	metrics = []monitoring.PrometheusLogKey{
		monitoring.NewPrometheusNameKey("txpool_received"),

		monitoring.NewPrometheusNameKey("txpool_valid"),
		monitoring.NewPrometheusNameKey("txpool_invalid"),
		monitoring.NewPrometheusNameKey("txpool_underpriced"),
		monitoring.NewPrometheusNameKey("txpool_overflowed"),

		monitoring.NewPrometheusNameKey("txpool_pending"),
		monitoring.NewPrometheusNameKey("txpool_queued"),

		monitoring.NewPrometheusNameKey("system_cpu_procload"),
		monitoring.NewPrometheusNameKey("system_memory_used"),
		monitoring.NewPrometheusNameKey("db_size"),
		monitoring.NewPrometheusNameKey("statedb_disksize"),
	}
)

func init() {
	for _, metric := range metrics {
		metric := metric
		metricsFactory := func(monitor *monitoring.Monitor) monitoring.Source[monitoring.Node, monitoring.Series[monitoring.Time, float64]] {
			return NewPromLogSource(monitor, metric)
		}
		if err := monitoring.RegisterSource(toMetric(metric), metricsFactory); err != nil {
			panic(fmt.Sprintf("failed to register metric source: %v", err))
		}
	}
}

// PromLogSource is a generic metric source for all metrics obtained via Prometheus API
// from the Nodes. It is configured with the Prometheus metric of interest,
// and it listens for incoming metric data of all running Nodes.
type PromLogSource struct {
	*utils.SyncedSeriesSource[monitoring.Node, monitoring.Time, float64]
}

// NewPromLogSource creates a new instance, which checks all network Nodes for Prometheus metrics.
// The metric for which this instance is registered is captured and stored in time series separately for each Node.
// This source will represent a new metric, which will have the same name as the metric to get from prometheus.
// If the prometheus metric has quantile, the suffix '_q<num>', e.g. '_q0.999', will be added to the new metric name.
func NewPromLogSource(monitor *monitoring.Monitor, prometheusMetric monitoring.PrometheusLogKey) *PromLogSource {
	p := &PromLogSource{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(toMetric(prometheusMetric)),
	}
	monitor.PrometheusLogProvider().RegisterLogListener(prometheusMetric, p)
	return p
}

func (p *PromLogSource) OnLog(node monitoring.Node, time monitoring.Time, value float64) {
	series := p.GetOrAddSubject(node)
	if err := series.Append(time, value); err != nil {
		log.Printf("cannot add to series: %s", err)
	}
}

func toMetric(prometheusMetric monitoring.PrometheusLogKey) monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.Time, float64]] {
	name := prometheusMetric.Name
	if prometheusMetric.Quantile != monitoring.QuantileEmpty {
		name = fmt.Sprintf("%s_q%s", name, prometheusMetric.Quantile)
	}
	return monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.Time, float64]]{
		Name:        name,
		Description: fmt.Sprintf("Prometheus metric for %s", name),
	}
}

package monitoring

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// PrometheusLogValue is one measured value obtained from Prometheus.
// It contains the measured value itself, and  the metrics' categorisation
// falling into one of the: counter, gauge, summary.
// For the summary, the underlying value can come from the Meter or Timer. If the metric is summary,
// it contains additionally percentile.
// For more information about metrics type, the reader can have a look at: https://geth.ethereum.org/docs/monitoring/metrics
type PrometheusLogValue struct {
	PrometheusLogKey
	metricType PrometheusMetricType
	value      float64
}

func (p PrometheusLogValue) String() string {
	if p.metricType == summaryPrometheusMetricType {
		return fmt.Sprintf("%s_%s: %.2f (q: %.5f)", p.name, p.metricType, p.value, p.quantile)
	} else {
		return fmt.Sprintf("%s_%s: %.2f", p.name, p.metricType, p.value)
	}
}

// PrometheusMetricType is one of the counter, gauge, summary.
type PrometheusMetricType string

const (
	counterPrometheusMetricType PrometheusMetricType = "counter" // a counter can be increased or decreased
	gaugePrometheusMetricType                        = "gauge"   // a gauge works as the counter, but can be also set to a direct value
	summaryPrometheusMetricType                      = "summary" // summary is either a Meter or a Timer, they both measure throughput
)

// ParsePrometheusLogReader reads text logs from the input reader, and produces the output slice of parsed representation of the log.
// The reader is expected to contain a Prometheus TXT Log stream,
// which is parsed and converted into PrometheusLogValue structs.
// The prometheus log format is described at: https://prometheus.io/docs/instrumenting/exposition_formats/
// however, Opera/Geth uses only a subset of metric types: https://geth.ethereum.org/docs/monitoring/metrics
// In brief, the TXT format is as follows:
// A line can start with a comment: '# TYPE <metric_name> <metric_type>',
// which describes the name and type of the metric. The type is one of 'counter', 'gauge' or 'summary'.
// Next lines must follow for the metrics with the same name, and be one of the two forms:
// either: <metric_name> {quantile="<val>"} <value>
// or: <metric_name> <value>
// The variant with the quantile is used only when the metric type is 'summary', and in this case
// measured values are categorised into quantiles.
// For both variants, the value contains an actual value measured for this particular metric.
// A line starting with the `# TYPE` marks the beginning of a new metric.
func ParsePrometheusLogReader(reader io.Reader) ([]PrometheusLogValue, error) {
	res := make([]PrometheusLogValue, 0, 2000)
	scanner := bufio.NewScanner(reader)
	var nextType PrometheusMetricType
	var currentName string
	var errs []error
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) >= 2 {
			if tokens[1] == "TYPE" {
				currentName = tokens[2]
				nextType = PrometheusMetricType(tokens[3])
			} else if tokens[0] == currentName {
				val := PrometheusLogValue{PrometheusLogKey: PrometheusLogKey{name: currentName}, metricType: nextType}
				if err := fillValue(tokens, &val); err != nil {
					errs = append(errs, err)
				} else {
					res = append(res, val)
				}
			} else {
				errs = append(errs, fmt.Errorf("unecpected line starting with, %s -> %s", currentName, tokens))
			}

		}
	}

	return res, errors.Join(scanner.Err(), errors.Join(errs...))
}

var (
	quantileReg = regexp.MustCompile(`.*{quantile=".*"}`)
)

// fillValue analyses input tokens and stores values of quantile and the metrics value
// into the PrometheusLogValue.
func fillValue(tokens []string, dest *PrometheusLogValue) error {
	var valueStr string
	// the format of the array is either:
	// - metric_name metric_value
	// - metric_name quantile_value metric_value
	if len(tokens) >= 3 && quantileReg.MatchString(tokens[1]) {
		q, err := strconv.ParseFloat(strings.Split(tokens[1], "\"")[1], 32)
		if err != nil {
			return fmt.Errorf("cannot parse quantile from: %s", tokens[1])
		}
		dest.quantile = float32(q)
		valueStr = tokens[2]
	} else {
		valueStr = tokens[1]
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("cannot parse value from: %s", valueStr)
	}

	dest.value = value

	return nil
}

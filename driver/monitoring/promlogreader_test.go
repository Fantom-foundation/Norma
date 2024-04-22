// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package monitoring

import (
	"strings"
	"testing"
)

func TestParsePrometheusLog(t *testing.T) {
	prometheusValues, err := ParsePrometheusLogReader(strings.NewReader(testTxtPrometheusLog))
	if err != nil {
		t.Errorf("error: %s", err)
	}

	var count int
	for i, value := range prometheusValues {
		if got, want := value, expectedValues[i]; got != want {
			t.Errorf("prometheus txt log value does not match: %v ! %v", got, want)
		}
		count++
	}

	if got, want := len(prometheusValues), len(expectedValues); got != want {
		t.Errorf("not all logs entries were extracted: %d != %d", got, want)
	}
}

func TestParseCorruptedLineShouldNotPanic(t *testing.T) {
	txt :=
		"c o r r u p t e d\n" +
			"# TYPE chain_execution summary\n" +
			"chain_execution a b c d e f\n" +
			"chain_execution\n" +
			"chain_execution {quantile=\"0.95\"}\n"

	prometheusValues, err := ParsePrometheusLogReader(strings.NewReader(txt))
	if err == nil {
		t.Errorf("parsing should fail")
	}

	var count int
	for range prometheusValues {
		count++
	}

	if len(prometheusValues) != 0 {
		t.Errorf("should read no data")
	}
}

var (
	// test log
	testTxtPrometheusLog = "# TYPE chain_execution_count counter\n" +
		"chain_execution_count 231\n" +
		"\n" +
		"# TYPE chain_execution summary\n" +
		"chain_execution {quantile=\"0.5\"} 7.8168292e+07\n" +
		"chain_execution {quantile=\"0.75\"} 1.22758041e+08\n" +
		"chain_execution {quantile=\"0.95\"} 1.991017745999998e+08\n" +
		"chain_execution {quantile=\"0.99\"} 3.821057185599999e+08\n" +
		"chain_execution {quantile=\"0.999\"} 4.4993e+08\n" +
		"chain_execution {quantile=\"0.9999\"} 4.4993e+08\n" +
		"\n" +
		"# TYPE chain_block_age gauge\n" +
		"chain_block_age 815212292\n"

	expectedValues = []PrometheusLogValue{
		{PrometheusLogKey{"chain_execution_count", ""}, counterPrometheusMetricType, 231},
		{PrometheusLogKey{"chain_execution", "0.5"}, summaryPrometheusMetricType, 7.8168292e+07},
		{PrometheusLogKey{"chain_execution", "0.75"}, summaryPrometheusMetricType, 1.22758041e+08},
		{PrometheusLogKey{"chain_execution", "0.95"}, summaryPrometheusMetricType, 1.991017745999998e+08},
		{PrometheusLogKey{"chain_execution", "0.99"}, summaryPrometheusMetricType, 3.821057185599999e+08},
		{PrometheusLogKey{"chain_execution", "0.999"}, summaryPrometheusMetricType, 4.4993e+08},
		{PrometheusLogKey{"chain_execution", "0.9999"}, summaryPrometheusMetricType, 4.4993e+08},
		{PrometheusLogKey{"chain_block_age", ""}, gaugePrometheusMetricType, 815212292},
	}
)

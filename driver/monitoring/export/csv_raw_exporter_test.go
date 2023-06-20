package export_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/app"
	"github.com/Fantom-foundation/Norma/driver/monitoring/export"
	netmon "github.com/Fantom-foundation/Norma/driver/monitoring/network"
	nodemon "github.com/Fantom-foundation/Norma/driver/monitoring/node"
	"golang.org/x/exp/constraints"
)

func TestPrintMultiSourceMultiSectionCsvRows(t *testing.T) {
	time1, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:15.080]")
	time2, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:15.537]")
	time3, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:16.027]")
	time4, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:16.512]")
	time5, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:17.003]")
	time6, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:38:15.080]")

	now := time.Now()
	time1y := time.Date(now.Year(), time1.Month(), time1.Day(), time1.Hour(), time1.Minute(), time1.Second(), time1.Nanosecond(), time.UTC)
	time2y := time.Date(now.Year(), time2.Month(), time2.Day(), time2.Hour(), time2.Minute(), time2.Second(), time2.Nanosecond(), time.UTC)
	time3y := time.Date(now.Year(), time3.Month(), time3.Day(), time3.Hour(), time3.Minute(), time3.Second(), time3.Nanosecond(), time.UTC)
	time4y := time.Date(now.Year(), time4.Month(), time4.Day(), time4.Hour(), time4.Minute(), time4.Second(), time4.Nanosecond(), time.UTC)
	time5y := time.Date(now.Year(), time5.Month(), time5.Day(), time5.Hour(), time5.Minute(), time5.Second(), time5.Nanosecond(), time.UTC)
	time6y := time.Date(now.Year(), time6.Month(), time6.Day(), time6.Hour(), time6.Minute(), time6.Second(), time6.Nanosecond(), time.UTC)

	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Time]{}
	_ = s1.Append(monitoring.BlockNumber(1), time1y)
	_ = s1.Append(monitoring.BlockNumber(2), time2y)
	_ = s1.Append(monitoring.BlockNumber(3), time3y)

	s2 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Time]{}
	_ = s2.Append(monitoring.BlockNumber(1), time4y)
	_ = s2.Append(monitoring.BlockNumber(2), time5y)
	_ = s2.Append(monitoring.BlockNumber(3), time6y)

	s3 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s3.Append(monitoring.BlockNumber(1), 10*time.Second)
	_ = s3.Append(monitoring.BlockNumber(2), 20*time.Second)
	_ = s3.Append(monitoring.BlockNumber(3), 30*time.Second)

	s4 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s4.Append(monitoring.BlockNumber(1), 15*time.Second)
	_ = s4.Append(monitoring.BlockNumber(2), 25*time.Second)
	_ = s4.Append(monitoring.BlockNumber(3), 30*time.Second)

	s5 := &monitoring.SyncedSeries[monitoring.Time, int]{}
	_ = s5.Append(monitoring.NewTime(time1y), 110)
	_ = s5.Append(monitoring.NewTime(time2y), 120)

	s6 := &monitoring.SyncedSeries[monitoring.BlockNumber, int]{}
	_ = s6.Append(monitoring.BlockNumber(1), 17)
	_ = s6.Append(monitoring.BlockNumber(2), 21)
	_ = s6.Append(monitoring.BlockNumber(3), 35)

	s7 := &monitoring.SyncedSeries[monitoring.Time, int]{}
	_ = s7.Append(monitoring.NewTime(time1y), 11)
	_ = s7.Append(monitoring.NewTime(time2y), 12)

	s8 := &monitoring.SyncedSeries[monitoring.Time, int]{}
	_ = s8.Append(monitoring.NewTime(time1y), 11)
	_ = s8.Append(monitoring.NewTime(time2y), 13)

	s9 := &monitoring.SyncedSeries[int, int]{}
	_ = s9.Append(100, 110)
	_ = s9.Append(200, 213)

	n1 := monitoring.Node("A")
	n2 := monitoring.Node("B")

	// section 1
	source1 := newSource[monitoring.Node, monitoring.BlockNumber, time.Time, monitoring.BlockSeries[time.Time]](nodemon.BlockCompletionTime)
	source1.put(n1, s1)
	source1.put(n2, s2)

	// section 1 - next column
	source2 := newSource[monitoring.Node, monitoring.BlockNumber, time.Duration, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source2.put(n1, s3)
	source2.put(n2, s4)

	// section 2
	source3 := newSource[monitoring.Network, monitoring.Time, int, monitoring.TimeSeries[int]](netmon.NumberOfNodes)
	source3.put(monitoring.Network{}, s5)

	// section 3
	source4 := newSource[monitoring.Network, monitoring.BlockNumber, int, monitoring.BlockSeries[int]](netmon.BlockNumberOfTransactions)
	source4.put(monitoring.Network{}, s6)

	// section 4
	source5 := newSource[monitoring.Node, monitoring.Time, int, monitoring.TimeSeries[int]](nodemon.NodeBlockHeight)
	source5.put(n1, s7)
	source5.put(n2, s8)

	// section 5
	source6 := newSource[monitoring.App, int, int, monitoring.Series[int, int]](app.ReceivedTransactions)
	source6.put("app-1", s9)

	// add in the CSV
	var builder strings.Builder
	csv := monitoring.NewWriterChain(&stringBuilderCloser{&builder})
	_ = export.AddNodeBlockSeriesSource[time.Time](csv, source1, export.TimeConverter{})
	_ = export.AddNodeBlockSeriesSource[time.Duration](csv, source2, export.DurationConverter{})
	_ = export.AddNetworkBlockSeriesSource[int](csv, source4, export.DirectConverter[int]{})
	_ = export.AddNetworkTimeSeriesSource[int](csv, source3, export.DirectConverter[int]{})
	_ = export.AddNodeTimeSeriesSource[int](csv, source5, export.DirectConverter[int]{})
	_ = export.AddAppSeriesSource[int, int](csv, source6, export.DirectConverter[int]{}, export.DirectConverter[int]{})
	_ = csv.Close()

	expected :=
		"BlockCompletionTime, network, A, , , 1, , 1683192855080000000\n" +
			"BlockCompletionTime, network, A, , , 2, , 1683192855537000000\n" +
			"BlockCompletionTime, network, A, , , 3, , 1683192856027000000\n" +
			"BlockCompletionTime, network, B, , , 1, , 1683192856512000000\n" +
			"BlockCompletionTime, network, B, , , 2, , 1683192857003000000\n" +
			"BlockCompletionTime, network, B, , , 3, , 1683193095080000000\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 1, , 10000000000\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 2, , 20000000000\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 3, , 30000000000\n" +
			"BlockEventAndTxsProcessingTime, network, B, , , 1, , 15000000000\n" +
			"BlockEventAndTxsProcessingTime, network, B, , , 2, , 25000000000\n" +
			"BlockEventAndTxsProcessingTime, network, B, , , 3, , 30000000000\n" +
			"BlockNumberOfTransactions, network, , , , 1, , 17\n" +
			"BlockNumberOfTransactions, network, , , , 2, , 21\n" +
			"BlockNumberOfTransactions, network, , , , 3, , 35\n" +
			"NumberOfNodes, network, , , 1683192855080000000, , , 110\n" +
			"NumberOfNodes, network, , , 1683192855537000000, , , 120\n" +
			"NodeBlockHeight, network, A, , 1683192855080000000, , , 11\n" +
			"NodeBlockHeight, network, A, , 1683192855537000000, , , 12\n" +
			"NodeBlockHeight, network, B, , 1683192855080000000, , , 11\n" +
			"NodeBlockHeight, network, B, , 1683192855537000000, , , 13\n" +
			"ReceivedTransactions, network, , app-1, , , 100, 110\n" +
			"ReceivedTransactions, network, , app-1, , , 200, 213\n"

	if !strings.Contains(builder.String(), expected) {
		t.Errorf("strings do not match:\n %s \n is not \n%s", expected, builder.String())
	}
}

func TestRegisterSources(t *testing.T) {
	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s1.Append(monitoring.BlockNumber(1), 10*time.Second)
	_ = s1.Append(monitoring.BlockNumber(2), 20*time.Second)
	_ = s1.Append(monitoring.BlockNumber(3), 30*time.Second)

	n1 := monitoring.Node("A")
	source := newSource[monitoring.Node, monitoring.BlockNumber, time.Duration, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source.put(n1, s1)

	// register twice the same
	var builder strings.Builder
	csv := monitoring.NewWriterChain(&stringBuilderCloser{&builder})
	csv.Add(func() error {
		return export.AddNodeBlockSeriesSource[time.Duration](csv, source, export.DirectConverter[time.Duration]{})
	})
	csv.Add(func() error {
		return export.AddNodeBlockSeriesSource[time.Duration](csv, source, export.DurationConverter{})
	})
	_ = csv.Close()

	expected :=
		"BlockEventAndTxsProcessingTime, network, A, , , 1, , 10s\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 2, , 20s\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 3, , 30s\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 1, , 10000000000\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 2, , 20000000000\n" +
			"BlockEventAndTxsProcessingTime, network, A, , , 3, , 30000000000\n"

	if !strings.Contains(builder.String(), expected) {
		t.Errorf("strings do not match:\n %s \n is not \n%s", expected, builder.String())
	}
}

type source[S comparable, K constraints.Ordered, T any, X monitoring.Series[K, T]] struct {
	metric monitoring.Metric[S, X]
	data   map[S]X
}

func newSource[S comparable, K constraints.Ordered, T any, X monitoring.Series[K, T]](metric monitoring.Metric[S, X]) *source[S, K, T, X] {
	return &source[S, K, T, X]{
		data:   make(map[S]X, 10),
		metric: metric,
	}
}

func (d *source[S, K, T, X]) GetSubjects() []S {
	subjects := make([]S, 0, len(d.data))
	for k := range d.data {
		subjects = append(subjects, k)
	}
	return subjects
}

func (d *source[S, K, T, X]) GetMetric() monitoring.Metric[S, X] {
	return d.metric
}

func (d *source[S, K, T, X]) GetData(subject S) (X, bool) {
	res, exists := d.data[subject]
	return res, exists
}

func (d *source[S, K, T, X]) Shutdown() error {
	return nil
}

func (d *source[S, K, T, X]) put(s S, val X) {
	d.data[s] = val
}

type stringBuilderCloser struct {
	builder *strings.Builder
}

func (s *stringBuilderCloser) Write(b []byte) (int, error) {
	return s.builder.Write(b)
}

func (s *stringBuilderCloser) Close() error {
	// noop
	return nil
}

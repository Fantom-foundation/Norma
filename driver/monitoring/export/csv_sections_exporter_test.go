package export

import (
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	netmon "github.com/Fantom-foundation/Norma/driver/monitoring/network"
	nodemon "github.com/Fantom-foundation/Norma/driver/monitoring/node"
	"strings"
	"testing"
	"time"
)

func TestPrintMultiSourceMultiSectionCsv(t *testing.T) {
	time1, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:15.080]")
	time2, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:15.537]")
	time3, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:16.027]")
	time4, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:16.512]")
	time5, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:34:17.003]")
	time6, _ := time.Parse("[01-02|15:04:05.000]", "[05-04|09:38:15.080]")

	now := time.Now()
	time1y := time.Date(now.Year(), time1.Month(), time1.Day(), time1.Hour(), time1.Minute(), time1.Second(), time1.Nanosecond(), time.UTC)
	time2y := time.Date(now.Year(), time2.Month(), time2.Day(), time2.Hour(), time2.Minute(), time2.Second(), time2.Nanosecond(), time.UTC)

	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Time]{}
	_ = s1.Append(monitoring.BlockNumber(1), time1)
	_ = s1.Append(monitoring.BlockNumber(2), time2)
	_ = s1.Append(monitoring.BlockNumber(3), time3)

	s2 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Time]{}
	_ = s2.Append(monitoring.BlockNumber(1), time4)
	_ = s2.Append(monitoring.BlockNumber(2), time5)
	_ = s2.Append(monitoring.BlockNumber(3), time6)

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

	n1 := monitoring.Node("A")
	n2 := monitoring.Node("B")

	// section 1
	source1 := newSource[monitoring.Node, monitoring.BlockSeries[time.Time]](nodemon.BlockCompletionTime)
	source1.put(n1, s1)
	source1.put(n2, s2)

	// section 1 - next column
	source2 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source2.put(n1, s3)
	source2.put(n2, s4)

	// section 2
	source3 := newSource[monitoring.Network, monitoring.TimeSeries[int]](netmon.NumberOfNodes)
	source3.put(monitoring.Network{}, s5)

	// section 3
	source4 := newSource[monitoring.Network, monitoring.BlockSeries[int]](netmon.BlockNumberOfTransactions)
	source4.put(monitoring.Network{}, s6)

	// section 4
	source5 := newSource[monitoring.Node, monitoring.TimeSeries[int]](nodemon.NodeBlockHeight)
	source5.put(n1, s7)
	source5.put(n2, s8)

	// add in the CSV
	var builder strings.Builder
	csv := NewCsvExporter(&stringBuilderCloser{&builder})
	AddEmptySection(csv, 3)
	AddSection(csv)
	AddNodeBlockSeriesSource[time.Time](csv, source1, TimeConverter{})
	AddNodeBlockSeriesSource[time.Duration](csv, source2, DurationConverter{})
	AddNetworkBlockSeriesSource[int](csv, source4, DirectConverter[int]{})
	AddEmptySection(csv, 1)
	AddSection(csv)
	AddNetworkTimeSeriesSource[int](csv, source3, DirectConverter[int]{})
	AddNodeTimeSeriesSource[int](csv, source5, DirectConverter[int]{})
	_ = csv.Flush()

	expected := ",\n,\n,\n" +
		", BlockCompletionTime, , BlockEventAndTxsProcessingTime, , BlockNumberOfTransactions,\n" +
		"monitoring.BlockNumber, A, B, A, B, {},\n" +
		"1, 09:34:15.08, 09:34:16.512, 10000000000, 15000000000, 17,\n" +
		"2, 09:34:15.537, 09:34:17.003, 20000000000, 25000000000, 21,\n" +
		"3, 09:34:16.027, 09:38:15.08, 30000000000, 30000000000, 35,\n" +
		",\n" +
		", NumberOfNodes, NodeBlockHeight, ,\n" +
		"monitoring.Time, {}, A, B,\n" +
		"09:34:15.08, 110, 11, 11,\n" +
		"09:34:15.537, 120, 12, 13,\n"

	if expected != builder.String() {
		t.Errorf("strings do not match:\n %s \n is not \n %s", expected, builder.String())
	}
}

func TestSeriesSizesDoNotMatch(t *testing.T) {
	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s1.Append(monitoring.BlockNumber(1), 10*time.Second)
	_ = s1.Append(monitoring.BlockNumber(2), 20*time.Second)
	_ = s1.Append(monitoring.BlockNumber(3), 30*time.Second)

	s2 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s2.Append(monitoring.BlockNumber(1), 11*time.Second)

	n1 := monitoring.Node("A")
	n2 := monitoring.Node("B")

	source1 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source1.put(n1, s1)
	source1.put(n2, s2)

	var builder strings.Builder
	csv := NewCsvExporter(&stringBuilderCloser{&builder})
	AddSection(csv)
	AddNodeBlockSeriesSource[time.Duration](csv, source1, DurationConverter{})
	_ = csv.Flush()

	expected := ", BlockEventAndTxsProcessingTime, ,\n" +
		"monitoring.BlockNumber, A, B,\n" +
		"1, 10000000000, 11000000000,\n" +
		"2, 20000000000, ,\n" +
		"3, 30000000000, ,\n"

	if expected != builder.String() {
		t.Errorf("strings do not match:\n %s \n is not \n %s", expected, builder.String())
	}
}

func TestSeriesKeyDoNotMatch(t *testing.T) {
	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s1.Append(monitoring.BlockNumber(1), 10*time.Second)
	_ = s1.Append(monitoring.BlockNumber(2), 20*time.Second)
	_ = s1.Append(monitoring.BlockNumber(3), 30*time.Second)

	s2 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s2.Append(monitoring.BlockNumber(10), 11*time.Second)
	_ = s2.Append(monitoring.BlockNumber(20), 12*time.Second)
	_ = s2.Append(monitoring.BlockNumber(30), 13*time.Second)

	n1 := monitoring.Node("A")
	n2 := monitoring.Node("B")

	source1 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source1.put(n1, s1)
	source1.put(n2, s2)

	var builder strings.Builder
	csv := NewCsvExporter(&stringBuilderCloser{&builder})
	AddSection(csv)
	AddNodeBlockSeriesSource[time.Duration](csv, source1, DurationConverter{})
	_ = csv.Flush()

	expected := ", BlockEventAndTxsProcessingTime, , ,\n" +
		"monitoring.BlockNumber, A, monitoring.BlockNumber, B,\n" +
		"1, 10000000000, 10, 11000000000,\n" +
		"2, 20000000000, 20, 12000000000,\n" +
		"3, 30000000000, 30, 13000000000,\n"

	if expected != builder.String() {
		t.Errorf("strings do not match:\n %s \n is not \n %s", expected, builder.String())
	}
}

func TestSubjectSeriesKeyDoNotMatch(t *testing.T) {
	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s1.Append(monitoring.BlockNumber(1), 10*time.Second)
	_ = s1.Append(monitoring.BlockNumber(2), 20*time.Second)
	_ = s1.Append(monitoring.BlockNumber(3), 30*time.Second)

	s2 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s2.Append(monitoring.BlockNumber(10), 11*time.Second)
	_ = s2.Append(monitoring.BlockNumber(20), 12*time.Second)
	_ = s2.Append(monitoring.BlockNumber(30), 13*time.Second)

	s3 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s3.Append(monitoring.BlockNumber(10), 21*time.Second)
	_ = s3.Append(monitoring.BlockNumber(20), 22*time.Second)

	s4 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s4.Append(monitoring.BlockNumber(20), 31*time.Second)
	_ = s4.Append(monitoring.BlockNumber(30), 32*time.Second)

	n1 := monitoring.Node("A")

	source1 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source1.put(n1, s1)

	source2 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source2.put(n1, s2)

	source3 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source3.put(n1, s3)

	source4 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source4.put(n1, s4)

	var builder strings.Builder
	csv := NewCsvExporter(&stringBuilderCloser{&builder})
	AddSection(csv)
	AddNodeBlockSeriesSource[time.Duration](csv, source1, DurationConverter{})
	AddNodeBlockSeriesSource[time.Duration](csv, source2, DurationConverter{})
	AddNodeBlockSeriesSource[time.Duration](csv, source3, DurationConverter{})
	AddNodeBlockSeriesSource[time.Duration](csv, source4, DurationConverter{})
	_ = csv.Flush()

	expected := ", BlockEventAndTxsProcessingTime, , BlockEventAndTxsProcessingTime, BlockEventAndTxsProcessingTime, BlockEventAndTxsProcessingTime,\n" +
		"monitoring.BlockNumber, A, monitoring.BlockNumber, A, A, A,\n" +
		"1, 10000000000, 10, 11000000000, 21000000000,\n" +
		"2, 20000000000, 20, 12000000000, 22000000000, 31000000000,\n" +
		"3, 30000000000, 30, 13000000000, , 32000000000,\n"

	if expected != builder.String() {
		t.Errorf("strings do not match:\n %s \n is not \n %s", expected, builder.String())
	}
}

func TestSeriesKeyShift(t *testing.T) {
	s1 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s1.Append(monitoring.BlockNumber(1), 10*time.Second)
	_ = s1.Append(monitoring.BlockNumber(2), 20*time.Second)
	_ = s1.Append(monitoring.BlockNumber(3), 30*time.Second)
	_ = s1.Append(monitoring.BlockNumber(4), 40*time.Second)
	_ = s1.Append(monitoring.BlockNumber(5), 50*time.Second)

	s2 := &monitoring.SyncedSeries[monitoring.BlockNumber, time.Duration]{}
	_ = s2.Append(monitoring.BlockNumber(3), 33*time.Second)
	_ = s2.Append(monitoring.BlockNumber(5), 55*time.Second)

	n1 := monitoring.Node("A")
	n2 := monitoring.Node("B")

	source1 := newSource[monitoring.Node, monitoring.BlockSeries[time.Duration]](nodemon.BlockEventAndTxsProcessingTime)
	source1.put(n1, s1)
	source1.put(n2, s2)

	var builder strings.Builder
	csv := NewCsvExporter(&stringBuilderCloser{&builder})
	AddSection(csv)
	AddNodeBlockSeriesSource[time.Duration](csv, source1, DurationConverter{})
	_ = csv.Flush()

	expected := ", BlockEventAndTxsProcessingTime, ,\n" +
		"monitoring.BlockNumber, A, B,\n" +
		"1, 10000000000,\n" +
		"2, 20000000000,\n" +
		"3, 30000000000, 33000000000,\n" +
		"4, 40000000000,\n" +
		"5, 50000000000, 55000000000,\n"

	if expected != builder.String() {
		t.Errorf("strings do not match:\n %s \n is not \n %s", expected, builder.String())
	}
}

type source[S comparable, T any] struct {
	metric monitoring.Metric[S, T]
	data   map[S]T
}

func newSource[S comparable, T any](metric monitoring.Metric[S, T]) *source[S, T] {
	return &source[S, T]{
		data:   make(map[S]T, 10),
		metric: metric,
	}
}

func (d *source[S, T]) GetSubjects() []S {
	subjects := make([]S, 0, len(d.data))
	for k := range d.data {
		subjects = append(subjects, k)
	}
	return subjects
}

func (d *source[S, T]) GetMetric() monitoring.Metric[S, T] {
	return d.metric
}

func (d *source[S, T]) GetData(subject S) (T, bool) {
	res, exists := d.data[subject]
	return res, exists
}

func (d *source[S, T]) Shutdown() error {
	return nil
}

func (d *source[S, T]) put(s S, val T) {
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

package export

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"golang.org/x/exp/constraints"
)

// CsvRecord summarizes the content of a single line in the exported raw CSV
// metric dump. It is used as the output format of data sources when exporting
// data to facilitate future extensions / modifications of the format.
type CsvRecord struct {
	Metric              string // must not be empty
	Network, Node, App  string // may be empty
	Time, Block, Worker *int64 // one must be set
	Value               string // must not be empty!
}

// WriteCsvHeader writes a header line defining the fields of the CsvRecord.
func WriteCsvHeader(out io.Writer) error {
	_, err := out.Write([]byte("metric,network,node,app,time,block,workers,value\n"))
	return err
}

// WriteTo writes the content of a CsvRecord to the given writer.
func (r *CsvRecord) WriteTo(out io.Writer) (int64, error) {
	toStr := func(x *int64) string {
		if x == nil {
			return ""
		}
		return fmt.Sprintf("%d", *x)
	}

	line := strings.Join([]string{
		r.Metric,
		r.Network,
		r.Node,
		r.App,
		toStr(r.Time),
		toStr(r.Block),
		toStr(r.Worker),
		r.Value,
	}, ", ") + "\n"
	n, err := out.Write([]byte(line))
	return int64(n), err
}

// setSubject sets the subject fields in the row to the given value.
func (r *CsvRecord) setSubject(subject any) *CsvRecord {
	r.Network = "network"
	switch value := subject.(type) {
	case monitoring.Network:
		// nothing to do
	case monitoring.Node:
		r.Node = string(value)
	case monitoring.App:
		r.App = string(value)
	case monitoring.Account:
		r.App = string(value.App)
		var worker int64 = int64(value.Id)
		r.Worker = &worker
	default:
		panic(fmt.Sprintf("unsupported subject value encountered: %v (type: %v)", subject, reflect.TypeOf(subject)))
	}
	return r
}

// setKey sets the key fields in the row to the given value.
func (r *CsvRecord) setKey(key any) *CsvRecord {
	switch value := key.(type) {
	case monitoring.BlockNumber:
		block := int64(value)
		r.Block = &block
	case monitoring.Time:
		time := int64(value.Time().UTC().UnixNano())
		r.Time = &time
	case int:
		worker := int64(value)
		r.Worker = &worker
	default:
		panic(fmt.Sprintf("unsupported key value encountered: %v (type: %v)", key, reflect.TypeOf(key)))
	}
	return r
}

// setValue sets the value field in the row to the given value.
func (r *CsvRecord) setValue(value any) *CsvRecord {
	switch v := value.(type) {
	case int:
		r.Value = fmt.Sprintf("%d", v)
	case float32:
		r.Value = fmt.Sprintf("%v", v)
	case string:
		r.Value = v
	case time.Time:
		r.Value = fmt.Sprintf("%d", v.UTC().UnixNano())
	case time.Duration:
		r.Value = fmt.Sprintf("%d", v.Nanoseconds())
	default:
		panic(fmt.Sprintf("unsupported value encountered: %v (type: %v)", value, reflect.TypeOf(value)))
	}
	return r
}

// AddSeriesData writes all data from the given source to the given writer in CSV format.
func AddSeriesData[S any, K constraints.Ordered, T any](
	exporter io.Writer,
	source monitoring.Source[S, monitoring.Series[K, T]],
) error {
	forEacher := monitoring.NewSourceRowsForEacher[S, K, T, monitoring.Series[K, T]](source)
	var errs []error
	forEacher.ForEachRow(func(row monitoring.Row[S, K, T, monitoring.Series[K, T]]) {
		line := CsvRecord{Metric: row.Metric.Name}
		line.setSubject(row.Subject).setKey(row.Position).setValue(row.Value)
		if _, err := line.WriteTo(exporter); err != nil {
			errs = append(errs, err)
		}
	})

	return errors.Join(errs...)
}

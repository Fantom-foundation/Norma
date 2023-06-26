package monitoring

import (
	"fmt"
	"io"
	"strings"
)

// CsvRecord summarizes the content of a single line in the exported raw CSV
// metric dump. It is used as the output format of data sources when exporting
// data to facilitate future extensions / modifications of the format.
type CsvRecord struct {
	Record
	Metric string // must not be empty
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

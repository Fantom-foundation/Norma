package export

import (
	"errors"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"golang.org/x/exp/constraints"
	"io"
	"sort"
	"strings"
)

// CsvExporter allows for accumulating data in a form of CSV lines and  export them to the output writer.
// The CSV is composed of monitoring.Source instances, which can be organised into sections.
// It is expected that the sources contain monitoring.Series as its values. No other types
// are supported at the moment.
// Every section can include many sources and subjects as long as series' keys (i.e. the X-axis) are compatible.
// The series keys from the first series are printed in the first column representing the X-axis, and series
// values (i.e. the Y-axis) are printed in the next column. Further series values are printed in consecutive columns,
// while keys are ignored. In other worlds, it is assumed that next series have compatible keys with the already printed
// first column.
// Note that it is on purpose not checked that the series in the section have all exactly the same keys.
// It is up to the user to insert comparable data. For instance, when a sampling is performed,
// it is hardly possible to have exactly same timestamps at all series, but as long as the sampling period
// was set the same for all series, they may be assumed compatible.
// Furthermore, if the series following the first series are longer or shorter than the first one,
// the series are is either trimmed or not printed fully.
// Finally, every section has a header with the title of the first column matching the X-axis type and subject + metric
// name being the header of next (Y-axis) columns.
type CsvExporter struct {
	writer      io.WriteCloser
	lines       []string
	firstColumn bool
	currentLine int
	seriesSize  int
	errors      []error
}

// NewCsvExporter creates a new CSV exporter that writes the CSV lines into the input writer.
func NewCsvExporter(writer io.WriteCloser) *CsvExporter {
	return &CsvExporter{
		writer:      writer,
		lines:       make([]string, 0, 1000),
		errors:      make([]error, 0, 5),
		firstColumn: true,
	}
}

// AddEmptySection inserts the input amount of empty lines
func AddEmptySection(exporter *CsvExporter, lines int) {
	if err := exporter.flush(); err != nil {
		exporter.errors = append(exporter.errors, err)
	}

	for i := 0; i < lines; i++ {
		exporter.append(",")
	}

	exporter.currentLine = 0
}

// AddSection starts a new section, which will produce lines below the current section.
// This method has to be called before consecutive calls to AddSource().
// The section is typed to the type of the source's series key, which will be print
// as the first column header.
func AddSection[K constraints.Ordered](exporter *CsvExporter) {
	if err := exporter.flush(); err != nil {
		exporter.errors = append(exporter.errors, err)
	}

	// add header line
	exporter.append("")
	var firstColumn K
	exporter.append(fmt.Sprintf("%T", firstColumn))

	exporter.currentLine = 0
}

// AddSource adds lines for the input source's series. It is expected that the source's values
// are monitoring.Series. No other type of value is supported.
// If this is the first source added in this section, the series keys (i.e. X-axis) are printed
// in the first column. The values are printed next to it in the following column
// (i.e. the Y-axis). If further sources are added, only the series' values are printed to next columns.
// It means that it is expected that all sources added in this section have matching series' keys (X-axis).
// If the series do not have exactly the same keys, they should be at least comparable. For instance,
// sampling may produce different timestamps for every series, but it may be considered fine as long as the sampling
// period is the same.
func AddSource[S any, K constraints.Ordered, T any, X monitoring.Series[K, T]](
	exporter *CsvExporter,
	source monitoring.Source[S, X],
	kConverter converter[K],
	tConverter converter[T],
	cmp comparator[S],
) {

	subjects := source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return cmp.Compare(&subjects[i], &subjects[j]) < 0 })
	subjectsCount := len(subjects)

	// add headers - metric name first + sub-subHeader - subject names
	var subHeader, header strings.Builder
	for i, subject := range subjects {
		if i == 0 {
			header.WriteString(fmt.Sprintf(",%v", source.GetMetric().Name))
		} else {
			header.WriteString(",")
		}
		subHeader.WriteString(fmt.Sprintf(",%v", subject))
	}
	exporter.append(header.String())
	exporter.append(subHeader.String())

	// get the sizes, i.e. the number of data points
	// store in order to arrays as subjects are sorted
	data := make([]monitoring.Series[K, T], 0, subjectsCount)
	sizes := make([]int, 0, subjectsCount)
	var longestSeriesIndex int
	var size int
	for i, subject := range subjects {
		series, exists := source.GetData(subject)
		if !exists {
			panic(fmt.Errorf("series for subject: %v does not exist", subject))
		}
		data = append(data, series)
		currentSize := series.Size()
		if currentSize > size {
			size = currentSize
			longestSeriesIndex = i
		}
		sizes = append(sizes, currentSize)
	}
	if exporter.firstColumn {
		exporter.seriesSize = size
	}

	// continue with data
	for i := 0; i < exporter.seriesSize; i++ {
		var line strings.Builder
		position := data[longestSeriesIndex].GetAt(i).Position
		if exporter.firstColumn {
			line.WriteString(kConverter.Convert(position))
		}
		for j := 0; j < subjectsCount; j++ {
			if i >= sizes[j] {
				// series is shorter, replace with empty cells.
				line.WriteString(",")
			} else {
				point := data[j].GetAt(i)
				val := point.Value
				line.WriteString(fmt.Sprintf(", %s", tConverter.Convert(val)))
			}
		}

		exporter.append(line.String())
	}

	exporter.currentLine = 0
	exporter.firstColumn = false
}

// AddNodeBlockSeriesSource is a shortcut for AddSource(), which adds a source where the subject
// is a monitoring.Node and the series' key is monitoring.BlockNumber
func AddNodeBlockSeriesSource[T any](
	exporter *CsvExporter,
	source monitoring.Source[monitoring.Node, monitoring.BlockSeries[T]],
	yConverter converter[T],
) {
	AddSource[monitoring.Node, monitoring.BlockNumber, T, monitoring.BlockSeries[T]](exporter, source, DirectConverter[monitoring.BlockNumber]{}, yConverter, OrderedTypeComparator[monitoring.Node]{})
}

// AddNodeTimeSeriesSource is a shortcut for AddSource(), which adds a source where the subject
// is a monitoring.Node and the series' key is monitoring.Time
func AddNodeTimeSeriesSource[T any](
	exporter *CsvExporter,
	source monitoring.Source[monitoring.Node, monitoring.TimeSeries[T]],
	yConverter converter[T],
) {
	AddSource[monitoring.Node, monitoring.Time, T, monitoring.TimeSeries[T]](exporter, source, MonitoringTimeConverter{}, yConverter, OrderedTypeComparator[monitoring.Node]{})
}

// AddNetworkBlockSeriesSource is a shortcut for AddSource(), which adds a source where the subject
// is a monitoring.Network and the series' key is monitoring.BlockNumber
func AddNetworkBlockSeriesSource[T any](
	exporter *CsvExporter,
	source monitoring.Source[monitoring.Network, monitoring.BlockSeries[T]],
	yConverter converter[T],
) {
	AddSource[monitoring.Network, monitoring.BlockNumber, T, monitoring.BlockSeries[T]](exporter, source, DirectConverter[monitoring.BlockNumber]{}, yConverter, NoopComparator[monitoring.Network]{})
}

// AddNetworkTimeSeriesSource is a shortcut for AddSource(), which adds a source where the subject
// is a monitoring.Network and the series' key is monitoring.Time
func AddNetworkTimeSeriesSource[T any](
	exporter *CsvExporter,
	source monitoring.Source[monitoring.Network, monitoring.TimeSeries[T]],
	yConverter converter[T],
) {
	AddSource[monitoring.Network, monitoring.Time, T, monitoring.TimeSeries[T]](exporter, source, MonitoringTimeConverter{}, yConverter, NoopComparator[monitoring.Network]{})
}

func (c *CsvExporter) append(line string) {
	if c.currentLine >= len(c.lines) {
		c.lines = append(c.lines, line)
	} else {
		c.lines[c.currentLine] += line
	}

	c.currentLine++
}

func (c *CsvExporter) flush() error {
	for _, line := range c.lines {
		if _, err := c.writer.Write([]byte(line + "\n")); err != nil {
			return err
		}
	}

	c.lines = c.lines[0:0]
	c.seriesSize = 0
	c.currentLine = 0
	c.firstColumn = true
	return nil
}

// Flush writes so far accumulated lines into the output writer as CSV lines.
func (c *CsvExporter) Flush() error {
	err := c.flush()
	return errors.Join(err, errors.Join(c.errors...))
}

func (c *CsvExporter) Close() error {
	return errors.Join(errors.Join(c.errors...), c.flush(), c.writer.Close())
}

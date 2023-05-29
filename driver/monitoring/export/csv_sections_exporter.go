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

// CsvExporter allows for accumulating data in a form of CSV lines and export them to the output writer.
// The CSV is composed of monitoring.Source instances, which can be organised into sections.
// It is expected that the sources contain monitoring.Series as its values. No other types
// are supported at the moment.
// Every section can include many sources and subjects.
// The series keys from the first series are printed in the first column representing the X-axis, and series
// values (i.e. the Y-axis) are printed in the next column. Further series values are printed in consecutive columns,
// The algorithm tries to match the X-axis of inserted series of one source. It checks if the consecutive
// series has the same or higher series' keys as the series already printed in the first column.
// It possibly moves or skips values in the series to match the first column.
// If no single element would be printed because no match is found,
// an extra column containing keys (X-axis) of this series pretending the values is print. In other words,
// if subjects of the source contain series with incompatible keys (X-axis) the X-axis is printed extra
// for each series.
// The matching algorithm is simple, it checks if the position to be inserted is higher or equal to current position.
// It assumes that the series positions are growing, i.e. finding next equal or higher position assures that
// the finding match fits into the window between current position and following position. It is useful for instance,
// when the X-axis represents a sampling period, then the keys do not have to exactly match, but they will
// fill in the same sampling window.
// Furthermore, if the series following the first series are longer or shorter than the first one,
// the series are is either trimmed or not printed fully.
// Finally, every section has a header with the title of the first column matching the X-axis type and subject + metric
// name being the header of next (Y-axis) columns.
type CsvExporter struct {
	writer          io.WriteCloser
	lines           []*strings.Builder
	firstColumn     []int
	isFirstColumn   bool
	firstColumnSize int
	errors          []error
}

// NewCsvExporter creates a new CSV exporter that writes the CSV lines into the input writer.
func NewCsvExporter(writer io.WriteCloser) *CsvExporter {
	return &CsvExporter{
		writer:        writer,
		lines:         make([]*strings.Builder, 0, 1000),
		firstColumn:   make([]int, 0, 1000),
		errors:        make([]error, 0, 5),
		isFirstColumn: true,
	}
}

// AddEmptySection inserts the input amount of empty lines
func AddEmptySection(exporter *CsvExporter, lines int) {
	if err := exporter.flush(); err != nil {
		exporter.errors = append(exporter.errors, err)
	}
	for i := 0; i < lines; i++ {
		exporter.newLine()
		exporter.concat(i, ",")
	}
}

// AddSection starts a new section, which will produce lines below the current section.
// This method has to be called before consecutive calls to AddSource().
func AddSection(exporter *CsvExporter) {
	if err := exporter.flush(); err != nil {
		exporter.errors = append(exporter.errors, err)
	}
}

// AddSource adds lines for the input source's series. It is expected that the source's values
// are monitoring.Series. No other type of value is supported.
// If this is the first source added in this section, the series keys (i.e. X-axis) are printed
// in the first column. The values are printed next to it in the following column
// (i.e. the Y-axis). If further sources are added, only the series' values are printed to next columns
// as long as series keys can be matched with the first column.
// If matching is not found, the keys are printed in next column prepending the series values in further column.
func AddSource[S any, K constraints.Integer, T any, X monitoring.Series[K, T]](
	exporter *CsvExporter,
	source monitoring.Source[S, X],
	kConverter converter[K],
	tConverter converter[T],
	cmp comparator[S],
) {
	subjects := source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return cmp.Compare(&subjects[i], &subjects[j]) < 0 })

	// get the sizes, i.e. the number of data points
	// store in order to arrays as subjects are sorted
	data := make([][]monitoring.DataPoint[K, T], 0, len(subjects))
	sizes := make([]int, 0, len(subjects))
	var longestSeriesIndex int
	var size int
	for i, subject := range subjects {
		series, exists := source.GetData(subject)
		dataPoints := copySeries[K, T](series)
		if !exists {
			panic(fmt.Errorf("series for subject: %v does not exist", subject))
		}
		data = append(data, dataPoints)
		currentSize := len(dataPoints)
		if currentSize > size {
			size = currentSize
			longestSeriesIndex = i
		}
		sizes = append(sizes, currentSize)
	}

	if exporter.isFirstColumn {
		exporter.firstColumnSize = size
		// two lines for header
		exporter.newLine()
		exporter.newLine()

		// memorise first column positions and reserve lines
		for i := 0; i < exporter.firstColumnSize; i++ {
			exporter.newLine()
			exporter.firstColumn = append(exporter.firstColumn, int(data[longestSeriesIndex][i].Position))
		}
	}

	// continue with data
	for i := 0; i < len(subjects); i++ {
		var seriesMatch bool
		var headerPrint bool
		for !seriesMatch {
			// it contains two indexes, 'k' and 'j'
			// 'j' represents current row, which is incremented every loop
			// 'k' indexes position in current series, which is incremented every loop when the position matches with position at index 'j'
			// when positions do not much, 'j' is incremented until a matching position is found at index 'k', i.e. the series is shifted when needed
			// if no match is found, this loop repeats while printing the positions (x-axis) again, i.e. next column will be bound to this new axis.
			var j, k int
			for j < exporter.firstColumnSize {

				// print position for first column
				if exporter.isFirstColumn {
					exporter.concat(j+2, fmt.Sprintf("%v, ", kConverter.Convert(K(exporter.firstColumn[j]))))
				}

				// try to move to next row when position does not match, try the series may be shifted
				// matching is simply checking if the position is greater than current position.
				if k < sizes[i] && int(data[i][k].Position) > exporter.firstColumn[j] {
					j++
					continue
				}

				if k >= sizes[i] {
					// series is shorter, replace with empty cells.
					exporter.concat(j+2, ", ")
				} else {
					point := data[i][k]
					exporter.concat(j+2, fmt.Sprintf("%s, ", tConverter.Convert(point.Value)))
					seriesMatch = true
				}

				// move next
				j++
				k++
			}

			// if no single row was matched, print X-axis again while printing this subject
			if !seriesMatch {
				exporter.isFirstColumn = true
				// memorise first column positions for current series
				exporter.firstColumn = exporter.firstColumn[0:0]
				for r := 0; r < exporter.firstColumnSize; r++ {
					exporter.firstColumn = append(exporter.firstColumn, int(data[i][r].Position))
				}
			}

			// print header
			if !headerPrint {
				if exporter.isFirstColumn {
					var firstColumn K
					exporter.concat(0, ", ")
					exporter.concat(1, fmt.Sprintf("%T, ", firstColumn))
				}
				if i == 0 {
					exporter.concat(0, fmt.Sprintf("%v, ", source.GetMetric().Name))
				} else {
					exporter.concat(0, ", ")
				}
				exporter.concat(1, fmt.Sprintf("%v, ", subjects[i]))
			}
			headerPrint = true
		}

		exporter.isFirstColumn = false
	}
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

func (c *CsvExporter) newLine() {
	c.lines = append(c.lines, &strings.Builder{})
}

func (c *CsvExporter) concat(lineNum int, line string) {
	c.lines[lineNum].WriteString(line)
}

func copySeries[K constraints.Ordered, T any](series monitoring.Series[K, T]) []monitoring.DataPoint[K, T] {
	lastPoint := series.GetLatest()
	if lastPoint == nil {
		return []monitoring.DataPoint[K, T]{}
	}

	var k K
	dataPoints := series.GetRange(k, lastPoint.Position)
	res := make([]monitoring.DataPoint[K, T], 0, len(dataPoints)+1)
	for _, point := range dataPoints {
		res = append(res, point)
	}

	return append(res, *lastPoint)
}

func (c *CsvExporter) flush() error {
	for _, line := range c.lines {
		if _, err := c.writer.Write([]byte(strings.TrimSpace(line.String()) + "\n")); err != nil {
			return err
		}
	}

	c.lines = c.lines[0:0]
	c.firstColumnSize = 0
	c.isFirstColumn = true
	c.firstColumn = c.firstColumn[0:0]
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

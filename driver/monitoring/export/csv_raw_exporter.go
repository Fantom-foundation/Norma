package export

import (
	"errors"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"io"
	"time"
)

// AddNodeBlockSeriesSource iterates all series of the input source and prints them into the input writer
// as a CSV line.
// This method is typed for sources where the subject
// is a monitoring.Node and the series' key is monitoring.BlockNumber
func AddNodeBlockSeriesSource[T any](
	exporter io.Writer,
	source monitoring.Source[monitoring.Node, monitoring.BlockSeries[T]],
	yConverter Converter[T],
) error {
	kConverter := DirectConverter[monitoring.BlockNumber]{}
	forEacher := monitoring.NewSourceRowsForEacher[monitoring.Node, monitoring.BlockNumber, T, monitoring.BlockSeries[T]](source)
	var errs []error
	forEacher.ForEachRow(func(row monitoring.Row[monitoring.Node, monitoring.BlockNumber, T, monitoring.BlockSeries[T]]) {
		line := fmt.Sprintf("%s, network, %v, , %v, %v", row.Metric.Name, row.Subject, kConverter.Convert(row.Position), yConverter.Convert(row.Value))
		if _, err := exporter.Write([]byte(line + "\n")); err != nil {
			errs = append(errs, err)
		}
	}, monitoring.OrderedTypeComparator[monitoring.Node]{})

	return errors.Join(errs...)
}

// AddNodeTimeSeriesSource iterates all series of the input source and prints them into the input writer
// as a CSV line.
// This method is typed for sources where the subject
// is a monitoring.Node and the series' key is monitoring.Time
func AddNodeTimeSeriesSource[T any](
	exporter io.Writer,
	source monitoring.Source[monitoring.Node, monitoring.TimeSeries[T]],
	yConverter Converter[T],
) error {
	kConverter := MonitoringTimeConverter{}
	forEacher := monitoring.NewSourceRowsForEacher[monitoring.Node, monitoring.Time, T, monitoring.TimeSeries[T]](source)
	var errs []error
	forEacher.ForEachRow(func(row monitoring.Row[monitoring.Node, monitoring.Time, T, monitoring.TimeSeries[T]]) {
		line := fmt.Sprintf("%s, network, %v, %v, , %v", row.Metric.Name, row.Subject, kConverter.Convert(row.Position), yConverter.Convert(row.Value))
		if _, err := exporter.Write([]byte(line + "\n")); err != nil {
			errs = append(errs, err)
		}
	}, monitoring.OrderedTypeComparator[monitoring.Node]{})

	return errors.Join(errs...)
}

// AddNetworkBlockSeriesSource iterates all series of the input source and prints them into the input writer
// as a CSV line.
// This method is typed for sources where the subject
// is a monitoring.Network and the series' key is monitoring.BlockNumber
func AddNetworkBlockSeriesSource[T any](
	exporter io.Writer,
	source monitoring.Source[monitoring.Network, monitoring.BlockSeries[T]],
	yConverter Converter[T],
) error {
	kConverter := DirectConverter[monitoring.BlockNumber]{}
	forEacher := monitoring.NewSourceRowsForEacher[monitoring.Network, monitoring.BlockNumber, T, monitoring.BlockSeries[T]](source)
	var errs []error
	forEacher.ForEachRow(func(row monitoring.Row[monitoring.Network, monitoring.BlockNumber, T, monitoring.BlockSeries[T]]) {
		line := fmt.Sprintf("%s, network, , , %v, %v", row.Metric.Name, kConverter.Convert(row.Position), yConverter.Convert(row.Value))
		if _, err := exporter.Write([]byte(line + "\n")); err != nil {
			errs = append(errs, err)
		}
	}, monitoring.NoopComparator[monitoring.Network]{})

	return errors.Join(errs...)
}

// AddNetworkTimeSeriesSource iterates all series of the input source and prints them into the input writer
// as a CSV line.
// This method is typed for sources where the subject
// is a monitoring.Network and the series' key is monitoring.Time
func AddNetworkTimeSeriesSource[T any](
	exporter io.Writer,
	source monitoring.Source[monitoring.Network, monitoring.TimeSeries[T]],
	yConverter Converter[T],
) error {
	kConverter := MonitoringTimeConverter{}
	forEacher := monitoring.NewSourceRowsForEacher[monitoring.Network, monitoring.Time, T, monitoring.TimeSeries[T]](source)
	var errs []error
	forEacher.ForEachRow(func(row monitoring.Row[monitoring.Network, monitoring.Time, T, monitoring.TimeSeries[T]]) {
		line := fmt.Sprintf("%s, network, , %v, , %v", row.Metric.Name, kConverter.Convert(row.Position), yConverter.Convert(row.Value))
		if _, err := exporter.Write([]byte(line + "\n")); err != nil {
			errs = append(errs, err)
		}
	}, monitoring.NoopComparator[monitoring.Network]{})

	fmt.Printf("%v", time.Unix(0, 80000000))

	return errors.Join(errs...)
}

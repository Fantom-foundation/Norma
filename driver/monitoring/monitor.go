package monitoring

import (
	"errors"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"os"
)

// Monitor instances are handling the life-cycle of sets of data sources for a
// a configurable set of metrics. Instances are to be created using the
// NewMonitor() factory below and required to be shut down in the end.
//
// Monitor instances provide means for registering metric sources and for
// obtaining data for respective metrics. The implementation aims at keeping
// metric access type save. However, it is not possible to define generic
// methods in Go. Thus, several methods interacting with Monitor instances
// are free functions (see implementations below).
type Monitor struct {
	Network         driver.Network
	NodeLogProvider NodeLogProvider
	sources         map[string]source
	Writer          WriterChain
}

// NewMonitor creates a new Monitor instance without any registered sources.
func NewMonitor(network driver.Network, csvFile *os.File) *Monitor {
	return &Monitor{network, NewNodeLogDispatcher(network), map[string]source{}, NewWriterChain(csvFile)}
}

// Shutdown disconnects all sources, stopping the collection of data. This
// should be called before abandoning the Monitor instance.
func (m *Monitor) Shutdown() error {
	var errs = []error{}
	for _, source := range m.sources {
		if err := source.Shutdown(); err != nil {
			errs = append(errs, err)
		}
	}

	errs = append(errs, m.Writer.Close())
	return errors.Join(errs...)
}

// InstallSource installs a new source on the given monitor. The provided factory
// is used to create a new source instance, of which the monitor takes ownership.
// In particular, the monitor will stop it during the Shutdown of the monitor.
func InstallSource[S any, T any](monitor *Monitor, factory SourceFactory[S, T]) error {
	metric := factory.GetMetric()
	_, present := monitor.sources[metric.Name]
	if present {
		return fmt.Errorf("source for metric %s already present", metric.Name)
	}
	monitor.sources[metric.Name] = factory.CreateSource(monitor)
	return nil
}

// IsSupported checks whether there is a source registered for the given metric.
func IsSupported[S any, T any](monitor *Monitor, metric Metric[S, T]) bool {
	_, present := monitor.sources[metric.Name]
	return present
}

// GetSubjects retrieves all subjects with available data for the given metric.
func GetSubjects[S any, T any](monitor *Monitor, metric Metric[S, T]) []S {
	source := monitor.sources[metric.Name]
	if source == nil {
		return nil
	}
	return source.(Source[S, T]).GetSubjects()
}

// GetData retrieves access to the data collected for a given metric or nil, if
// the defined metric for the given subject is not available.
func GetData[S any, T any](monitor *Monitor, subject S, metric Metric[S, T]) (t T, exists bool) {
	source := monitor.sources[metric.Name]
	if source == nil {
		return t, false
	}
	return source.(Source[S, T]).GetData(subject)
}

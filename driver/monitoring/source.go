package monitoring

// Source is a provider of monitoring data for a single metric. A source is
// an active component collecting data asynchroniously and retaining it for
// the duration of its life cycle.
type Source[S any, T any] interface {
	source
	// GetMetric returns the provided metric.
	GetMetric() Metric[S, T]

	// GetSubjects returns a list of subjects for which monitoring data is
	// available in this source. The list is expected to grow monotonies.
	GetSubjects() []S

	// GetData obtains the monitoring data retained for a selected subject
	// or nil if there is no such data.
	GetData(S) *T
}

// source is a type-erased base type for sources.
type source interface {
	// Shutdown stops the collection of data. Already collected data shall
	// remain available, but no new data is collected.
	Shutdown() error
}

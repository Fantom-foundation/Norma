package monitoring

import (
	"golang.org/x/exp/constraints"
)

// SourceRowsForEacher wraps a Source, of which values are Series. It implements
// a for-each iteration of the data stored in this Source in a form
// of Rows.
type SourceRowsForEacher[S any, K constraints.Ordered, T any, X Series[K, T]] struct {
	source Source[S, X]
}

// NewSourceRowsForEacher creates a new instance, which wraps a Source. The Source has to be bound to a Series.
func NewSourceRowsForEacher[S any, K constraints.Ordered, T any, X Series[K, T]](source Source[S, X]) *SourceRowsForEacher[S, K, T, X] {
	return &SourceRowsForEacher[S, K, T, X]{source}
}

// Row is a struct containing fields of one row exported from the Source.
// It is used for export of data from the Source to print it in the console, CSV file, etc.
type Row[S any, K constraints.Ordered, T any, X Series[K, T]] struct {
	Metric   Metric[S, X]
	Subject  S
	Position K
	Value    T
}

// ForEachRow iterates all rows representing data stored in the wrapped source.
func (s *SourceRowsForEacher[S, K, T, X]) ForEachRow(f func(Row[S, K, T, X])) {
	metrics := s.source.GetMetric()
	subjects := s.source.GetSubjects()

	for _, subject := range subjects {
		series, _ := s.source.GetData(subject)
		last := series.GetLatest()
		if last != nil {
			var first K
			for _, point := range series.GetRange(first, last.Position) {
				row := Row[S, K, T, X]{
					Metric:   metrics,
					Subject:  subject,
					Position: point.Position,
					Value:    point.Value,
				}
				f(row)
			}
			// include last element
			row := Row[S, K, T, X]{
				Metric:   metrics,
				Subject:  subject,
				Position: last.Position,
				Value:    last.Value,
			}
			f(row)
		}
	}
}

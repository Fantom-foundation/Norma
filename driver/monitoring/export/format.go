package export

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"golang.org/x/exp/constraints"
	"time"
)

// comparator is used for comparing subjects added as sources.
type comparator[T any] interface {
	Compare(a, b *T) int
}

// OrderedTypeComparator compares constraints.Ordered by using <, > operators.
type OrderedTypeComparator[T constraints.Ordered] struct{}

func (OrderedTypeComparator[T]) Compare(a, b *T) int {
	if *a > *b {
		return 1
	}
	if *a < *b {
		return -1
	}

	return 0
}

// NoopComparator performs no comparison, it is used for types that cannot be ordered.
type NoopComparator[T any] struct{}

func (NoopComparator[T]) Compare(_, _ *T) int {
	return 0
}

// converter is used for converting input types into strings that is suitable to CSV format.
type converter[T any] interface {
	Convert(t T) string
}

// TimeConverter converts the time.Time timestamp to print only time including milliseconds.
type TimeConverter struct{}

func (TimeConverter) Convert(time time.Time) string {
	return time.Format("15:04:05.999")
}

// MonitoringTimeConverter converts the monitoring.Time the same as TimeConverter.
type MonitoringTimeConverter struct{}

func (MonitoringTimeConverter) Convert(time monitoring.Time) string {
	return TimeConverter{}.Convert(time.Time().UTC())
}

// DurationConverter converts time.Duration to print Nanoseconds.
type DurationConverter struct{}

func (DurationConverter) Convert(time time.Duration) string {
	return fmt.Sprintf("%d", time.Nanoseconds())
}

// DirectConverter uses default string conversion to any input types.
type DirectConverter[T any] struct{}

func (DirectConverter[T]) Convert(t T) string {
	return fmt.Sprintf("%v", t)
}

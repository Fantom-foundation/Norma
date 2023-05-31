package export

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"time"
)

// Converter is used for converting input types into strings that is suitable to CSV format.
type Converter[T any] interface {
	Convert(t T) string
}

// TimeConverter converts the time.Time timestamp to print only time including milliseconds.
type TimeConverter struct{}

func (TimeConverter) Convert(time time.Time) string {
	return fmt.Sprintf("%d", uint64(time.UnixNano()))
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

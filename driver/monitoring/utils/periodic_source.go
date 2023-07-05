package utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
)

// Sensor is an abstraction of some input device capable of probing a node
// for some metric of type T.
type Sensor[T any] interface {
	ReadValue() (T, error)
}

// PeriodicDataSource is a generic data source periodically querying
// node-associated sensors for data.
type PeriodicDataSource[S comparable, T any] struct {
	*SyncedSeriesSource[S, monitoring.Time, T]
	period time.Duration
	stop   chan bool  // used to signal per-node collectors about the shutdown
	done   chan error // used to signal collector shutdown to source
}

// NewPeriodicDataSource creates a new data source managing per-node sensor
// instances for a given metric and periodically collecting data from those.
func NewPeriodicDataSource[S comparable, T any](
	metric monitoring.Metric[S, monitoring.Series[monitoring.Time, T]],
	monitor *monitoring.Monitor,
) *PeriodicDataSource[S, T] {
	return NewPeriodicDataSourceWithPeriod(metric, monitor, time.Second)
}

// NewPeriodicDataSourceWithPeriod is the same as NewPeriodicDataSource but with
// a customizable sampling periode.
func NewPeriodicDataSourceWithPeriod[S comparable, T any](
	metric monitoring.Metric[S, monitoring.Series[monitoring.Time, T]],
	monitor *monitoring.Monitor,
	period time.Duration,
) *PeriodicDataSource[S, T] {
	stop := make(chan bool)
	done := make(chan error)

	res := &PeriodicDataSource[S, T]{
		SyncedSeriesSource: NewSyncedSeriesSource(metric),
		period:             period,
		stop:               stop,
		done:               done,
	}

	return res
}

func (s *PeriodicDataSource[S, T]) Shutdown() error {
	close(s.stop)
	<-s.done
	return s.SyncedSeriesSource.Shutdown()
}

func (s *PeriodicDataSource[S, T]) AddSubject(subject S, sensor Sensor[T]) error {
	data, err := s.NewSubject(subject)
	if err != nil {
		return err
	}

	// Start background routine collecting sensor data.
	go func() {
		var err error
		defer func() {
			s.done <- err
		}()

		// Introduce random sampling offsets to avoid load peaks and to
		// eliminate steps in aggregated metrics.
		time.Sleep(time.Duration(float32(s.period) * rand.Float32()))

		var errs []error
		ticker := time.NewTicker(s.period)
		defer ticker.Stop()
		for {
			select {
			case now := <-ticker.C:
				value, err := sensor.ReadValue()
				if err != nil {
					errs = append(errs, err)
				} else {
					data.Append(monitoring.NewTime(now), value)
				}
			case <-s.stop:
				err = errors.Join(errs...)
				return
			}
		}
	}()

	return nil
}

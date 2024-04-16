// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

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
	period   time.Duration
	subjects map[S]process // used to signal removal of a subject
}

// process represents a structure that encapsulates a stop signal and potential done.
// It is used to communicate the stop signal and any done that occur during the execution of a subject.
// The stop signal is sent through the 'stop' channel of type 'chan bool'.
// Any done that occur are stored in the 'done' field of type 'error'.
// This structure is typically used in the context of controlling the execution of a subject.
type process struct {
	stop chan bool
	done chan error
}

// Stop stops the instance by closing the stop channel and
// returning the error received from the done channel.
func (s process) Stop() error {
	close(s.stop)
	return <-s.done
}

// Done sends the provided error on the done channel and closes it.
func (s process) Done(err error) {
	s.done <- err
	close(s.done)
}

// NotifyStop returns the stop channel of the process instance.
// This channel notifies that the process should stop.
func (s process) NotifyStop() chan bool {
	return s.stop
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
	res := &PeriodicDataSource[S, T]{
		SyncedSeriesSource: NewSyncedSeriesSource(metric),
		period:             period,
		subjects:           make(map[S]process),
	}

	return res
}

func (s *PeriodicDataSource[S, T]) Shutdown() error {
	// stop all subjects and drain potential done
	var err error
	for _, stop := range s.subjects {
		err = errors.Join(err, stop.Stop())
	}

	return errors.Join(err, s.SyncedSeriesSource.Shutdown())
}

func (s *PeriodicDataSource[S, T]) AddSubject(subject S, sensor Sensor[T]) error {
	data, err := s.NewSubject(subject)
	if err != nil {
		return err
	}

	subjectStop := process{make(chan bool), make(chan error, 1)}
	s.subjects[subject] = subjectStop

	// Start background routine collecting sensor data.
	go func() {
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
					if err := data.Append(monitoring.NewTime(now), value); err != nil {
						errs = append(errs, err)
					}
				}
			case <-subjectStop.NotifyStop():
				subjectStop.Done(errors.Join(errs...))
				return
			}
		}
	}()

	return nil
}

func (s *PeriodicDataSource[S, T]) RemoveSubject(subject S) error {
	subjectStop, exists := s.subjects[subject]
	if exists {
		delete(s.subjects, subject)
		return subjectStop.Stop()
	}

	return nil
}

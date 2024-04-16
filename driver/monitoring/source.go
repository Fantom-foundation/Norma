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

package monitoring

import (
	"errors"
	"fmt"
)

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

	// GetData obtains the monitoring data retained for a selected subject.
	// The other argument returns false if the value for the subject does not exist, true otherwise.
	// When true is returned, the first return parameter is not nil, otherwise it is unspecified.
	GetData(S) (T, bool)
}

// source is a type-erased base type for sources. While its methods should
// be public, the interface itself is only intended to be used internally to
// store multiple sources of different generic types in a common map.
type source interface {
	// Shutdown stops the collection of data. Already collected data shall
	// remain available, but no new data is collected.
	Shutdown() error

	// ForEachRecord requires each source to enumerate all collected records.
	// This is mainly intended to provide a uniform export interface.
	ForEachRecord(func(r Record))
}

// SourceFactory is a generic interface for metric sources. It is used to
// register metrics and their sources in Norma's monitoring system.
type SourceFactory[S any, T any] interface {
	// GetMetric returns the metric the source is providing.
	GetMetric() Metric[S, T]
	// CreateSource creates a new source instance collecting data within
	// the given monitoring.
	CreateSource(monitor *Monitor) Source[S, T]
}

// RegisterFactory registers a new source factory in a global registry. It is
// intended to be called in initialization code to announce the availability
// of metric sources.
func RegisterFactory[S any, T any](factory SourceFactory[S, T]) error {
	metric := factory.GetMetric()
	_, present := sourceInstallers[metric.Name]
	if present {
		return fmt.Errorf("metric collision: multiple sources for metric '%s' encountered", metric.Name)
	}
	sourceInstallers[metric.Name] = &sourceAdapter[S, T]{factory}
	return nil
}

// RegisterSource is a convenience variant of RegisterFactory above, accepting
// a metric and a factory function for registering a source.
func RegisterSource[S any, T any](metric Metric[S, T], factory func(*Monitor) Source[S, T]) error {
	return RegisterFactory[S, T](&genericSourceFactory[S, T]{metric, factory})
}

// InstallAllRegisteredSources installs one instance of every registered source
// in the given monitor. The resulting error represents the union of all
// errors that occured during the source creation and installation.
func InstallAllRegisteredSources(monitor *Monitor) error {
	var errs []error
	for _, installer := range sourceInstallers {
		if err := installer.installIn(monitor); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// InstallSourceFor installs a data source for the given metric in the given monitor.
func InstallSourceFor[S any, T any](metric Metric[S, T], monitor *Monitor) error {
	installer, exists := sourceInstallers[metric.Name]
	if !exists {
		return fmt.Errorf("no definition registered for metric %s", metric.Name)
	}
	return installer.installIn(monitor)
}

// sourceInstallers is the internal global registry of metric sources.
var sourceInstallers = map[string]sourceInstaller{}

// sourceInstaller is an internal interface for type-erased metric source
// installers.
type sourceInstaller interface {
	installIn(monitor *Monitor) error
}

// sourceAdapter is a type-safe adapter bridging the gap between source
// factories and source installers.
type sourceAdapter[S any, T any] struct {
	factory SourceFactory[S, T]
}

func (i *sourceAdapter[S, T]) installIn(monitor *Monitor) error {
	return InstallSource(monitor, i.factory)
}

type genericSourceFactory[S any, T any] struct {
	metric  Metric[S, T]
	factory func(*Monitor) Source[S, T]
}

func (f *genericSourceFactory[S, T]) GetMetric() Metric[S, T] {
	return f.metric
}

func (f *genericSourceFactory[S, T]) CreateSource(monitor *Monitor) Source[S, T] {
	return f.factory(monitor)
}

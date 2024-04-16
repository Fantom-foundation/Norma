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

package appmon

import (
	"log"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
)

// SensorFactory is a factory for sensors targeting selected applications.
type SensorFactory[T any] interface {
	CreateSensor(driver.Application) (utils.Sensor[T], error)
}

// periodicAppDataSource is a generic data source periodically querying
// node-associated sensors for data.
type periodicAppDataSource[T any] struct {
	*utils.PeriodicDataSource[mon.App, T]
	factory SensorFactory[T]
}

// NewPeriodicAppDataSource creates a new data source managing per-app sensor
// instances for a given metric and periodically collecting data from those.
func NewPeriodicAppDataSource[T any](
	metric mon.Metric[mon.App, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	factory SensorFactory[T],
) mon.Source[mon.App, mon.Series[mon.Time, T]] {
	return newPeriodicAppDataSource(metric, monitor, time.Second, factory)
}

// newPeriodicAppDataSource is the same as NewPeriodicAppDataSource but with
// a customizable sampling periode.
func newPeriodicAppDataSource[T any](
	metric mon.Metric[mon.App, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	period time.Duration,
	factory SensorFactory[T],
) mon.Source[mon.App, mon.Series[mon.Time, T]] {
	res := &periodicAppDataSource[T]{
		PeriodicDataSource: utils.NewPeriodicDataSourceWithPeriod(metric, monitor, period),
		factory:            factory,
	}

	monitor.Network().RegisterListener(res)
	for _, app := range monitor.Network().GetActiveApplications() {
		res.AfterApplicationCreation(app)
	}

	return res
}

func (s *periodicAppDataSource[T]) AfterNodeCreation(driver.Node) {
	// ignored
}

func (s *periodicAppDataSource[T]) AfterNodeRemoval(driver.Node) {
	// ignored
}

func (s *periodicAppDataSource[T]) AfterApplicationCreation(app driver.Application) {
	label := app.Config().Name
	sensor, err := s.factory.CreateSensor(app)
	if err != nil {
		log.Printf("failed to create sensor for metric %v / app %s: %v", s.GetMetric().Name, label, err)
		return
	}
	s.AddSubject(mon.App(label), sensor)
}

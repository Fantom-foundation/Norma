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

package user

import (
	"log"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
)

// SensorFactory is a factory for sensors targeting selected users.
type SensorFactory[T any] interface {
	CreateSensor(driver.Application, int) (utils.Sensor[T], error)
}

// periodicUserDataSource is a generic data source periodically querying
// user-associated sensors for data.
type periodicUserDataSource[T any] struct {
	*utils.PeriodicDataSource[mon.User, T]
	factory SensorFactory[T]
}

// NewPeriodicUserDataSource creates a new data source managing per-user
// sensor instances for a given metric and periodically collecting data from those.
func NewPeriodicUserDataSource[T any](
	metric mon.Metric[mon.User, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	factory SensorFactory[T],
) mon.Source[mon.User, mon.Series[mon.Time, T]] {
	return newPeriodicUserDataSource(metric, monitor, time.Second, factory)
}

// newPeriodicUserDataSource is the same as NewPeriodicUserDataSource but with
// a customizable sampling periode.
func newPeriodicUserDataSource[T any](
	metric mon.Metric[mon.User, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	period time.Duration,
	factory SensorFactory[T],
) mon.Source[mon.User, mon.Series[mon.Time, T]] {
	res := &periodicUserDataSource[T]{
		PeriodicDataSource: utils.NewPeriodicDataSourceWithPeriod(metric, monitor, period),
		factory:            factory,
	}

	monitor.Network().RegisterListener(res)
	for _, app := range monitor.Network().GetActiveApplications() {
		res.AfterApplicationCreation(app)
	}

	return res
}

func (s *periodicUserDataSource[T]) AfterNodeCreation(driver.Node) {
	// ignored
}

func (s *periodicUserDataSource[T]) AfterNodeRemoval(driver.Node) {
	// ignored
}

func (s *periodicUserDataSource[T]) AfterApplicationCreation(app driver.Application) {
	label := mon.App(app.Config().Name)
	for i := 0; i < app.Config().Users; i++ {
		sensor, err := s.factory.CreateSensor(app, i)
		if err != nil {
			log.Printf("failed to create sensor for metric %v / app %s / user %d: %v", s.GetMetric().Name, label, i, err)
			return
		}
		s.AddSubject(mon.User{
			App: label,
			Id:  i,
		}, sensor)
	}
}

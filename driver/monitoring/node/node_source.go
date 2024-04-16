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

package nodemon

import (
	"log"
	"time"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
)

// SensorFactory is a factory for sensors targeting selected nodes.
type SensorFactory[T any] interface {
	CreateSensor(driver.Node) (utils.Sensor[T], error)
}

// periodicNodeDataSource is a generic data source periodically querying
// node-associated sensors for data.
type periodicNodeDataSource[T any] struct {
	*utils.PeriodicDataSource[monitoring.Node, T]
	factory SensorFactory[T]
}

// NewPeriodicNodeDataSource creates a new data source managing per-node sensor
// instances for a given metric and periodically collecting data from those.
func NewPeriodicNodeDataSource[T any](
	metric mon.Metric[mon.Node, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	factory SensorFactory[T],
) mon.Source[mon.Node, mon.Series[mon.Time, T]] {
	return newPeriodicNodeDataSource(metric, monitor, time.Second, factory)
}

// newPeriodicNodeDataSource is the same as NewPeriodicNodeDataSource but with
// a customizable sampling periode.
func newPeriodicNodeDataSource[T any](
	metric mon.Metric[mon.Node, mon.Series[mon.Time, T]],
	monitor *mon.Monitor,
	period time.Duration,
	factory SensorFactory[T],
) mon.Source[mon.Node, mon.Series[mon.Time, T]] {
	res := &periodicNodeDataSource[T]{
		PeriodicDataSource: utils.NewPeriodicDataSourceWithPeriod(metric, monitor, period),
		factory:            factory,
	}

	monitor.Network().RegisterListener(res)
	for _, node := range monitor.Network().GetActiveNodes() {
		res.AfterNodeCreation(node)
	}

	return res
}

func (s *periodicNodeDataSource[T]) AfterNodeCreation(node driver.Node) {
	label := node.GetLabel()
	sensor, err := s.factory.CreateSensor(node)
	if err != nil {
		log.Printf("failed to create sensor for metric %v / node %s: %v", s.GetMetric().Name, label, err)
	}
	s.AddSubject(mon.Node(label), sensor)
}

func (s *periodicNodeDataSource[T]) AfterNodeRemoval(node driver.Node) {
	label := node.GetLabel()
	s.RemoveSubject(mon.Node(label))
}

func (s *periodicNodeDataSource[T]) AfterApplicationCreation(driver.Application) {
	// ignored
}

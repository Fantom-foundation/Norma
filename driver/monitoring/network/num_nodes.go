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

package netmon

import (
	"fmt"
	"time"

	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
)

// NumberOfNodes retains a time-series for the number of nodes in the network
// run by Norma. This includes all types of nodes.
var NumberOfNodes = mon.Metric[mon.Network, mon.Series[mon.Time, int]]{
	Name:        "NumberOfNodes",
	Description: "The number of connected nodes at various times.",
}

func init() {
	if err := mon.RegisterSource(NumberOfNodes, NewNumNodesSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// numNodesSource is a monitoring data source tracking the number of active
// nodes in a network environment.
type numNodesSource struct {
	*utils.SyncedSeriesSource[mon.Network, mon.Time, int]
	data *mon.SyncedSeries[mon.Time, int]
	stop chan<- bool
	done <-chan bool
}

// NewNumNodesSource creates a new data source periodically collecting data on
// the number of nodes in the network.
func NewNumNodesSource(monitor *mon.Monitor) mon.Source[mon.Network, mon.Series[mon.Time, int]] {
	return newNumNodesSource(monitor, time.Second)
}

// newNumNodesSource creates a new data source periodically collecting data on
// the number of nodes in the network.
func newNumNodesSource(monitor *mon.Monitor, period time.Duration) mon.Source[mon.Network, mon.Series[mon.Time, int]] {
	stop := make(chan bool)
	done := make(chan bool)

	res := &numNodesSource{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(NumberOfNodes),
		stop:               stop,
		done:               done,
	}
	res.data = res.GetOrAddSubject(mon.Network{})

	go func() {
		defer close(done)
		ticker := time.NewTicker(period)
		for {
			select {
			case now := <-ticker.C:
				numNodes := len(monitor.Network().GetActiveNodes())
				res.data.Append(mon.NewTime(now), numNodes)
			case <-stop:
				return
			}
		}
	}()

	return res
}

func (s *numNodesSource) Shutdown() error {
	close(s.stop)
	<-s.done
	return s.SyncedSeriesSource.Shutdown()
}

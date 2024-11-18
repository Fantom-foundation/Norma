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
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/Fantom-foundation/go-opera/opera/contracts/sfc"
)

// EpochProgress retains a time-series for the epoch number relative to sfc.
var EpochProgress = mon.Metric[mon.Network, mon.Series[mon.Time, int]]{
	Name:        "EpochProgress",
	Description: "The epoch number wrt SFC as a time-series.",
}

func init() {
	if err := mon.RegisterSource(EpochProgress, NewEpochProgressSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// epochProgressSource is a monitoring data source tracking the epoch number.
type epochProgressSource struct {
	*utils.SyncedSeriesSource[mon.Network, mon.Time, int]
	data *mon.SyncedSeries[mon.Time, int]
	stop chan<- bool
	done <-chan bool
}

// NewEpochProgressSource creates a new data source periodically collecting data on
// the current epoch number.
func NewEpochProgressSource(monitor *mon.Monitor) mon.Source[mon.Network, mon.Series[mon.Time, int]] {
	return newEpochProgressSource(monitor, time.Second)
}

// newEpochProgressSource creates a new data source periodically collecting
// data on the the list of ValidatorId participated in previous epoch sealing.
func newEpochProgressSource(monitor *mon.Monitor, period time.Duration) mon.Source[mon.Network, mon.Series[mon.Time, int]] {
	stop := make(chan bool)
	done := make(chan bool)

	res := &epochProgressSource{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(EpochProgress),
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
				// if not here, may fail when node ends
				rpcClient, err := monitor.Network().DialRandomRpc()
				if err != nil {
					return
				}

				sfcc, err := contract.NewSFC(sfc.ContractAddress, rpcClient)
				if err != nil {
					return
				}

				currentEpoch, err := sfcc.CurrentEpoch(nil)
				if err != nil {
					return
				}

				res.data.Append(mon.NewTime(now), int(currentEpoch.Int64()))
			case <-stop:
				return
			}
		}
	}()

	return res
}

func (s *epochProgressSource) Shutdown() error {
	close(s.stop)
	<-s.done
	return s.SyncedSeriesSource.Shutdown()
}

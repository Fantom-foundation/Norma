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
	"math/big"
	"time"

	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/Fantom-foundation/go-opera/opera/contracts/sfc"
)

// PreviousEpochValidators retains a time-series for the list of validators in previous
// epoch in the network run by Norma.
var PreviousEpochValidators = mon.Metric[mon.Network, mon.Series[mon.Time, []int]]{
	Name:        "PreviousEpochValidators",
	Description: "The list of ValidatorId participated in previous epoch sealing.",
}

func init() {
	if err := mon.RegisterSource(PreviousEpochValidators, NewPreviousEpochValidatorsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// previousEpochValidatorsSource is a monitoring data source tracking
// the list of ValidatorId participated in previous epoch sealing.
type previousEpochValidatorsSource struct {
	*utils.SyncedSeriesSource[mon.Network, mon.Time, []int]
	data *mon.SyncedSeries[mon.Time, []int]
	stop chan<- bool
	done <-chan bool
}

// NewNumNodesSource creates a new data source periodically collecting data on
// the number of nodes in the network.
func NewPreviousEpochValidatorsSource(monitor *mon.Monitor) mon.Source[mon.Network, mon.Series[mon.Time, []int]] {
	return newPreviousEpochValidatorsSource(monitor, time.Second)
}

// newPreviousEpochValidatorsSource creates a new data source periodically collecting
// data on the the list of ValidatorId participated in previous epoch sealing.
func newPreviousEpochValidatorsSource(monitor *mon.Monitor, period time.Duration) mon.Source[mon.Network, mon.Series[mon.Time, []int]] {
	stop := make(chan bool)
	done := make(chan bool)

	res := &previousEpochValidatorsSource{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(PreviousEpochValidators),
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

				currentEpoch, err := sfcc.CurretEpoch(nil)
				if err != nil {
					return
				}

				bIds, err := sfcc.GetEpochValidatorIDs(nil, currentEpoch)
				if err != nil {
					return
				}
				ids := toInts(bIds)
				fmt.Println(ids)

				res.data.Append(mon.NewTime(now), ids)
			case <-stop:
				return
			}
		}
	}()

	return res
}

func toInts(bigs []*big.Int) []int {
	res := make([]int, len(bigs))
	for i, b := range bigs {
		res[i] = int(b.Int64())
	}
	return res
}

func (s *previousEpochValidatorsSource) Shutdown() error {
	close(s.stop)
	<-s.done
	return s.SyncedSeriesSource.Shutdown()
}

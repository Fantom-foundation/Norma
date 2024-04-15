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
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
	opera "github.com/Fantom-foundation/Norma/driver/node"
	"github.com/ethereum/go-ethereum/rpc"
)

// NodeBlockHeight collects a per-node time series of its current block height.
var NodeBlockHeight = mon.Metric[mon.Node, mon.Series[mon.Time, int]]{
	Name:        "NodeBlockHeight",
	Description: "The block height of nodes at various times.",
}

func init() {
	if err := mon.RegisterSource(NodeBlockHeight, NewNodeBlockHeightSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// NewNodeBlockHeightSource creates a new data source periodically collecting data on
// the block height at various nodes over time.
func NewNodeBlockHeightSource(monitor *mon.Monitor) mon.Source[mon.Node, mon.Series[mon.Time, int]] {
	return newNodeBlockHeightSource(monitor, time.Second)
}

func newNodeBlockHeightSource(monitor *mon.Monitor, period time.Duration) mon.Source[mon.Node, mon.Series[mon.Time, int]] {
	return newPeriodicNodeDataSource[int](NodeBlockHeight, monitor, period, &blockProgressSensorFactory{})
}

type blockProgressSensorFactory struct{}

func (f *blockProgressSensorFactory) CreateSensor(node driver.Node) (utils.Sensor[int], error) {
	url := node.GetServiceUrl(&opera.OperaRpcService)
	if url == nil {
		return nil, fmt.Errorf("node does not export an RPC server")
	}
	rpcClient, err := rpc.DialContext(context.Background(), string(*url))
	if err != nil {
		return nil, err
	}
	return &blockProgressSensor{rpcClient}, nil
}

type blockProgressSensor struct {
	rpcClient *rpc.Client
}

func (s *blockProgressSensor) ReadValue() (int, error) {
	var blockNumber string
	err := s.rpcClient.Call(&blockNumber, "eth_blockNumber")
	if err != nil {
		return 0, err
	}
	blockNumber = strings.TrimPrefix(blockNumber, "0x")
	value, err := strconv.ParseInt(blockNumber, 16, 32)
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

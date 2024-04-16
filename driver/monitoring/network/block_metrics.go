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
	"log"

	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
)

var (
	// BlockNumberOfTransactions is a metric capturing number of transactions for each block of a node.
	BlockNumberOfTransactions = monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]]{
		Name:        "BlockNumberOfTransactions",
		Description: "The number of transactions processed in a block",
	}

	// BlockGasUsed is a metric capturing Gas used for each block of a node.
	BlockGasUsed = monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]]{
		Name:        "BlockGasUsed",
		Description: "The gas used in a block",
	}
)

func init() {
	if err := monitoring.RegisterSource(BlockNumberOfTransactions, newNumberOfTransactionsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}

	if err := monitoring.RegisterSource(BlockGasUsed, newGasUsedSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// BlockNetworkMetricSource is a metric source that captures block properties where the Metric is the subject
type BlockNetworkMetricSource[T any] struct {
	*utils.SyncedSeriesSource[monitoring.Network, monitoring.BlockNumber, T]
	series           *monitoring.SyncedSeries[monitoring.BlockNumber, T]
	getBlockProperty func(b monitoring.Block) T
	monitor          *monitoring.Monitor
	lastBlock        int // track last block added in the series not to add duplicated block heights
}

// NewNumberOfTransactionsSource creates a metric capturing number of transactions for each block of a network
func NewNumberOfTransactionsSource(monitor *monitoring.Monitor) *BlockNetworkMetricSource[int] {
	f := func(b monitoring.Block) int {
		return b.Txs
	}
	return newBlockNetworkMetricsSource[int](monitor, f, BlockNumberOfTransactions)
}

// NewGasUsedSource creates a metric capturing Gas used for each block of a network.
func NewGasUsedSource(monitor *monitoring.Monitor) *BlockNetworkMetricSource[int] {
	f := func(b monitoring.Block) int {
		return b.GasUsed
	}
	return newBlockNetworkMetricsSource[int](monitor, f, BlockGasUsed)
}

// newNumberOfTransactionsSource is the same as its public counterpart, it only returns the Source interface instead of the struct to be used in factories
func newNumberOfTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]] {
	return NewNumberOfTransactionsSource(monitor)
}

// newGasUsedSource is the same as its public counterpart, it only returns the Source interface instead of the struct to be used in factories
func newGasUsedSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Network, monitoring.Series[monitoring.BlockNumber, int]] {
	return NewGasUsedSource(monitor)
}

// newBlockNodeMetricsSource creates a new data source periodically collecting data from the Node log
func newBlockNetworkMetricsSource[T any](
	monitor *monitoring.Monitor,
	getBlockProperty func(b monitoring.Block) T,
	metric monitoring.Metric[monitoring.Network, monitoring.Series[monitoring.BlockNumber, T]]) *BlockNetworkMetricSource[T] {

	m := &BlockNetworkMetricSource[T]{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(metric),
		getBlockProperty:   getBlockProperty,
		monitor:            monitor,
		lastBlock:          -1,
	}
	m.series = m.GetOrAddSubject(monitoring.Network{})
	monitor.NodeLogProvider().RegisterLogListener(m)
	return m
}

func (s *BlockNetworkMetricSource[T]) Shutdown() error {
	s.monitor.NodeLogProvider().UnregisterLogListener(s)
	return s.SyncedSeriesSource.Shutdown()
}

func (s *BlockNetworkMetricSource[T]) OnBlock(_ monitoring.Node, block monitoring.Block) {
	if block.Height > s.lastBlock {
		if err := s.series.Append(monitoring.BlockNumber(block.Height), s.getBlockProperty(block)); err != nil {
			log.Printf("error to add to the series: %s", err)
		}
		s.lastBlock = block.Height
	}
}

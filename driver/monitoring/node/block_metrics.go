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
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
	"github.com/ethereum/go-ethereum/log"
)

var (
	// BlockCompletionTime is a metric capturing time of the block finalisation.
	BlockCompletionTime = monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.BlockNumber, time.Time]]{
		Name:        "BlockCompletionTime",
		Description: "Time the block was completed",
	}

	BlockEventAndTxsProcessingTime = monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.BlockNumber, time.Duration]]{
		Name:        "BlockEventAndTxsProcessingTime",
		Description: "Time to process a block, it applies all lachesis events, applies all transactions, and commits stateDB",
	}
)

func init() {
	if err := monitoring.RegisterSource(BlockCompletionTime, newBlockTimeSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
	if err := monitoring.RegisterSource(BlockEventAndTxsProcessingTime, newBlockProcessingTimeSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// BlockNodeMetricSource is a metric source that captures block properties where the Node is the subject
type BlockNodeMetricSource[T any] struct {
	*utils.SyncedSeriesSource[monitoring.Node, monitoring.BlockNumber, T]
	getBlockProperty func(b monitoring.Block) T
	monitor          *monitoring.Monitor
	subjects         map[monitoring.Node]bool
}

// NewBlockTimeSource creates a metric capturing time of the block finalisation for each Node.
func NewBlockTimeSource(monitor *monitoring.Monitor) *BlockNodeMetricSource[time.Time] {
	f := func(b monitoring.Block) time.Time {
		return b.Time
	}
	return newBlockNodeMetricsSource[time.Time](monitor, f, BlockCompletionTime)
}

// newBlockTimeSource is the same as its public counterpart, it only returns the struct instead of the Source interface
func newBlockTimeSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Node, monitoring.Series[monitoring.BlockNumber, time.Time]] {
	return NewBlockTimeSource(monitor)
}

// NewBlockProcessingTimeSource creates a metric capturing time of the block finalisation for each Node.
func NewBlockProcessingTimeSource(monitor *monitoring.Monitor) *BlockNodeMetricSource[time.Duration] {
	f := func(b monitoring.Block) time.Duration {
		return b.ProcessingTime
	}
	return newBlockNodeMetricsSource[time.Duration](monitor, f, BlockEventAndTxsProcessingTime)
}

// newBlockProcessingTimeSource is the same as its public counterpart, it only returns the struct instead of the Source interface
func newBlockProcessingTimeSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Node, monitoring.Series[monitoring.BlockNumber, time.Duration]] {
	return NewBlockProcessingTimeSource(monitor)
}

// newBlockNodeMetricsSource creates a new data source periodically collecting data from the Node log
func newBlockNodeMetricsSource[T any](
	monitor *monitoring.Monitor,
	getBlockProperty func(b monitoring.Block) T,
	metric monitoring.Metric[monitoring.Node, monitoring.Series[monitoring.BlockNumber, T]]) *BlockNodeMetricSource[T] {

	m := &BlockNodeMetricSource[T]{
		SyncedSeriesSource: utils.NewSyncedSeriesSource(metric),
		getBlockProperty:   getBlockProperty,
		monitor:            monitor,
		subjects:           make(map[monitoring.Node]bool, 10),
	}

	monitor.NodeLogProvider().RegisterLogListener(m)
	monitor.Network().RegisterListener(m)
	for _, node := range monitor.Network().GetActiveNodes() {
		m.AfterNodeCreation(node)
	}

	return m
}

func (s *BlockNodeMetricSource[T]) Shutdown() error {
	s.monitor.NodeLogProvider().UnregisterLogListener(s)
	return s.SyncedSeriesSource.Shutdown()
}

func (s *BlockNodeMetricSource[T]) OnBlock(node monitoring.Node, block monitoring.Block) {
	series := s.GetOrAddSubject(node)
	if err := series.Append(monitoring.BlockNumber(block.Height), s.getBlockProperty(block)); err != nil {
		log.Error("error to add to the series: %s", err)
	}
}

// These are added so that BlockNodeMetricSource can track live nodes only
// e.g. no longer display latest data for nodes that have already ended.
// This is done without touching parents.data, which tracks all historical data.

func (s *BlockNodeMetricSource[T]) AfterNodeCreation(node driver.Node) {
	label := node.GetLabel()
	s.AddSubject(mon.Node(label))
}

func (s *BlockNodeMetricSource[T]) AfterNodeRemoval(node driver.Node) {
	label := node.GetLabel()
	s.RemoveSubject(mon.Node(label))
}

func (s *BlockNodeMetricSource[T]) AfterApplicationCreation(driver.Application) {
	//ignored
}

func (s *BlockNodeMetricSource[T]) GetSubjects() []monitoring.Node {
	res := make([]monitoring.Node, 0, len(s.subjects))
	for subject := range s.subjects {
		res = append(res, subject)
	}
	return res
}

func (s *BlockNodeMetricSource[T]) AddSubject(node monitoring.Node) error {
	s.subjects[node] = true
	return nil
}

func (s *BlockNodeMetricSource[T]) RemoveSubject(node monitoring.Node) error {
	if _, exists := s.subjects[node]; exists {
		delete(s.subjects, node)
	}
	return nil
}

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

// Metric defines a metric in the monitoring system. The type S is the type of
// subject this metric is collected for (e.g., a node or the full network), and
// the type T is the type of value produced by this metric.
//
// The subject of a metric is the object the metric's property is to be
// associated to. There are two major subjects:
//
//   - monitoring.Network ... to be used for network-wide properties like the
//     number of transactions in a block or the utilized gas. Those metrics are
//     consistent throughout the network and do not require any finer
//     granularity.
//
//   - monitoring.Node ... to be used for node-level properties like the time a
//     block was completed, or the CPU usage at a givne time.
//
// Metric data is typically organized in data series, of which there are two
// main types: monitoring.TimeSeries and monitoring.BlockSeries. The former is
// associating a value to various points in time (using absolute time-stamps).
// The latter associates a value to various block-numbers.
type Metric[S any, T any] struct {
	Name        string // used for unique identification of a metric
	Description string // a description documenting the details of the metric
}

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

import (
	"golang.org/x/exp/constraints"
)

// DataPoint is one entry of a data series.
type DataPoint[K constraints.Ordered, T any] struct {
	Position K
	Value    T
}

// Series is a generic interface for arbitrarily indexed sequences of values.
// The type K is the index type, the type T the value associated to the keys.
type Series[K constraints.Ordered, T any] interface {
	// GetRange captures a snapshot of all points collected for the half-open
	// interval [from,to).
	GetRange(from, to K) []DataPoint[K, T]
	// GetLatest retrieves the latest collected data point or nil if no data
	// was collected.
	GetLatest() *DataPoint[K, T]
}

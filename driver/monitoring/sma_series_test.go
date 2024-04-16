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
	"fmt"
	"testing"
)

func TestSmaSeriesIsSeries(t *testing.T) {
	inst := SmaSeries[int, int]{}
	var _ Series[int, int] = &inst
}

func TestSmaSeriesComputeSMA(t *testing.T) {
	input := []float64{1, 3, 5, 10, 30, 50, 100, 300, 500}
	extraInput := []float64{1000, 3000, 5000}

	windowSizes := map[int][]float64{
		2: {1, 2, 4, 7.5, 20, 40, 75, 200, 400, 750, 2000, 4000},
		4: {1, 2, 3, 4.75, 12, 23.75, 47.5, 120, 237.5, 475, 1200, 2375},
		8: {1, 2, 3, 4.75, 9.8, 16.5, 28.42857143, 62.375, 124.75, 249.375, 623.75, 1247.5},
	}

	for window, expected := range windowSizes {
		t.Run(fmt.Sprintf("window: %d", window), func(t *testing.T) {

			series := &SyncedSeries[float64, float64]{}
			for i, val := range input {
				_ = series.Append(float64(i), val)
			}
			smaSeries := NewSMASeries[float64, float64](series, window)

			var start float64
			var size int
			latest := smaSeries.GetLatest()
			points := smaSeries.GetRange(start, latest.Position)
			for i, point := range points {
				if got, want := fmt.Sprintf("%.2f", point.Value), fmt.Sprintf("%.2f", expected[i]); got != want {
					t.Errorf("values do not match: %v != %v", got, want)
				}
				size++
			}

			// check last value
			if got, want := (*smaSeries.GetLatest()).Value, expected[size]; got != want {
				t.Errorf("values do not match: %v != %v", got, want)
			}

			// add more values, SMA will be computed for an extra values
			for i, val := range extraInput {
				if err := series.Append(float64(i+size+1), val); err != nil {
					t.Errorf("cannot insert value: %s", err)
				}
			}
			latest2 := smaSeries.GetLatest()
			extraPoints := append(smaSeries.GetRange(latest.Position, latest2.Position), *latest2)
			shift := size
			for i, point := range extraPoints {
				if got, want := fmt.Sprintf("%.2f", point.Value), fmt.Sprintf("%.2f", expected[i+shift]); got != want {
					t.Errorf("values do not match: %v != %v", got, want)
				}
				size++
			}

			if size != len(expected) {
				t.Errorf("sizes do not match: %d != %d", size, len(expected))
			}
		})
	}
}

func TestSmaSeriesComputeSMASingleValue(t *testing.T) {
	series := &SyncedSeries[float64, float64]{}
	smaSeries := NewSMASeries[float64, float64](series, 4)

	// add value, SMA will be computed for an extra value
	if err := series.Append(float64(0), 1000); err != nil {
		t.Errorf("cannot add value: %s", err)
	}
	if got, want := (*smaSeries.GetLatest()).Value, float64(1000); got != want {
		t.Errorf("values do not match: %v != %v", got, want)
	}
}

func TestSmaSeriesComputeSMANegativeValues(t *testing.T) {
	series := &SyncedSeries[int, float64]{}
	smaSeries := NewSMASeries[int, float64](series, 4)

	values := []int{-5, -6, -7, -8, 8, 7, 6, 5}
	for i, val := range values {
		if err := series.Append(i, float64(val)); err != nil {
			t.Errorf("cannot add value: %s", err)
		}
	}

	if got, want := (*smaSeries.GetLatest()).Value, 6.5; got != want {
		t.Errorf("values do not match: %v != %v", got, want)
	}
}

func TestSmaSeriesEmpty(t *testing.T) {
	series := &SyncedSeries[float64, float64]{}
	smaSeries := NewSMASeries[float64, float64](series, 4)

	if point := smaSeries.GetLatest(); point != nil {
		t.Errorf("point should not exist")
	}
}

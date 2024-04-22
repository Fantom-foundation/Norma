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
	"bytes"
	"testing"
)

func TestCsvExport_HeaderCanBeWritten(t *testing.T) {
	buffer := new(bytes.Buffer)
	if err := WriteCsvHeader(buffer); err != nil {
		t.Fatalf("failed to write CSV header: %v", err)
	}
	want := "run,metric,network,node,app,time,block,workers,value\n"
	if got := buffer.String(); want != got {
		t.Errorf("invalid header, got %s, wanted %s", got, want)
	}
}

func TestCsvExport_WriteLineWithoutOptionals(t *testing.T) {
	record := Record{
		Network: "N",
		Node:    "n",
		App:     "a",
		Value:   "v",
	}

	buffer := new(bytes.Buffer)
	line := CsvRecord{
		Metric: "m",
		Run:    "r",
		Record: record,
	}
	if _, err := line.WriteTo(buffer); err != nil {
		t.Errorf("failed to encode subject: %v", err)
		return
	}

	want := "r, m, N, n, a, , , , v\n"
	if got := buffer.String(); got != want {
		t.Errorf("unexpected encoding, wanted `%v`, got `%v`", want, got)
	}
}

func TestCsvExport_WriteLineWithOptionals(t *testing.T) {
	toInt := func(x int64) *int64 {
		res := new(int64)
		*res = x
		return res
	}
	record := Record{
		Network: "N",
		Node:    "n",
		App:     "a",
		Time:    toInt(1),
		Block:   toInt(2),
		Worker:  toInt(3),
		Value:   "v",
	}

	buffer := new(bytes.Buffer)
	line := CsvRecord{
		Metric: "m",
		Run:    "r",
		Record: record,
	}
	if _, err := line.WriteTo(buffer); err != nil {
		t.Errorf("failed to encode subject: %v", err)
		return
	}

	want := "r, m, N, n, a, 1, 2, 3, v\n"
	if got := buffer.String(); got != want {
		t.Errorf("unexpected encoding, wanted `%v`, got `%v`", want, got)
	}
}

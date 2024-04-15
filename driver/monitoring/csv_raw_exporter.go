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
	"io"
	"strings"
)

// CsvRecord summarizes the content of a single line in the exported raw CSV
// metric dump. It is used as the output format of data sources when exporting
// data to facilitate future extensions / modifications of the format.
type CsvRecord struct {
	Record
	Metric string // must not be empty
	Run    string // The name of the evaluation run
}

// WriteCsvHeader writes a header line defining the fields of the CsvRecord.
func WriteCsvHeader(out io.Writer) error {
	_, err := out.Write([]byte("run,metric,network,node,app,time,block,workers,value\n"))
	return err
}

// WriteTo writes the content of a CsvRecord to the given writer.
func (r *CsvRecord) WriteTo(out io.Writer) (int64, error) {
	toStr := func(x *int64) string {
		if x == nil {
			return ""
		}
		return fmt.Sprintf("%d", *x)
	}

	line := strings.Join([]string{
		r.Run,
		r.Metric,
		r.Network,
		r.Node,
		r.App,
		toStr(r.Time),
		toStr(r.Block),
		toStr(r.Worker),
		r.Value,
	}, ", ") + "\n"
	n, err := out.Write([]byte(line))
	return int64(n), err
}

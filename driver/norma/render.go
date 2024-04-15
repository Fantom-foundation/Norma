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

package main

import (
	"fmt"
	"os"

	"github.com/Fantom-foundation/Norma/analysis/report"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma render`

var renderCommand = cli.Command{
	Action: render,
	Name:   "render",
	Usage:  "renders a report for given monitoring data",
}

func render(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() < 1 {
		return fmt.Errorf("requires measurment file path as argument")
	}

	input := args.First()
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	result, err := report.SingleEvalReport.Render(input, currentDir)
	if err != nil {
		return err
	}

	fmt.Printf("Generated report at file://%s/%s\n", currentDir, result)
	return nil
}

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

	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma check <scenario.yml>`

var checkCommand = cli.Command{
	Action: check,
	Name:   "check",
	Usage:  "checks a scenario configuration file for issues",
}

func check(ctx *cli.Context) (err error) {

	args := ctx.Args()
	if args.Len() < 1 {
		return fmt.Errorf("requires target file name as argument")
	}

	path := args.First()
	fmt.Printf("Trying to parse '%s' ...\n", path)

	scenario, err := parser.ParseFile(path)
	if err != nil {
		return err
	}

	err = scenario.Check()
	if err != nil {
		return err
	}
	fmt.Printf("All checks passed!\n")
	return nil
}

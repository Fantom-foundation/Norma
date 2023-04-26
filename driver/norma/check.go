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

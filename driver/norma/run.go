package main

import (
	"fmt"

	"github.com/Fantom-foundation/Norma/driver/executor"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma run <scenario.yml>`

var runCommand = cli.Command{
	Action: run,
	Name:   "run",
	Usage:  "runs a scenario",
}

func run(ctx *cli.Context) (err error) {

	args := ctx.Args()
	if args.Len() < 1 {
		return fmt.Errorf("requires scenario file as an argument")
	}

	path := args.First()
	fmt.Printf("Running '%s' ...\n", path)

	scenario, err := parser.ParseFile(path)
	if err != nil {
		return err
	}

	clock := executor.NewWallTimeClock()
	net, err := local.NewLocalNetwork()
	if err != nil {
		return err
	}

	err = executor.Run(clock, net, &scenario)
	if err != nil {
		return err
	}
	fmt.Printf("Execution completed successfully!\n")
	return nil
}

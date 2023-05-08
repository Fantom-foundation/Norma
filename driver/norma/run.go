package main

import (
	"fmt"

	"github.com/Fantom-foundation/Norma/driver"
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
	fmt.Printf("Reading '%s' ...\n", path)
	scenario, err := parser.ParseFile(path)
	if err != nil {
		return err
	}

	clock := executor.NewWallTimeClock()

	netConfig := driver.NetworkConfig{
		NumberOfValidators: 1,
	}
	if scenario.NumValidators != nil {
		netConfig.NumberOfValidators = *scenario.NumValidators
	}
	fmt.Printf("Createing network with %d validator(s) ...\n", netConfig.NumberOfValidators)
	net, err := local.NewLocalNetwork(&netConfig)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Printf("Shutting down network ...\n")
		if err := net.Shutdown(); err != nil {
			fmt.Printf("error during network shutdown:\n%v", err)
		}
	}()

	fmt.Printf("Running '%s' ...\n", path)
	err = executor.Run(clock, net, &scenario)
	if err != nil {
		return err
	}
	fmt.Printf("Execution completed successfully!\n")

	return nil
}

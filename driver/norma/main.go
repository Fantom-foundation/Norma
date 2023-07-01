package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma`

func main() {
	app := &cli.App{
		Name:      "Norma Network Runner TEST RUN",
		HelpName:  "norma",
		Usage:     "A set of tools for running network scenarios",
		Copyright: "(c) 2023 Fantom Foundation",
		Flags:     []cli.Flag{},
		Commands: []*cli.Command{
			&checkCommand,
			&runCommand,
			&purgeCommand,
			&renderCommand,
			&diffCommand,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"

	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma purge`

var purgeCommand = cli.Command{
	Action: purge,
	Name:   "purge",
	Usage:  "purges all resources taken by norma",
}

func purge(_ *cli.Context) error {
	fmt.Printf("Purging all resources...\n")
	err := docker.Purge()
	if err != nil {
		return err
	}
	fmt.Printf("Done.\n")
	return nil
}

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
	result, err := report.SingleEvalReport.Render(currentDir+"/"+input, currentDir)
	if err != nil {
		return err
	}

	fmt.Printf("Generated report at file://%s/%s\n", currentDir, result)
	return nil
}

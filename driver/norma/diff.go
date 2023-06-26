package main

import (
	"fmt"
	"os"

	"github.com/Fantom-foundation/Norma/analysis/report"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma diff`

var diffCommand = cli.Command{
	Action: diff,
	Name:   "diff",
	Usage:  "renders a report comparing the monitoring data of multiple evaluations",
}

func diff(ctx *cli.Context) error {
	args := ctx.Args()
	if args.Len() < 1 {
		return fmt.Errorf("requires at least one measurment file path as argument")
	}

	// Merge input files into a temporary file.
	file, err := os.CreateTemp("", "union_*.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	for _, src := range args.Slice() {
		content, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("failed to read file %v: %v", src, err)
		}
		if _, err := file.Write(content); err != nil {
			return err
		}
	}

	if err := file.Close(); err != nil {
		return err
	}
	defer os.Remove(file.Name())

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	result, err := report.MultiEvalReport.Render(file.Name(), currentDir)
	if err != nil {
		return err
	}

	fmt.Printf("Generated report at file://%s/%s\n", currentDir, result)
	return nil
}

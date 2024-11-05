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
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	_ "github.com/Fantom-foundation/Norma/driver/monitoring/app"
	_ "github.com/Fantom-foundation/Norma/driver/monitoring/user"
	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/Fantom-foundation/Norma/driver/runner"
	"github.com/urfave/cli/v2"
)

// Run with `go run ./driver/norma run <scenario.yml>`

var runCommand = cli.Command{
	Action: run,
	Name:   "run",
	Usage:  "runs a scenario",
	Flags: []cli.Flag{
		&evalLabel,
		&keepPrometheusRunning,
		&numValidators,
		&skipChecks,
		&skipReportRendering,
		&outputDirectory,
	},
}

var (
	evalLabel = cli.StringFlag{
		Name:  "label",
		Usage: "define a label for to be added to the monitoring data for this run. If empty, a random label is used.",
		Value: "",
	}
	outputDirectory = cli.StringFlag{
		Name:    "output-directory",
		Usage:   "define a directory at which the monitoring artifact will be saved.",
		Value:   "",
		Aliases: []string{"o"},
	}
	keepPrometheusRunning = cli.BoolFlag{
		Name:    "keep-prometheus-running",
		Usage:   "if set, the Prometheus instance will not be shut down after the run is complete.",
		Aliases: []string{"kpr"},
	}
	numValidators = cli.IntFlag{
		Name:  "num-validators",
		Usage: "overrides the number of validators specified in the scenario file.",
	}
	skipChecks = cli.BoolFlag{
		Name:  "skip-checks",
		Usage: "disables the final network consistency checks",
	}
	skipReportRendering = cli.BoolFlag{
		Name:  "skip-report-rendering",
		Usage: "disables the rendering of the final summary report",
	}
)

func run(ctx *cli.Context) (err error) {
	if num := ctx.Int(numValidators.Name); num != 0 {
		fmt.Printf("[DEPRECATED] --num-validator flag has been deprecated along with NumValidator configuration in scenarios.\n --num-validator %d will not have any effect when running the provided scenarios.", num)
	}

	label := ctx.String(evalLabel.Name)
	if label == "" {
		label = fmt.Sprintf("eval_%d", time.Now().Unix())
	}

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

	// if not configured, default to /tmp/norma_data_<label>_<timestamp> else /configured/path/norma_data_<l>_<t>
	outputDir, err := os.MkdirTemp(ctx.String(outputDirectory.Name), fmt.Sprintf("norma_data_%s_", label))
	if err != nil {
		return fmt.Errorf("Couldn't create temp dir for output; %w", err)
	}

	// create symlink as qol (_latest => _####) where #### is the randomly generated name
	symlink := filepath.Join(filepath.Dir(outputDir), fmt.Sprintf("norma_data_%s_latest", label))
	if _, err := os.Lstat(symlink); err == nil {
		os.Remove(symlink)
	}
	os.Symlink(outputDir, symlink)

	// Copy scenario yml to outputDir as well to provide context
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputDir, filepath.Base(path)), data, 0644)
	if err != nil {
		return err
	}

	err = runner.RunScenario(&scenario, runner.RunConfig{
		Label:                   label,
		OutputDirectory:         &outputDir,
		SkipReportRendering:     ctx.Bool(skipReportRendering.Name),
		SkipCheckNetworkPostRun: ctx.Bool(skipChecks.Name),
		KeepPrometheusRunning:   ctx.Bool(keepPrometheusRunning.Name),
	})
	if err != nil {
		return fmt.Errorf("failed during scenario run; %w", err)
	}

	return nil
}

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

	"github.com/Fantom-foundation/Norma/driver/checking"

	"github.com/Fantom-foundation/Norma/analysis/report"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/executor"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	_ "github.com/Fantom-foundation/Norma/driver/monitoring/app"
	prometheusmon "github.com/Fantom-foundation/Norma/driver/monitoring/prometheus"
	_ "github.com/Fantom-foundation/Norma/driver/monitoring/user"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/parser"
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

	if err := scenario.Check(); err != nil {
		return err
	}

	fmt.Printf("Starting evaluation %s\n", label)

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

	fmt.Printf("Monitoring data is written to %v\n", outputDir)

	// Copy scenario yml to outputDir as well to provide context
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputDir, filepath.Base(path)), data, 0644)
	if err != nil {
		return err
	}

	clock := executor.NewWallTimeClock()

	// Startup network.
	fmt.Printf("Creating network with: \n")
	fmt.Printf("    Network max block gas: %d\n", scenario.GetMaxBlockGas())
	fmt.Printf("    Network max epoch gas: %d\n", scenario.GetMaxEpochGas())

	net, err := local.NewLocalNetwork(&driver.NetworkConfig{
		NumberOfValidators: scenario.GetNumValidators(),
		MaxBlockGas:        scenario.GetMaxBlockGas(),
		MaxEpochGas:        scenario.GetMaxEpochGas(),
	})
	if err != nil {
		return err
	}
	defer func() {
		fmt.Printf("Shutting down network ...\n")
		if err := net.Shutdown(); err != nil {
			fmt.Printf("error during network shutdown:\n%v", err)
		}
	}()

	// Initialize monitoring environment.
	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{
		EvaluationLabel: label,
		OutputDir:       outputDir,
	})
	if err != nil {
		return err
	}
	defer func() {
		fmt.Printf("Shutting down data monitor ...\n")
		if err := monitor.Shutdown(); err != nil {
			fmt.Printf("error during monitor shutdown:\n%v\n", err)
		}
		fmt.Printf("Monitoring data was written to %v\n", outputDir)
		fmt.Printf("Raw data was exported to %s\n", monitor.GetMeasurementFileName())

		if !ctx.Bool(skipReportRendering.Name) {
			fmt.Printf("Rendering summary report (may take a few minutes the first time if R packages need to be installed) ...\n")
			if file, err := report.SingleEvalReport.Render(monitor.GetMeasurementFileName(), outputDir); err != nil {
				fmt.Printf("Report generation failed:\n%v\n", err)
			} else {
				fmt.Printf("Summary report was exported to file://%s/%s\n", outputDir, file)
			}
		} else {
			fmt.Printf("Report rendering skipped (--%s)\n", skipReportRendering.Name)
			fmt.Printf("To render report run `norma render %s`\n", monitor.GetMeasurementFileName())
		}
	}()

	// Install monitoring sensory.
	if err := monitoring.InstallAllRegisteredSources(monitor); err != nil {
		return err
	}

	// Run prometheus.
	fmt.Printf("Starting Prometheus ...\n")
	prom, err := prometheusmon.Start(net, net.GetDockerNetwork())
	if err != nil {
		fmt.Printf("error starting Prometheus:\n%v", err)
	}
	defer func() {
		if !ctx.Bool(keepPrometheusRunning.Name) && prom != nil {
			fmt.Printf("Shutting down Prometheus ...\n")
			if err := prom.Shutdown(); err != nil {
				fmt.Printf("error during Prometheus shutdown:\n%v", err)
			}
		}
	}()

	// Run scenario.
	fmt.Printf("Running '%s' ...\n", path)
	logger := startProgressLogger(monitor, net)
	defer logger.shutdown()
	err = executor.Run(clock, net, &scenario, outputDir)
	if err != nil {
		return err
	}
	fmt.Printf("Execution completed successfully!\n")

	if !ctx.Bool(skipChecks.Name) {
		fmt.Printf("Checking network consistency ...\n")
		err = checking.CheckNetworkConsistency(net)
		if err != nil {
			return fmt.Errorf("checking the network consistency failed: %v", err)
		}
		fmt.Printf("Network checks succeed.\n")
	} else {
		fmt.Printf("Network checks skipped (--%s)\n", skipChecks.Name)
	}

	return nil
}

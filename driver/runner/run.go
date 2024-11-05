package runner

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/analysis/report"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/checking"
	"github.com/Fantom-foundation/Norma/driver/executor"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	netmon "github.com/Fantom-foundation/Norma/driver/monitoring/network"
	nodemon "github.com/Fantom-foundation/Norma/driver/monitoring/node"
	prometheusmon "github.com/Fantom-foundation/Norma/driver/monitoring/prometheus"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/parser"
	"golang.org/x/exp/constraints"
)

type RunConfig struct {
	Label                   string
	OutputDirectory         *string
	SkipReportRendering     bool
	SkipCheckNetworkPostRun bool
	KeepPrometheusRunning   bool
}

// RunScenario runs the provided scenario
func RunScenario(scenario *parser.Scenario, config RunConfig) error {
	var outdir string
	if config.OutputDirectory != nil {
		outdir = *config.OutputDirectory
	} else {
		od, err := os.MkdirTemp(config.Label, "*")
		if err != nil {
			return fmt.Errorf("failed to create tmp out dir; %w", err)
		}
		outdir = od
	}

	if err := scenario.Check(); err != nil {
		return err
	}

	fmt.Printf("Starting evaluation %s\n", config.Label)
	fmt.Printf("Monitoring data is written to %v\n", outdir)

	static, dynamic := scenario.GetStaticDynamicValidatorCount()
	mandatory := scenario.GetMandatoryValidatorCount()

	// start network
	net, err := local.NewLocalNetwork(&driver.NetworkConfig{
		TotalNumberOfValidators: static + dynamic + mandatory,
		NumberOfValidators:      mandatory,
		MaxBlockGas:             scenario.GetMaxBlockGas(),
		MaxEpochGas:             scenario.GetMaxEpochGas(),
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
		EvaluationLabel: config.Label,
		OutputDir:       outdir,
	})
	if err != nil {
		return err
	}
	defer func() {
		fmt.Printf("Shutting down data monitor ...\n")
		if err := monitor.Shutdown(); err != nil {
			fmt.Printf("error during monitor shutdown:\n%v\n", err)
		}

		fmt.Printf("Monitoring data was written to %v\n", outdir)
		fmt.Printf("Raw data was exported to %s\n", monitor.GetMeasurementFileName())

		if !config.SkipReportRendering && config.OutputDirectory != nil {
			fmt.Printf("Rendering summary report (may take a few minutes the first time if R packages need to be installed) ...\n")
			if file, err := report.SingleEvalReport.Render(monitor.GetMeasurementFileName(), outdir); err != nil {
				fmt.Printf("Report generation failed:\n%v\n", err)
			} else {
				fmt.Printf("Summary report was exported to file://%s/%s\n", outdir, file)
			}
		} else {
			fmt.Printf("Report rendering skipped\n")
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
		if !config.KeepPrometheusRunning && prom != nil {
			fmt.Printf("Shutting down Prometheus ...\n")
			if err := prom.Shutdown(); err != nil {
				fmt.Printf("error during Prometheus shutdown:\n%v", err)
			}
		}
	}()

	// Run scenario.
	fmt.Printf("Running '%s' ...\n", scenario.Name)

	logger := startProgressLogger(monitor)
	defer logger.shutdown()

	clock := executor.NewWallTimeClock()
	err = executor.Run(clock, net, scenario, outdir, logger.epochTracker)
	if err != nil {
		return err
	}
	fmt.Printf("Execution completed successfully!\n")

	if !config.SkipCheckNetworkPostRun {
		fmt.Printf("Checking network consistency ...\n")
		err = checking.CheckNetworkConsistency(net)
		if err != nil {
			return fmt.Errorf("checking the network consistency failed: %v", err)
		}
		fmt.Printf("Network checks succeed.\n")
	} else {
		fmt.Printf("Network checks skipped.\n")
	}

	return nil
}

type progressLogger struct {
	monitor      *monitoring.Monitor
	epochTracker map[monitoring.Node]string
	stop         chan<- bool
	done         <-chan bool
}

func startProgressLogger(monitor *monitoring.Monitor) *progressLogger {
	epochTracker := map[monitoring.Node]string{}
	stop := make(chan bool)
	done := make(chan bool)

	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				logState(monitor, epochTracker)
			}
		}
	}()

	return &progressLogger{
		monitor,
		epochTracker,
		stop,
		done,
	}
}

func (l *progressLogger) shutdown() {
	close(l.stop)
	<-l.done
}

func logState(monitor *monitoring.Monitor, epochTracker map[monitoring.Node]string) {
	numNodes := getNumNodes(monitor)
	blockStatuses := getBlockStatuses(monitor, epochTracker)
	txPers := getTxPerSec(monitor)
	txs := getNumTxs(monitor)
	gas := getGasUsed(monitor)
	processingTimes := getBlockProcessingTimes(monitor)

	log.Printf("Nodes: %s, block heights: %v, tx/s: %v, txs: %v, gas: %s, block processing: %v", numNodes, blockStatuses, txPers, txs, gas, processingTimes)
}

func getNumNodes(monitor *monitoring.Monitor) string {
	data, exists := monitoring.GetData(monitor, monitoring.Network{}, netmon.NumberOfNodes)
	return getLastValAsString[monitoring.Time, int](exists, data)
}

func getNumTxs(monitor *monitoring.Monitor) string {
	data, exists := monitoring.GetData(monitor, monitoring.Network{}, netmon.BlockNumberOfTransactions)
	return getLastValAsString[monitoring.BlockNumber, int](exists, data)
}

func getTxPerSec(monitor *monitoring.Monitor) []string {
	metric := nodemon.TransactionsThroughput
	return getLastValAllSubjects[monitoring.BlockNumber, float32](monitor, metric, nil)
}

func getGasUsed(monitor *monitoring.Monitor) string {
	data, exists := monitoring.GetData(monitor, monitoring.Network{}, netmon.BlockGasUsed)
	return getLastValAsString[monitoring.BlockNumber, int](exists, data)
}

func getBlockStatuses(monitor *monitoring.Monitor, epochTracker map[monitoring.Node]string) []string {
	metric := nodemon.NodeBlockStatus
	return getLastValAllSubjects[monitoring.Time, monitoring.BlockStatus, monitoring.Series[monitoring.Time, monitoring.BlockStatus]](monitor, metric, epochTracker)
}

func getBlockProcessingTimes(monitor *monitoring.Monitor) []string {
	metric := nodemon.BlockEventAndTxsProcessingTime
	return getLastValAllSubjects[monitoring.BlockNumber, time.Duration, monitoring.Series[monitoring.BlockNumber, time.Duration]](monitor, metric, nil)
}

func getLastValAllSubjects[K constraints.Ordered, T any, X monitoring.Series[K, T]](monitor *monitoring.Monitor, metric monitoring.Metric[monitoring.Node, X], epochTracker map[monitoring.Node]string) []string {
	nodes := monitoring.GetSubjects(monitor, metric)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i] < nodes[j] })

	res := make([]string, 0, len(nodes))
	for _, node := range nodes {
		data, exists := monitoring.GetData(monitor, node, metric)
		var d string = getLastValAsString[K, T](exists, data)
		res = append(res, d)

		if epochTracker != nil {
			epochTracker[node] = strings.Split(d, "/")[0]
		}
	}
	return res
}

func getLastValAsString[K constraints.Ordered, T any](exists bool, series monitoring.Series[K, T]) string {
	if !exists || series == nil {
		return "N/A"
	}
	point := series.GetLatest()
	if point == nil {
		return "N/A"
	}
	return fmt.Sprintf("%v", point.Value)
}

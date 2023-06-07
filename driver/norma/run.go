package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/executor"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	_ "github.com/Fantom-foundation/Norma/driver/monitoring/app"
	netmon "github.com/Fantom-foundation/Norma/driver/monitoring/network"
	nodemon "github.com/Fantom-foundation/Norma/driver/monitoring/node"
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
		&DbImpl,
	},
}

var (
	DbImpl = cli.StringFlag{
		Name:  "db-impl",
		Usage: "select the DB implementation to use (geth or carmen)",
		Value: "carmen",
	}
)

func run(ctx *cli.Context) (err error) {

	db := strings.ToLower(ctx.String(DbImpl.Name))
	if db == "carmen" || db == "go-file" {
		db = "go-file"
	} else if db != "geth" {
		return fmt.Errorf("unknown value fore --%v flag: %v", DbImpl.Name, db)
	}

	args := ctx.Args()
	if args.Len() < 1 {
		return fmt.Errorf("requires scenario file as an argument")
	}

	outputDir, err := os.MkdirTemp("", "norma_data_")
	if err != nil {
		return err
	}
	fmt.Printf("Monitoring data is written to %v\n", outputDir)

	path := args.First()
	fmt.Printf("Reading '%s' ...\n", path)
	scenario, err := parser.ParseFile(path)
	if err != nil {
		return err
	}

	clock := executor.NewWallTimeClock()

	// Startup network.
	netConfig := driver.NetworkConfig{
		NumberOfValidators:    1,
		StateDbImplementation: db,
	}
	if scenario.NumValidators != nil {
		netConfig.NumberOfValidators = *scenario.NumValidators
	}
	fmt.Printf("Creating network with %d validator(s) using the `%v` DB implementation ...\n", netConfig.NumberOfValidators, netConfig.StateDbImplementation)
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

	// Initialize monitoring environment.
	csvPath := "./output.csv"
	csvFile, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	monitor := monitoring.NewMonitor(net, monitoring.MonitorConfig{
		OutputDir: outputDir,
	}, monitoring.NewWriterChain(csvFile))
	defer func() {
		// TODO: dump data before shutting down monitor
		fmt.Printf("Shutting down data monitor ...\n")
		if err := monitor.Shutdown(); err != nil {
			fmt.Printf("error during monitor shutdown:\n%v", err)
		}
		_ = csvFile.Close()
		fmt.Printf("Monitoring data was written to %v\n", outputDir)
	}()

	// Install monitoring sensory.
	if err := monitoring.InstallAllRegisteredSources(monitor); err != nil {
		return err
	}

	// Run scenario.
	fmt.Printf("Running '%s' ...\n", path)
	logger := startProgressLogger(monitor)
	defer logger.shutdown()
	err = executor.Run(clock, net, &scenario)
	if err != nil {
		return err
	}
	fmt.Printf("Execution completed successfully!\n")

	return nil
}

type progressLogger struct {
	monitor *monitoring.Monitor
	stop    chan<- bool
	done    <-chan bool
}

func startProgressLogger(monitor *monitoring.Monitor) *progressLogger {
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
				logState(monitor)
			}
		}
	}()

	return &progressLogger{
		monitor,
		stop,
		done,
	}
}

func (l *progressLogger) shutdown() {
	close(l.stop)
	<-l.done
}

func logState(monitor *monitoring.Monitor) {
	numNodes := getNumNodes(monitor)
	blockHeights := getBlockHeights(monitor)
	txPers := getTxPerSec(monitor)
	txs := getNumTxs(monitor)
	gas := getGasUsed(monitor)
	processingTimes := getBlockProcessingTimes(monitor)
	log.Printf("Nodes: %s, block heights: %v, tx/s: %v, txs: %v, gas: %s, block processing: %v", numNodes, blockHeights, txPers, txs, gas, processingTimes)
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
	return getLastValAllSubjects[monitoring.BlockNumber, float32](monitor, metric)
}

func getGasUsed(monitor *monitoring.Monitor) string {
	data, exists := monitoring.GetData(monitor, monitoring.Network{}, netmon.BlockGasUsed)
	return getLastValAsString[monitoring.BlockNumber, int](exists, data)
}

func getBlockHeights(monitor *monitoring.Monitor) []string {
	metric := nodemon.NodeBlockHeight
	return getLastValAllSubjects[monitoring.Time, int, monitoring.TimeSeries[int]](monitor, metric)
}

func getBlockProcessingTimes(monitor *monitoring.Monitor) []string {
	metric := nodemon.BlockEventAndTxsProcessingTime
	return getLastValAllSubjects[monitoring.BlockNumber, time.Duration, monitoring.BlockSeries[time.Duration]](monitor, metric)
}

func getLastValAllSubjects[K constraints.Ordered, T any, X monitoring.Series[K, T]](monitor *monitoring.Monitor, metric monitoring.Metric[monitoring.Node, X]) []string {
	nodes := monitoring.GetSubjects(monitor, metric)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i] < nodes[j] })

	res := make([]string, 0, len(nodes))
	for _, node := range nodes {
		data, exists := monitoring.GetData(monitor, node, metric)
		res = append(res, getLastValAsString[K, T](exists, data))
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

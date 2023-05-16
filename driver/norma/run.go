package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/executor"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
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

	// Startup network.
	netConfig := driver.NetworkConfig{
		NumberOfValidators: 1,
	}
	if scenario.NumValidators != nil {
		netConfig.NumberOfValidators = *scenario.NumValidators
	}
	fmt.Printf("Creating network with %d validator(s) ...\n", netConfig.NumberOfValidators)
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
	monitor := monitoring.NewMonitor(net)
	defer func() {
		// TODO: dump data before shutting down monitor
		fmt.Printf("Shutting down data monitor ...\n")
		if err := monitor.Shutdown(); err != nil {
			fmt.Printf("error during monitor shutdown:\n%v", err)
		}
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
	log.Printf("Block height of %s nodes in the network: %v", numNodes, blockHeights)
}

func getNumNodes(monitor *monitoring.Monitor) string {
	data := monitoring.GetData(monitor, monitoring.Network{}, netmon.NumberOfNodes)
	if data == nil {
		return "N/A"
	}
	point := (*data).GetLatest()
	if point == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d", point.Value)
}

func getBlockHeights(monitor *monitoring.Monitor) []string {
	metric := nodemon.NodeBlockHeight
	nodes := monitoring.GetSubjects(monitor, metric)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i] < nodes[j] })

	res := make([]string, 0, len(nodes))
	for _, node := range nodes {
		data := monitoring.GetData(monitor, node, metric)
		if data == nil {
			res = append(res, "N/A")
			continue
		}
		point := (*data).GetLatest()
		if point == nil {
			res = append(res, "N/A")
			continue
		}
		res = append(res, fmt.Sprintf("%d", point.Value))
	}
	return res
}

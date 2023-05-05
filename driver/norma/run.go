package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/executor"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/ethereum/go-ethereum/rpc"
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

	fmt.Printf("Createing network ...\n")
	net, err := local.NewLocalNetwork()
	if err != nil {
		return err
	}
	defer func() {
		fmt.Printf("Shutting down network ...\n")
		if err := net.Shutdown(); err != nil {
			fmt.Printf("error during network shutdown:\n%v", err)
		}
	}()

	// Starting monitoring.
	monitor := StartMonitor(net)
	defer monitor.Stop()

	fmt.Printf("Running '%s' ...\n", path)
	err = executor.Run(clock, net, &scenario)
	if err != nil {
		return err
	}
	fmt.Printf("Execution completed successfully!\n")

	return nil
}

// ----------------------------------------------------------------------------

// TODO: move stuff below to monitoring package

type Monitor struct {
	stop chan<- bool
	done <-chan bool
}

func StartMonitor(net driver.Network) Monitor {
	stop := make(chan bool)
	done := make(chan bool)

	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				listBlockHeights(net)
			case <-stop:
				return
			}
		}
	}()

	return Monitor{stop, done}
}

func (m *Monitor) Stop() {
	close(m.stop)
	<-m.done
}

func listBlockHeights(net driver.Network) {
	type Entry struct {
		id    string
		block int
	}

	var data = []Entry{}
	for _, node := range net.GetActiveNodes() {
		id, err := node.GetNodeID()
		if err != nil {
			fmt.Printf("failed to obtain node id")
		}
		block, err := getBlockHeight(node)
		if err != nil {
			fmt.Printf("failed to obtain block height from node: %v", err)
		} else {
			data = append(data, Entry{string(id), block})
		}
	}

	// Sort nodes by ID to have a stable order.
	sort.Slice(data, func(i, j int) bool { return data[i].id < data[j].id })

	heights := make([]int, 0, len(data))
	for _, entry := range data {
		heights = append(heights, entry.block)
	}
	log.Printf("Block heights: %v", heights)
}

func getBlockHeight(node driver.Node) (int, error) {
	url := node.GetRpcServiceUrl()
	if url == nil {
		return 0, fmt.Errorf("node does not export an RPC server")
	}
	rpcClient, err := rpc.DialContext(context.Background(), string(*url))
	if err != nil {
		return 0, err
	}
	var result string
	err = rpcClient.Call(&result, "eth_blockNumber")
	if err != nil {
		return 0, err
	}
	result = strings.TrimPrefix(result, "0x")
	res, err := strconv.ParseInt(result, 16, 32)
	return int(res), err
}

package main

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/Fantom-foundation/Norma/driver/network/rpc"
	rpc2 "github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/core/types"
	erpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli/v2"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Run with `go run ./driver/norma feed --rpc http://127.0.0.1:18545/ --rate 5 --users 2 --app erc20 --private-key file --network-id 0xFA3`

var feedCommand = cli.Command{
	Action: feed,
	Name:   "feed",
	Usage:  "feeds transactions into external RPC nodes",
	Flags: []cli.Flag{
		&rpcUrlFlag,
		&rateFlag,
		&usersFlag,
		&appFlag,
		&privateKeyFlag,
		&networkIdFlag,
	},
}

var (
	rpcUrlFlag = cli.StringSliceFlag{
		Name:     "rpc",
		Usage:    "RPC URL to use (use multiple times for multiple targets)",
		Required: true,
	}
	rateFlag = cli.Float64Flag{
		Name:     "rate",
		Usage:    "Transactions per second",
		Required: true,
	}
	usersFlag = cli.IntFlag{
		Name:     "users",
		Usage:    "Amount of sending addresses",
		Required: true,
	}
	appFlag = cli.StringFlag{
		Name:     "app",
		Usage:    "Type of the blockchain app (erc20, uniswap)",
		Required: true,
	}
	privateKeyFlag = cli.StringFlag{
		Name:      "private-key",
		Usage:     "File with private key for priming (in hex, no 0x prefix)",
		Required:  true,
		TakesFile: true,
	}
	networkIdFlag = cli.Int64Flag{
		Name:     "network-id",
		Usage:    "Network ID of the network",
		Required: true,
	}
)

func feed(ctx *cli.Context) (err error) {
	primaryKey, err := os.ReadFile(ctx.String(privateKeyFlag.Name))
	if err != nil {
		return fmt.Errorf("failed to read private key; %v", err)
	}

	primaryAccount, err := app.NewAccount(0, strings.TrimSpace(string(primaryKey)), ctx.Int64(networkIdFlag.Name))
	if err != nil {
		return fmt.Errorf("failed to create primary account; %v", err)
	}

	net := NewFeederNetworkConnection(ctx.StringSlice(rpcUrlFlag.Name))
	defer net.Close()

	rpcClient, err := net.DialRandomRpc()
	if err != nil {
		return fmt.Errorf("failed to connect to RPC to initialize the application; %v", err)
	}
	defer rpcClient.Close()

	if err := primaryAccount.LoadNonceFromNetwork(rpcClient); err != nil {
		return fmt.Errorf("failed to load nonce of the primary account from the network; %v", err)
	}

	application, err := app.NewApplication(ctx.String(appFlag.Name), rpcClient, primaryAccount, ctx.Int(usersFlag.Name))
	if err != nil {
		return fmt.Errorf("failed to initialize on-chain app; %v", err)
	}

	sh := shaper.NewConstantShaper(ctx.Float64(rateFlag.Name))

	appController, err := controller.NewAppController(application, sh, ctx.Int(usersFlag.Name), net)
	if err != nil {
		return err
	}

	log.Printf("running...\n")
	return appController.Run(context.Background())
}

func NewFeederNetworkConnection(rpcUrls []string) *FeederNetworkConnection {
	rpcWorkerPool := rpc.NewRpcWorkerPool()

	for _, rpcUrl := range rpcUrls {
		rpcWorkerPool.AddRpcNode(driver.URL(rpcUrl))
	}

	return &FeederNetworkConnection{
		rpcWorkerPool: rpcWorkerPool,
		rpcUrls:       rpcUrls,
	}
}

type FeederNetworkConnection struct {
	rpcWorkerPool *rpc.RpcWorkerPool
	rpcUrls       []string
}

func (n *FeederNetworkConnection) SendTransaction(tx *types.Transaction) {
	n.rpcWorkerPool.SendTransaction(tx)
}

func (n *FeederNetworkConnection) DialRandomRpc() (rpc2.RpcClient, error) {
	rpcUrl := n.rpcUrls[rand.Intn(len(n.rpcUrls))]
	rpcClient, err := network.RetryReturn(network.DefaultRetryAttempts, 1*time.Second, func() (*erpc.Client, error) {
		return erpc.DialContext(context.Background(), rpcUrl)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC for %s; %v", rpcUrl, err)
	}
	return rpc2.WrapRpcClient(rpcClient), nil
}

func (n *FeederNetworkConnection) Close() error {
	return n.rpcWorkerPool.Close()
}

package local

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// LocalNetwork is a Docker based network running each individual node
// within its own, dedicated Docker Container.
type LocalNetwork struct {
	docker  *docker.Client
	nodes   map[driver.NodeID]*node.OperaNode
	apps    []driver.Application
	primary *node.OperaNode // first node generated, always the only validator for now
}

func NewLocalNetwork() (driver.Network, error) {
	client, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	return &LocalNetwork{
		docker: client,
		nodes:  map[driver.NodeID]*node.OperaNode{},
		apps:   []driver.Application{},
	}, nil
}

func (n *LocalNetwork) CreateNode(config *driver.NodeConfig) (driver.Node, error) {
	// TODO: support more than one validator node
	isValidator := len(n.nodes) == 0 // for now, only the first node is a validator
	node, err := node.StartOperaDockerNode(n.docker, isValidator)
	if err != nil {
		return nil, err
	}

	id, err := node.GetNodeID()
	for _, other := range n.nodes {
		if err := other.AddPeer(id); err != nil {
			return nil, err
		}
	}

	if len(n.nodes) == 0 {
		n.primary = node
	}

	n.nodes[id] = node
	return node, err
}

// reasureAccountPrivateKey is an account with tokens that can be used to
// initiate test applications and accounts.
const treasureAccountPrivateKey = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7" // Fakenet validator 1

const fakeNetworkID = 0xfa3

type localApplication struct {
	controller *controller.AppController
	cancel     context.CancelFunc
}

func (a *localApplication) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	go func() {
		err := a.controller.Run(ctx)
		if err != nil {
			log.Printf("Failed to run load generator: %v", err)
		}
	}()
	return nil
}

func (a *localApplication) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	a.cancel = nil
	return nil
}

func (n *LocalNetwork) CreateApplication(config *driver.ApplicationConfig) (driver.Application, error) {

	if n.primary == nil {
		return nil, fmt.Errorf("network is empty")
	}

	url := n.primary.GetRpcServiceUrl()
	if url == nil {
		return nil, fmt.Errorf("primary node is not running an RPC server")
	}

	rpcClient, err := ethclient.Dial(string(*url))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(treasureAccountPrivateKey)
	if err != nil {
		return nil, err
	}

	counterGenerator, err := generator.NewCounterTransactionGenerator(privateKey, big.NewInt(fakeNetworkID))
	if err != nil {
		return nil, err
	}

	constantShaper := shaper.NewConstantShaper(config.Rate)

	sourceDriver := controller.NewAppController(counterGenerator, constantShaper, rpcClient)
	err = sourceDriver.Init()
	if err != nil {
		return nil, err
	}

	app := &localApplication{
		controller: sourceDriver,
	}

	n.apps = append(n.apps, app)
	return app, nil
}

func (n *LocalNetwork) GetActiveNodes() []driver.Node {
	res := make([]driver.Node, 0, len(n.nodes))
	for _, node := range n.nodes {
		if node.IsRunning() {
			res = append(res, node)
		}
	}
	return res
}

func (n *LocalNetwork) Shutdown() error {
	var errs []error
	// First stop all generators.
	for _, app := range n.apps {
		// TODO: shutdown apps in parallel.
		if err := app.Stop(); err != nil {
			errs = append(errs, err)
		}
	}
	n.apps = n.apps[:0]

	// Second, shut down the nodes.
	for _, node := range n.nodes {
		// TODO: shutdown nodes in parallel.
		if err := node.Stop(); err != nil {
			errs = append(errs, err)
		}
		if err := node.Cleanup(); err != nil {
			errs = append(errs, err)
		}
	}
	n.nodes = map[driver.NodeID]*node.OperaNode{}

	return errors.Join(errs...)
}

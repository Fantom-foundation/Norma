package local

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sync"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/node"
	opera "github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/crypto"
)

// LocalNetwork is a Docker based network running each individual node
// within its own, dedicated Docker Container.
type LocalNetwork struct {
	docker *docker.Client
	config driver.NetworkConfig

	// validators lists the validator nodes in the network. Validators
	// are created during network startup and run for the full duration
	// of the network.
	validators []*node.OperaNode

	// nodes provide a register for all nodes in the network, including
	// validator nodes created during startup.
	nodes map[driver.NodeID]*node.OperaNode

	// nodesMutex synchronizes access to the list of nodes.
	nodesMutex sync.Mutex

	// apps maintains a list of all applications created on the network.
	apps []driver.Application

	// appsMutex synchronizes access to the list of applications.
	appsMutex sync.Mutex

	// listeners is the set of registered NetworkListeners.
	listeners map[driver.NetworkListener]bool

	// listenerMutex is synching access to listeners
	listenerMutex sync.Mutex
}

func NewLocalNetwork(config *driver.NetworkConfig) (driver.Network, error) {
	client, err := docker.NewClient()
	if err != nil {
		return nil, err
	}

	// Create the empty network.
	net := &LocalNetwork{
		docker:    client,
		config:    *config,
		nodes:     map[driver.NodeID]*node.OperaNode{},
		apps:      []driver.Application{},
		listeners: map[driver.NetworkListener]bool{},
	}

	// Start all validators.
	nodeConfig := node.OperaNodeConfig{
		ValidatorId:   new(int),
		NetworkConfig: config,
	}
	var errs []error
	for i := 0; i < config.NumberOfValidators; i++ {
		// TODO: create nodes in parallel
		*nodeConfig.ValidatorId = i + 1
		nodeConfig.Label = fmt.Sprintf("Validator-%d", i+1)
		validator, err := net.createNode(&nodeConfig)
		if err != nil {
			errs = append(errs, err)
		} else {
			net.validators = append(net.validators, validator)
		}
	}

	// If starting the validators failed, the network statup should fail.
	if len(errs) > 0 {
		err := net.Shutdown()
		if err != nil {
			errs = append(errs, err)
		}
		return nil, errors.Join(errs...)
	}

	return net, nil
}

// createNode is an internal version of CreateNode enabling the creation
// of validator and non-validator nodes in the network.
func (n *LocalNetwork) createNode(nodeConfig *node.OperaNodeConfig) (*node.OperaNode, error) {
	node, err := node.StartOperaDockerNode(n.docker, nodeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start opera docker; %v", err)
	}

	n.nodesMutex.Lock()
	id, err := node.GetNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to get node id; %v", err)
	}
	for _, other := range n.nodes {
		if err = other.AddPeer(id); err != nil {
			n.nodesMutex.Unlock()
			return nil, fmt.Errorf("failed to add peer; %v", err)
		}
	}
	n.nodes[id] = node
	n.nodesMutex.Unlock()

	for _, listener := range n.getListeners() {
		listener.AfterNodeCreation(node)
	}

	return node, nil
}

// CreateNode creates non-validator nodes in the network.
func (n *LocalNetwork) CreateNode(config *driver.NodeConfig) (driver.Node, error) {
	return n.createNode(&node.OperaNodeConfig{
		Label:         config.Name,
		NetworkConfig: &n.config,
	})
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

	node, err := n.getRandomValidator()
	if err != nil {
		return nil, err
	}

	rpcUrl := node.GetServiceUrl(&opera.OperaWsService)
	if rpcUrl == nil {
		return nil, fmt.Errorf("primary node is not running an RPC server")
	}

	privateKey, err := crypto.HexToECDSA(treasureAccountPrivateKey)
	if err != nil {
		return nil, err
	}

	generatorFactory, err := generator.NewCounterGeneratorFactory(*rpcUrl, privateKey, big.NewInt(fakeNetworkID))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tx generator; %v", err)
	}

	constantShaper := shaper.NewConstantShaper(config.Rate)

	appController, err := controller.NewAppController(generatorFactory, constantShaper, config.Accounts)
	if err != nil {
		return nil, err
	}

	app := &localApplication{
		controller: appController,
	}

	n.appsMutex.Lock()
	n.apps = append(n.apps, app)
	n.appsMutex.Unlock()

	for _, listener := range n.getListeners() {
		listener.AfterApplicationCreation(app)
	}

	return app, nil
}

func (n *LocalNetwork) GetActiveNodes() []driver.Node {
	n.nodesMutex.Lock()
	defer n.nodesMutex.Unlock()
	res := make([]driver.Node, 0, len(n.nodes))
	for _, node := range n.nodes {
		if node.IsRunning() {
			res = append(res, node)
		}
	}
	return res
}

func (n *LocalNetwork) RegisterListener(listener driver.NetworkListener) {
	n.listenerMutex.Lock()
	n.listeners[listener] = true
	n.listenerMutex.Unlock()
}

func (n *LocalNetwork) UnregisterListener(listener driver.NetworkListener) {
	n.listenerMutex.Lock()
	delete(n.listeners, listener)
	n.listenerMutex.Unlock()
}

func (n *LocalNetwork) getListeners() []driver.NetworkListener {
	n.listenerMutex.Lock()
	res := make([]driver.NetworkListener, 0, len(n.listeners))
	for listener := range n.listeners {
		res = append(res, listener)
	}
	n.listenerMutex.Unlock()
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

func (n *LocalNetwork) getRandomValidator() (driver.Node, error) {
	if len(n.validators) == 0 {
		return nil, fmt.Errorf("network is empty")
	}
	return n.validators[rand.Intn(len(n.validators))], nil
}

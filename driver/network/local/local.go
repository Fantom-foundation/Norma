package local

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"

	rpc2 "github.com/Fantom-foundation/Norma/driver/rpc"

	"github.com/Fantom-foundation/Norma/driver/network/rpc"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/shaper"
)

// LocalNetwork is a Docker based network running each individual node
// within its own, dedicated Docker Container.
type LocalNetwork struct {
	docker         *docker.Client
	network        *docker.Network
	config         driver.NetworkConfig
	primaryAccount *app.Account

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

	rpcWorkerPool *rpc.RpcWorkerPool
}

func NewLocalNetwork(config *driver.NetworkConfig) (*LocalNetwork, error) {
	client, err := docker.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client; %v", err)
	}

	dn, err := client.CreateBridgeNetwork()
	if err != nil {
		return nil, fmt.Errorf("failed to create bridge network; %v", err)
	}

	// Create chain account, which will be used for the initialization
	primaryAccount, err := app.NewAccount(0, treasureAccountPrivateKey, fakeNetworkID)
	if err != nil {
		return nil, fmt.Errorf("failed to create primary account; %v", err)
	}

	// Create the empty network.
	net := &LocalNetwork{
		docker:         client,
		network:        dn,
		config:         *config,
		primaryAccount: primaryAccount,
		nodes:          map[driver.NodeID]*node.OperaNode{},
		apps:           []driver.Application{},
		listeners:      map[driver.NetworkListener]bool{},
		rpcWorkerPool:  rpc.NewRpcWorkerPool(),
	}

	// Let the RPC pool to start RPC workers when a node start.
	net.RegisterListener(net.rpcWorkerPool)

	// Start all validators.
	net.validators = make([]*node.OperaNode, config.NumberOfValidators)
	errs := make([]error, config.NumberOfValidators)
	var wg sync.WaitGroup
	for i := 0; i < config.NumberOfValidators; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			validatorId := i + 1
			nodeConfig := node.OperaNodeConfig{
				ValidatorId:      &validatorId,
				NetworkConfig:    config,
				Label:            fmt.Sprintf("_validator-%d", validatorId),
				VmImplementation: config.VmImplementation,
			}
			net.validators[i], errs[i] = net.createNode(&nodeConfig)
		}()
	}
	wg.Wait()

	// If starting the validators failed, the network startup should fail.
	if errors.Join(errs...) != nil {
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
	node, err := node.StartOperaDockerNode(n.docker, n.network, nodeConfig)
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

	n.listenerMutex.Lock()
	for listener := range n.listeners {
		listener.AfterNodeCreation(node)
	}
	n.listenerMutex.Unlock()

	return node, nil
}

// CreateNode creates non-validator nodes in the network.
func (n *LocalNetwork) CreateNode(config *driver.NodeConfig) (driver.Node, error) {
	return n.createNode(&node.OperaNodeConfig{
		Label:            config.Name,
		NetworkConfig:    &n.config,
		VmImplementation: n.config.VmImplementation,
	})
}

func (n *LocalNetwork) RemoveNode(node driver.Node) error {
	n.nodesMutex.Lock()
	id, err := node.GetNodeID()
	if err != nil {
		return fmt.Errorf("failed to get node id; %v", err)
	}

	delete(n.nodes, id)
	for _, other := range n.nodes {
		if err = other.RemovePeer(id); err != nil {
			n.nodesMutex.Unlock()
			return fmt.Errorf("failed to remove peer; %v", err)
		}
	}
	n.nodesMutex.Unlock()

	return nil
}

func (n *LocalNetwork) SendTransaction(tx *types.Transaction) {
	n.rpcWorkerPool.SendTransaction(tx)
}

func (n *LocalNetwork) DialRandomRpc() (rpc2.RpcClient, error) {
	nodes := n.GetActiveNodes()
	return nodes[rand.Intn(len(nodes))].DialRpc()
}

func (n *LocalNetwork) dialRandomValidatorRpc() (rpc2.RpcClient, error) {
	return n.validators[rand.Intn(len(n.validators))].DialRpc()
}

// treasureAccountPrivateKey is an account with tokens that can be used to
// initiate test applications and accounts.
const treasureAccountPrivateKey = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7" // Fakenet validator 1

const fakeNetworkID = 0xfa3

type localApplication struct {
	name       string
	controller *controller.AppController
	config     *driver.ApplicationConfig
	cancel     context.CancelFunc
	done       *sync.WaitGroup
}

func (a *localApplication) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	a.done.Add(1)
	go func() {
		defer a.done.Done()
		err := a.controller.Run(ctx)
		if err != nil {
			log.Printf("Failed to run load app: %v", err)
		}
	}()
	return nil
}

func (a *localApplication) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	a.cancel = nil
	log.Printf("waiting for application to stop: %s", a.name)
	a.done.Wait()
	log.Printf("application has stoped: %s", a.name)
	return nil
}

func (a *localApplication) Config() *driver.ApplicationConfig {
	return a.config
}

func (a *localApplication) GetNumberOfUsers() int {
	return a.controller.GetNumberOfUsers()
}

func (a *localApplication) GetSentTransactions(user int) (uint64, error) {
	return a.controller.GetTransactionsSentBy(user)
}

func (a *localApplication) GetReceivedTransactions() (uint64, error) {
	return a.controller.GetReceivedTransactions()
}

func (n *LocalNetwork) CreateApplication(config *driver.ApplicationConfig) (driver.Application, error) {
	rpcClient, err := n.dialRandomValidatorRpc()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC to initialize the application; %v", err)
	}
	defer rpcClient.Close()

	application, err := app.NewApplication(config.Type, rpcClient, n.primaryAccount, config.Users)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize on-chain app; %v", err)
	}

	sh, err := shaper.ParseRate(config.Rate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shaper; %v", err)
	}

	appController, err := controller.NewAppController(application, sh, config.Users, n)
	if err != nil {
		return nil, err
	}

	app := &localApplication{
		name:       config.Name,
		controller: appController,
		config:     config,
		done:       &sync.WaitGroup{},
	}

	n.appsMutex.Lock()
	n.apps = append(n.apps, app)
	n.appsMutex.Unlock()

	n.listenerMutex.Lock()
	for listener := range n.listeners {
		listener.AfterApplicationCreation(app)
	}
	n.listenerMutex.Unlock()

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

func (n *LocalNetwork) GetActiveApplications() []driver.Application {
	n.appsMutex.Lock()
	defer n.appsMutex.Unlock()
	return n.apps
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

	// Third, shut down the docker network.
	if n.network != nil {
		if err := n.network.Cleanup(); err != nil {
			errs = append(errs, err)
		}
	}

	errs = append(errs, n.rpcWorkerPool.Close())

	return errors.Join(errs...)
}

// GetDockerNetwork returns the underlying docker network.
func (n *LocalNetwork) GetDockerNetwork() *docker.Network {
	return n.network
}

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

package local

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/network/rpc"
	"github.com/Fantom-foundation/Norma/driver/node"
	rpcdriver "github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/Fantom-foundation/Norma/load/app"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
	"github.com/Fantom-foundation/go-opera/opera/contracts/sfc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	// TODO dynamically created validators are not added to this list. Can't be added without old epoch is sealed, because method randomValidatorRPC could return a non validator node.
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

	// nextAppId is the id to use for next created applications.
	nextAppId atomic.Uint32

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
	primaryAccount, err := app.NewAccount(0, treasureAccountPrivateKey, nil, fakeNetworkID)
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

// StartNode starts a node after it has been created.
func (n *LocalNetwork) StartNode(nd driver.Node) (driver.Node, error) {
	opera, ok := nd.(*node.OperaNode)
	if !ok {
		return nil, fmt.Errorf("trying to start non-sonic node")
	}
	return n.startNode(opera)
}

func (n *LocalNetwork) startNode(node *node.OperaNode) (*node.OperaNode, error) {
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

// createNode is an internal version of CreateNode enabling the creation
// of validator and non-validator nodes in the network.
func (n *LocalNetwork) createNode(nodeConfig *node.OperaNodeConfig) (*node.OperaNode, error) {
	node, err := node.StartOperaDockerNode(n.docker, n.network, nodeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start opera docker; %v", err)
	}
	return n.startNode(node)
}

// CreateNode creates nodes in the network during run.
func (n *LocalNetwork) CreateNode(config *driver.NodeConfig) (driver.Node, error) {
	newValId := 0
	if config.Validator {
		var err error
		newValId, err = n.registerValidatorNode(config)
		if err != nil {
			return nil, err
		}
	}

	if config.Cheater {
		_, err := n.createNode(&node.OperaNodeConfig{
			Label:            "cheater-" + config.Name,
			NetworkConfig:    &n.config,
			VmImplementation: n.config.VmImplementation,
			ValidatorId:      &newValId,
		})
		if err != nil {
			return nil, err
		}
	}

	return n.createNode(&node.OperaNodeConfig{
		Label:            config.Name,
		NetworkConfig:    &n.config,
		VmImplementation: n.config.VmImplementation,
		ValidatorId:      &newValId,
	})
}

// registerValidatorNode creates non genesis validator node in the network during run.
func (n *LocalNetwork) registerValidatorNode(config *driver.NodeConfig) (int, error) {
	newValId := 0
	//
	rpcClient, err := n.dialRandomGenesisValidatorRpc()
	if err != nil {
		return 0, err
	}
	defer rpcClient.Close()

	// get a representation of the deployed contract
	SFCContract, err := contract.NewSFC(sfc.ContractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get SFC contract representation; %v", err)
	}

	var lastValId *big.Int
	lastValId, err = SFCContract.LastValidatorID(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get validator count; %v", err)
	}

	newValId = int(lastValId.Int64()) + 1
	log.Printf("Creating validator node %s; id %d", config.Name, newValId)

	privateKeyECDSA := evmcore.FakeKey(uint32(newValId))
	validatorPubKey := validatorpk.PubKey{
		Raw:  crypto.FromECDSAPub(&privateKeyECDSA.PublicKey),
		Type: validatorpk.Types.Secp256k1,
	}

	validatorPrivateKeykHex := strings.TrimPrefix(hexutil.Encode(crypto.FromECDSA(privateKeyECDSA)), "0x")
	validatorAccount, err := app.NewAccount(newValId, validatorPrivateKeykHex, validatorPubKey.Bytes(), 0xfa3)
	if err != nil {
		return 0, fmt.Errorf("failed to create validator account: %v", err)
	}

	_, err = validatorAccount.CreateValidator(SFCContract, rpcClient)
	if err != nil {
		return 0, err
	}
	err = app.WaitUntilAccountNonceIs(crypto.PubkeyToAddress(privateKeyECDSA.PublicKey), 1, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("createValidator; failed to wait for nonce; %v", err)
	}

	lastValId, err = SFCContract.LastValidatorID(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get validator count; %v", err)
	}
	if newValId != int(lastValId.Int64()) {
		return 0, fmt.Errorf("failed to create validator %d", newValId)
	}

	return newValId, nil
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

	n.listenerMutex.Lock()
	for listener := range n.listeners {
		listener.AfterNodeRemoval(node)
	}
	n.listenerMutex.Unlock()

	return nil
}

func (n *LocalNetwork) KillNode(node driver.Node) error {
	return node.Kill()
}

func (n *LocalNetwork) SendTransaction(tx *types.Transaction) {
	n.rpcWorkerPool.SendTransaction(tx)
}

func (n *LocalNetwork) DialRandomRpc() (rpcdriver.RpcClient, error) {
	nodes := n.GetActiveNodes()
	return nodes[rand.Intn(len(nodes))].DialRpc()
}

// dialRandomGenesisValidatorRpc dials a random genesis validator node.
// When network starts and still doesn't have any traffic, then first transaction must come from validator node.
// Caused by: the regular nodes even when connected won't send transactions from their txpool,
// because they don't know whether they are on head or not if blockchain is empty.
func (n *LocalNetwork) dialRandomGenesisValidatorRpc() (rpcdriver.RpcClient, error) {
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
	rpcClient, err := n.dialRandomGenesisValidatorRpc()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC to initialize the application; %v", err)
	}
	defer rpcClient.Close()

	appId := n.nextAppId.Add(1)
	application, err := app.NewApplication(config.Type, rpcClient, n.primaryAccount, config.Users, 0, appId)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize on-chain app; %v", err)
	}

	sh, err := shaper.ParseRate(config.Rate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shaper; %v", err)
	}

	appController, err := controller.NewAppController(application, sh, config.Users, rpcClient, n)
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

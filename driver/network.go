package driver

import (
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockgen -source network.go -destination network_mock.go -package driver

// Network abstracts an execution environment for running scenarios.
// Implementations may run nodes and applications locally, in docker images, or
// remotely, on actual nodes. The interface is used by the scenario driver
// to execute scenario descriptions.
type Network interface {
	// CreateNode creates a new node instance running a network client based on
	// the given configuration. It is used by the scenario executor to add
	// nodes to the network as needed.
	CreateNode(config *NodeConfig) (Node, error)

	// RemoveNode removes node from the network
	RemoveNode(Node) error

	// CreateApplication creates a new application in this network, ready to
	// produce load as defined by its configuration.
	CreateApplication(config *ApplicationConfig) (Application, error)

	// GetActiveNodes obtains a list of active nodes in the network.
	GetActiveNodes() []Node

	// GetActiveApplications obtains a list of active apps in the network.
	GetActiveApplications() []Application

	// RegisterListener registers a listener to receive updates on network
	// changes, for instance, to update monitoring information. Registering
	// the same listener more than once will have no effect.
	RegisterListener(NetworkListener)

	// UnregisterListener removes the given listener from this network.
	UnregisterListener(NetworkListener)

	// Shutdown stops all applications and nodes in the network and frees
	// any potential other resources.
	Shutdown() error

	SendTransaction(tx *types.Transaction)

	DialRandomRpc() (app.RpcClient, error)
}

// NetworkConfig is a collection of network parameters to be used by factories
// creating network instances.
type NetworkConfig struct {
	// NumberOfValidators is the (static) number of validators in the network.
	NumberOfValidators int
	// The name of the StateDB implementation to be used by network nodes.
	StateDbImplementation string
}

// NetworkListener can be registered to networks to get callbacks whenever there
// are changes in the network.
type NetworkListener interface {
	// AfterNodeCreation is called whenever a new node has joined the network.
	AfterNodeCreation(Node)
	// AfterApplicationCreation is called after a new application has started.
	AfterApplicationCreation(Application)
}

type NodeConfig struct {
	Name string
	// TODO: add other parameters as needed
	//  - features to include on the node
	//  - state DB configuration
	//  - EVM configuration
}

type ApplicationConfig struct {
	Name string

	// Rate defines the Tx/s the source should produce while active.
	Rate float32

	// Accounts defines the amount of accounts sending transactions to the app.
	Accounts int

	// TODO: add other parameters as needed
	//  - application type
	//  - other traffic shapes
}

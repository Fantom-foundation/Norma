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

package driver

import (
	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/Fantom-foundation/Norma/driver/rpc"
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

	// RemoveNode ends the client gracefully and removes node from the network
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

	DialRandomRpc() (rpc.RpcClient, error)
}

// NetworkConfig is a collection of network parameters to be used by factories
// creating network instances.
type NetworkConfig struct {
	// NumberOfValidators is the (static) number of validators in the network.
	NumberOfValidators int
	// MaxBlockGas is the maximum gas limit for a block in the network.
	MaxBlockGas uint64
	// MaxEpochGas is the maximum gas limit for an epoch in the network.
	MaxEpochGas uint64
}

// NetworkListener can be registered to networks to get callbacks whenever there
// are changes in the network.
type NetworkListener interface {
	// AfterNodeCreation is called whenever a new node has joined the network.
	AfterNodeCreation(Node)
	// AfterNodeRemoval is called whenever a node is removed from the network.
	AfterNodeRemoval(Node)
	// AfterApplicationCreation is called after a new application has started.
	AfterApplicationCreation(Application)
}

type NodeConfig struct {
	Name         string
	Validator    bool
	Cheater      bool
	MountDatadir *string // mount node datadir to path if not nil
	MountGenesis *string // mount node genesis files to path if not nil
	// TODO: add other parameters as needed
	//  - features to include on the node
	//  - state DB configuration
	//  - EVM configuration
}

type ApplicationConfig struct {
	Name string

	// Type defines the on-chain app which should generate the traffic.
	Type string

	// Rate defines the Tx/s config the source should produce while active.
	Rate *parser.Rate

	// Users defines the number of users sending transactions to the app.
	Users int

	// TODO: add other parameters as needed
	//  - application type
}

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
	"io"

	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/Fantom-foundation/Norma/driver/rpc"
)

//go:generate mockgen -source node.go -destination node_mock.go -package driver

// Node is controlling a single node in a Norma network. It provides abstract
// control of a node, allowing it to be started (through an Environment),
// interact with the node, and shut it down.
type Node interface {
	// GetLabel returns a human-readable identifer for this node. Is is intended
	// to label data and should be unique within a single scenario run.
	GetLabel() string

	// Hostname returns the hostname of the host.
	Hostname() string

	// MetricsPort returns the port on which the node exposes its metrics.
	MetricsPort() int

	// IsRunning returns true if the node is still running, false if stopped.
	IsRunning() bool

	// GetNodeID returns an enode identifying this node within the Norma network.
	// An error shall be produced if no valid node ID could be obtained.
	GetNodeID() (NodeID, error)

	// GetServiceUrl returns the URL of a service running on the
	// represented node. May be nil if no such service is offered.
	GetServiceUrl(*network.ServiceDescription) *URL

	// DialRpc establish an RPC connection with the node and returns the RPC client.
	DialRpc() (rpc.RpcClient, error)

	// StreamLog provides a reader that is continuously providing the host log.
	// It is up to the caller to close the stream.
	StreamLog() (io.ReadCloser, error)

	// Stop shuts down this node gracefully, using its regular shutdown
	// procedure (not killed). After stopping the service, no more interactions
	// are expected to succeed.
	Stop() error

	// Kill shuts down this node disgracefully by using SigKill.
	Kill() error

	// Cleanup releases all underlying resources. After the cleanup no more
	// operations on this node are expected to succeed.
	Cleanup() error
}

// NodeID is a unique ID identifying each node. This identifier is used, for
// instance, to connect nodes within the network. In Opera, this ID is known
// as an 'enode' identifier.
type NodeID string

// URL is a mere alias type for a string supposed to encode a URL.
type URL string

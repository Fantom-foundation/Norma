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

package node

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	rpc2 "github.com/Fantom-foundation/Norma/driver/rpc"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/ethereum/go-ethereum/rpc"
)

var OperaRpcService = network.ServiceDescription{
	Name:     "OperaRPC",
	Port:     18545,
	Protocol: "http",
}

var OperaWsService = network.ServiceDescription{
	Name:     "OperaWs",
	Port:     18546,
	Protocol: "ws",
}

var OperaDebugService = network.ServiceDescription{
	Name:     "OperaPprof",
	Port:     6060,
	Protocol: "http",
}

var operaServices = network.ServiceGroup{}

func init() {
	if err := operaServices.RegisterService(&OperaRpcService); err != nil {
		panic(err)
	}
	if err := operaServices.RegisterService(&OperaWsService); err != nil {
		panic(err)
	}
	if err := operaServices.RegisterService(&OperaDebugService); err != nil {
		panic(err)
	}
}

const operaDockerImageName = "sonic"

// OperaNode implements the driver's Node interface by running a go-opera
// client on a generic host.
type OperaNode struct {
	host  network.Host
	label string
}

type OperaNodeConfig struct {
	// The label to be used to name this node. The label should not be empty.
	Label string
	// The ID of the validator, nil if the node should node be a validator.
	ValidatorId *int
	// The configuration of the network the configured node should be part of.
	NetworkConfig *driver.NetworkConfig
	// The EVM implementation to be used on this node.
	VmImplementation string
}

// labelPattern restricts labels for nodes to non-empty alpha-numerical strings
// with underscores and hyphens.
var labelPattern = regexp.MustCompile("[A-Za-z0-9_-]+")

// StartOperaDockerNode creates a new OperaNode running in a Docker container.
func StartOperaDockerNode(client *docker.Client, dn *docker.Network, config *OperaNodeConfig) (*OperaNode, error) {
	if !labelPattern.Match([]byte(config.Label)) {
		return nil, fmt.Errorf("invalid label for node: '%v'", config.Label)
	}

	shutdownTimeout := 1 * time.Second

	validatorId := "0"
	if config.ValidatorId != nil {
		validatorId = fmt.Sprintf("%d", *config.ValidatorId)
	}

	host, err := network.RetryReturn(network.DefaultRetryAttempts, 1*time.Second, func() (*docker.Container, error) {
		ports, err := network.GetFreePorts(len(operaServices.Services()))
		portForwarding := make(map[network.Port]network.Port, len(ports))
		for i, service := range operaServices.Services() {
			portForwarding[service.Port] = ports[i]
		}
		if err != nil {
			return nil, err
		}
		return client.Start(&docker.ContainerConfig{
			ImageName:       operaDockerImageName,
			ShutdownTimeout: &shutdownTimeout,
			PortForwarding:  portForwarding,
			Environment: map[string]string{
				"VALIDATOR_NUMBER": validatorId,
				"VALIDATORS_COUNT": fmt.Sprintf("%d", config.NetworkConfig.NumberOfValidators),
				"STATE_DB_IMPL":    config.NetworkConfig.StateDbImplementation,
				"VM_IMPL":          config.VmImplementation,
			},
			Network: dn,
		})
	})
	if err != nil {
		return nil, err
	}
	node := &OperaNode{
		host:  host,
		label: config.Label,
	}

	// Wait until the OperaNode inside the Container is ready.
	if err := network.Retry(network.DefaultRetryAttempts, 1*time.Second, func() error {
		_, err := node.GetNodeID()
		return err
	}); err == nil {
		return node, nil
	}

	// The node did not show up in time, so we consider the start to have failed.
	return nil, errors.Join(fmt.Errorf("failed to get node online"), node.host.Cleanup())
}

func (n *OperaNode) GetLabel() string {
	return n.label
}

// Hostname returns the hostname of the node.
// The hostname is accessible only inside the Docker network.
func (n *OperaNode) Hostname() string {
	return n.host.Hostname()
}

// MetricsPort returns the port on which the node exports its metrics.
// The port is accessible only inside the Docker network.
func (n *OperaNode) MetricsPort() int {
	return 6060
}

func (n *OperaNode) IsRunning() bool {
	return n.host.IsRunning()
}

func (n *OperaNode) GetServiceUrl(service *network.ServiceDescription) *driver.URL {
	addr := n.host.GetAddressForService(service)
	if addr == nil {
		return nil
	}
	url := driver.URL(fmt.Sprintf("%s://%s", service.Protocol, *addr))
	return &url
}

func (n *OperaNode) GetNodeID() (driver.NodeID, error) {
	url := n.GetServiceUrl(&OperaRpcService)
	if url == nil {
		return "", fmt.Errorf("node does not export an RPC server")
	}
	rpcClient, err := rpc.DialContext(context.Background(), string(*url))
	if err != nil {
		return "", err
	}
	var result struct {
		Enode string
	}
	err = rpcClient.Call(&result, "admin_nodeInfo")
	if err != nil {
		return "", err
	}
	return driver.NodeID(result.Enode), nil
}

func (n *OperaNode) StreamLog() (io.ReadCloser, error) {
	return n.host.StreamLog()
}

func (n *OperaNode) Stop() error {
	return n.host.Stop()
}

func (n *OperaNode) Cleanup() error {
	return n.host.Cleanup()
}

func (n *OperaNode) DialRpc() (rpc2.RpcClient, error) {
	url := n.GetServiceUrl(&OperaRpcService)
	if url == nil {
		return nil, fmt.Errorf("node %s does not export an RPC server", n.label)
	}

	rpcClient, err := network.RetryReturn(network.DefaultRetryAttempts, 1*time.Second, func() (*rpc.Client, error) {
		return rpc.DialContext(context.Background(), string(*url))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC for node %s; %v", n.label, err)
	}
	return rpc2.WrapRpcClient(rpcClient), nil
}

// AddPeer informs the client instance represented by the OperaNode about the
// existence of another node, to which it may establish a connection.
func (n *OperaNode) AddPeer(id driver.NodeID) error {
	rpcClient, err := n.DialRpc()
	if err != nil {
		return err
	}
	return network.Retry(network.DefaultRetryAttempts, 1*time.Second, func() error {
		return rpcClient.Call(nil, "admin_addPeer", id)
	})
}

// RemovePeer informs the client instance represented by the OperaNode
// that the input node is no more available in the network.
func (n *OperaNode) RemovePeer(id driver.NodeID) error {
	rpcClient, err := n.DialRpc()
	if err != nil {
		return err
	}
	return network.Retry(network.DefaultRetryAttempts, 1*time.Second, func() error {
		return rpcClient.Call(nil, "admin_removePeer", id)
	})
}

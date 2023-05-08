package node

import (
	"context"
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/ethereum/go-ethereum/rpc"
)

const OperaRPCPort = 18545

var OperaRpcService = network.ServiceDescription{
	Name: "OperaRPC",
	Port: OperaRPCPort,
}

const operaDockerImageName = "opera"

// OperaNode implements the driver's Node interface by running a go-opera
// client on a generic host.
type OperaNode struct {
	host network.Host
}

type OperaNodeConfig struct {
	// The ID of the validator, nil if the node should node be a validator.
	ValidatorId *int
	// The configuration of the network the configured node should be part of.
	NetworkConfig *driver.NetworkConfig
}

// StartOperaDockerNode creates a new OperaNode running in a Docker container.
func StartOperaDockerNode(client *docker.Client, config *OperaNodeConfig) (*OperaNode, error) {
	timeout := 1 * time.Second

	validatorId := "0"
	if config.ValidatorId != nil {
		validatorId = fmt.Sprintf("%d", *config.ValidatorId)
	}

	operaServicePort, err := network.GetFreePort()
	if err != nil {
		return nil, err
	}
	host, err := client.Start(&docker.ContainerConfig{
		ImageName:       operaDockerImageName,
		ShutdownTimeout: &timeout,
		PortForwarding: map[network.Port]network.Port{
			OperaRPCPort: operaServicePort,
		},
		Environment: map[string]string{
			"VALIDATOR_NUMBER": validatorId,
			"VALIDATORS_COUNT": fmt.Sprintf("%d", config.NetworkConfig.NumberOfValidators),
		},
	})
	if err != nil {
		return nil, err
	}
	node := &OperaNode{
		host: host,
	}

	// Wait until the OperaNode inside the Container is ready.
	for i := 0; i < 100; i++ {
		_, err := node.GetNodeID()
		if err == nil {
			return node, nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	// The node did not show up in time, so we consider the start to have failed.
	node.host.Cleanup()
	return nil, fmt.Errorf("failed to get node online")
}

func (n *OperaNode) GetRpcServiceUrl() *driver.URL {
	addr := n.host.GetAddressForService(&OperaRpcService)
	if addr == nil {
		return nil
	}
	url := driver.URL(fmt.Sprintf("http://%s", *addr))
	return &url
}

func (n *OperaNode) GetNodeID() (driver.NodeID, error) {
	url := n.GetRpcServiceUrl()
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

func (n *OperaNode) Stop() error {
	return n.host.Stop()
}

func (n *OperaNode) Cleanup() error {
	return n.host.Cleanup()
}

// AddPeer informs the client instance represented by the OperaNode about the
// existence of another node, to which it may establish a connection.
func (n *OperaNode) AddPeer(id driver.NodeID) error {
	url := n.GetRpcServiceUrl()
	if url == nil {
		return fmt.Errorf("node does not export an RPC server")
	}
	rpcClient, err := rpc.DialContext(context.Background(), string(*url))
	if err != nil {
		return err
	}
	return rpcClient.Call(nil, "admin_addPeer", id)
}

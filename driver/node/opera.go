package node

import (
	"context"
	"fmt"
	"io"
	"regexp"
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

const OperaWsPort = 18546

var OperaWsService = network.ServiceDescription{
	Name: "OperaWs",
	Port: OperaWsPort,
}

const OperaPprofPort = 6060

var OperaPprofService = network.ServiceDescription{
	Name: "OperaPprof",
	Port: OperaPprofPort,
}

const operaDockerImageName = "opera"

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
}

// labelPattern restricts labels for nodes to non-empty alpha-numerical strings
// with underscores and hyphens.
var labelPattern = regexp.MustCompile("[A-Za-z0-9_-]+")

// StartOperaDockerNode creates a new OperaNode running in a Docker container.
func StartOperaDockerNode(client *docker.Client, config *OperaNodeConfig) (*OperaNode, error) {
	if !labelPattern.Match([]byte(config.Label)) {
		return nil, fmt.Errorf("invalid label for node: '%v'", config.Label)
	}

	shutdownTimeout := 1 * time.Second

	validatorId := "0"
	if config.ValidatorId != nil {
		validatorId = fmt.Sprintf("%d", *config.ValidatorId)
	}

	ports, err := network.GetFreePorts(3)
	if err != nil {
		return nil, err
	}
	host, err := client.Start(&docker.ContainerConfig{
		ImageName:       operaDockerImageName,
		ShutdownTimeout: &shutdownTimeout,
		PortForwarding: map[network.Port]network.Port{
			OperaRPCPort:   ports[0],
			OperaWsPort:    ports[1],
			OperaPprofPort: ports[2],
		},
		Environment: map[string]string{
			"VALIDATOR_NUMBER": validatorId,
			"VALIDATORS_COUNT": fmt.Sprintf("%d", config.NetworkConfig.NumberOfValidators),
			"STATE_DB_IMPL":    config.NetworkConfig.StateDbImplementation,
		},
	})
	if err != nil {
		return nil, err
	}
	node := &OperaNode{
		host:  host,
		label: config.Label,
	}

	// Wait until the OperaNode inside the Container is ready. (3 minutes max)
	for i := 0; i < 3*60; i++ {
		_, err := node.GetNodeID()
		if err == nil {
			return node, nil
		}
		time.Sleep(time.Second)
	}

	// The node did not show up in time, so we consider the start to have failed.
	node.host.Cleanup()
	return nil, fmt.Errorf("failed to get node online")
}

func (n *OperaNode) GetLabel() string {
	return n.label
}

func (n *OperaNode) IsRunning() bool {
	return n.host.IsRunning()
}

func (n *OperaNode) GetHttpServiceUrl(service *network.ServiceDescription) *driver.URL {
	addr := n.host.GetAddressForService(service)
	if addr == nil {
		return nil
	}
	url := driver.URL(fmt.Sprintf("http://%s", *addr))
	return &url
}

func (n *OperaNode) GetWebsocketServiceUrl(service *network.ServiceDescription) *driver.URL {
	addr := n.host.GetAddressForService(service)
	if addr == nil {
		return nil
	}
	url := driver.URL(fmt.Sprintf("ws://%s", *addr))
	return &url
}

func (n *OperaNode) GetNodeID() (driver.NodeID, error) {
	url := n.GetHttpServiceUrl(&OperaRpcService)
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

// AddPeer informs the client instance represented by the OperaNode about the
// existence of another node, to which it may establish a connection.
func (n *OperaNode) AddPeer(id driver.NodeID) error {
	url := n.GetHttpServiceUrl(&OperaRpcService)
	if url == nil {
		return fmt.Errorf("node does not export an RPC server")
	}
	rpcClient, err := rpc.DialContext(context.Background(), string(*url))
	if err != nil {
		return err
	}
	return rpcClient.Call(nil, "admin_addPeer", id)
}

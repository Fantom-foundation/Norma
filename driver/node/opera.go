package node

import (
	"context"
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/ethereum/go-ethereum/rpc"
)

const OperaRPCPort = 18545

var OperaRpcService = driver.ServiceDescription{
	Name: "OperaPRC",
	Port: OperaRPCPort,
}

type OperaNode struct {
	*docker.Container
}

type PortManager interface {
	GetFreshPort() int
}

func StartOperaDockerNode(client *docker.Client, portManager PortManager, isValidator bool) (*OperaNode, error) {
	timeout := 1 * time.Second

	validatorFlag := "0"
	if isValidator {
		validatorFlag = "1"
	}

	host, err := client.Start(&docker.ContainerConfig{
		ImageName:       docker.OperaImageName,
		ShutdownTimeout: &timeout,
		PortForwarding: map[docker.Port]docker.Port{
			OperaRPCPort: docker.Port(portManager.GetFreshPort()),
		},
		Environment: map[string]string{
			"VALIDATOR_NUMBER": validatorFlag,
		},
	})
	if err != nil {
		return nil, err
	}
	node := &OperaNode{
		Container: host,
	}
	return node, nil
}

func (n *OperaNode) GetHost() driver.Host {
	return n.Container
}

func (n *OperaNode) GetNodeID() (driver.NodeID, error) {
	url := n.GetRpcServiceUrl()
	rpcClient, err := rpc.DialContext(context.Background(), url)
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

func (n *OperaNode) AddPeer(id driver.NodeID) error {
	url := n.GetRpcServiceUrl()
	rpcClient, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		return err
	}
	return rpcClient.Call(nil, "admin_addPeer", id)
}

func (n *OperaNode) GetRpcServiceUrl() string {
	return fmt.Sprintf("http://%s", n.GetAddressForService(&OperaRpcService))
}

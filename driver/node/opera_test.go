package node

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
)

func TestImplements(t *testing.T) {
	var inst OperaNode
	var _ driver.Node = &inst

}

func TestOperaNode_StartAndStop(t *testing.T) {
	docker, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	defer docker.Close()
	node, err := StartOperaDockerNode(docker, &OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	if err = node.host.Stop(); err != nil {
		t.Errorf("failed to stop Opera node: %v", err)
	}
}

func TestOperaNode_RpcServiceIsReadyAfterStartup(t *testing.T) {
	docker, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	defer docker.Close()
	node, err := StartOperaDockerNode(docker, &OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	if id, err := node.GetNodeID(); err != nil || len(id) == 0 {
		t.Errorf("failed to fetch NodeID from Opera node: '%v', err: %v", id, err)
	}
	if err = node.host.Stop(); err != nil {
		t.Errorf("failed to stop Opera node: %v", err)
	}
}

func TestOperaNode_StreamLog(t *testing.T) {
	docker, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	t.Cleanup(func() {
		_ = docker.Close()
	})

	node, err := StartOperaDockerNode(docker, &OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	t.Cleanup(func() {
		_ = node.Cleanup()
	})

	reader, err := node.StreamLog()
	if err != nil {
		t.Fatalf("cannot read logs: %e", err)
	}

	t.Cleanup(func() {
		_ = reader.Close()
	})

	done := make(chan bool)

	go func() {
		defer close(done)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Start/Switch to fullsync mode...") {
				done <- true
			}
		}
	}()

	var started bool
	select {
	case started = <-done:
	case <-time.After(10 * time.Second):
	}

	if !started {
		t.Errorf("expected log not found")
	}
}

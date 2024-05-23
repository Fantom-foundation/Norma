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
	"bufio"
	"fmt"
	"io"
	"net/http"
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
	t.Cleanup(func() {
		_ = docker.Close()
	})
	node, err := StartOperaDockerNode(docker, nil, &OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	t.Cleanup(func() {
		_ = node.Cleanup()
	})

	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	t.Cleanup(func() {
		_ = node.Cleanup()
	})
	if err = node.host.Stop(); err != nil {
		t.Errorf("failed to stop Opera node: %v", err)
	}
}

func TestOperaNode_RpcServiceIsReadyAfterStartup(t *testing.T) {
	docker, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	t.Cleanup(func() {
		_ = docker.Close()
	})
	node, err := StartOperaDockerNode(docker, nil, &OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	t.Cleanup(func() {
		_ = node.Cleanup()
	})

	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	t.Cleanup(func() {
		_ = node.Cleanup()
	})
	if id, err := node.GetNodeID(); err != nil || len(id) == 0 {
		t.Errorf("failed to fetch NodeID from Opera node: '%v', err: %v", id, err)
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

	node, err := StartOperaDockerNode(docker, nil, &OperaNodeConfig{
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
			if strings.Contains(line, "WebSocket enabled") {
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

func TestOperaNode_MetricsExposed(t *testing.T) {
	docker, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	t.Cleanup(func() {
		_ = docker.Close()
	})

	node, err := StartOperaDockerNode(docker, nil, &OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	t.Cleanup(func() {
		_ = node.Cleanup()
	})

	url := node.GetServiceUrl(&OperaDebugService)

	var apiWorks bool
	for i := 0; i < 100; i++ {
		resp, err := http.Get(fmt.Sprintf("%s/debug/metrics/prometheus", string(*url)))
		if err == nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err == nil && strings.Contains(string(bodyBytes), "# TYPE") {
				apiWorks = true
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !apiWorks {
		t.Errorf("monitoring API has not been available")
	}
}

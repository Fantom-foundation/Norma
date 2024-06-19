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

package prometheusmon

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver/network"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
)

func TestPrometheusCanBeRun(t *testing.T) {
	t.Parallel()
	net := createLocalNetwork(t)
	prom := startPrometheus(t, net)

	// test prometheus is running
	resp, err := http.Get(prom.GetUrl() + "/-/ready")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("prometheus is not running")
	}
}

func TestPrometheusCanBeShutdown(t *testing.T) {
	t.Parallel()
	net := createLocalNetwork(t)
	prom := startPrometheus(t, net)

	// shutdown prometheus
	err := prom.Shutdown()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// test prometheus is not running
	_, err = http.Get(prom.GetUrl() + "/-/ready")
	if err == nil {
		t.Errorf("prometheus is still running")
	}
}

func TestNodeCanBeAdded(t *testing.T) {
	t.Parallel()
	net := createLocalNetwork(t)
	prom := startPrometheus(t, net)

	// get a node
	nodes := net.GetActiveNodes()
	if len(nodes) == 0 {
		t.Fatalf("no active nodes")
	}
	node := nodes[0]

	// add node
	if err := prom.AddNode(node); err != nil {
		t.Fatalf("error: %v", err)
	}
	// wait for prometheus to reload config
	err := network.Retry(network.DefaultRetryAttempts, 1*time.Second, func() error {
		// verify node is added by calling prometheus API
		resp, err := http.Get(prom.GetUrl() + "/api/v1/targets")
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not HTTP OK: %d", resp.StatusCode)
		}

		// check response contains the node's label
		rawResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()

		if bytes.Contains(rawResponse, []byte(node.GetLabel())) {
			return fmt.Errorf("response does not contain Node Label: %s", rawResponse)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("node '%s' not added into prometheus: %s", node.GetLabel(), err)
	}
}

// startPrometheus starts a prometheus node and returns it.
func startPrometheus(t *testing.T, net *local.LocalNetwork) *Prometheus {
	prom, err := Start(net, net.GetDockerNetwork())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Cleanup(func() {
		_ = prom.Shutdown()
	})
	return prom
}

// createLocalNetwork creates a docker network and returns it.
func createLocalNetwork(t *testing.T) *local.LocalNetwork {
	config := driver.NetworkConfig{NumberOfValidators: 1}
	net, err := local.NewLocalNetwork(&config)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Cleanup(func() {
		_ = net.Shutdown()
	})
	return net
}

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

package local

import (
	"fmt"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"go.uber.org/mock/gomock"
)

func TestLocalNetworkIsNetwork(t *testing.T) {
	var net LocalNetwork
	var _ driver.Network = &net
}

func TestLocalNetwork_CanStartNodesAndShutThemDown(t *testing.T) {
	t.Parallel()
	config := driver.NetworkConfig{NumberOfValidators: 1}
	for _, N := range []int{1, 3} {
		N := N
		t.Run(fmt.Sprintf("num_nodes=%d", N), func(t *testing.T) {
			t.Parallel()
			net, err := NewLocalNetwork(&config)
			if err != nil {
				t.Fatalf("failed to create new local network: %v", err)
			}
			t.Cleanup(func() {
				_ = net.Shutdown()
			})

			nodes := []driver.Node{}
			for i := 0; i < N; i++ {
				node, err := net.CreateNode(&driver.NodeConfig{
					Name: fmt.Sprintf("T-%d", i),
				})
				if err != nil {
					t.Errorf("failed to create node: %v", err)
					continue
				}
				nodes = append(nodes, node)
			}

			for _, node := range nodes {
				if err := node.Stop(); err != nil {
					t.Errorf("failed to stop node: %v", err)
				}
			}

			for _, node := range nodes {
				if err := node.Cleanup(); err != nil {
					t.Errorf("failed to cleanup node: %v", err)
				}
			}
		})
	}
}

func TestLocalNetwork_CanStartApplicatonsAndShutThemDown(t *testing.T) {
	t.Parallel()
	config := driver.NetworkConfig{NumberOfValidators: 1}
	for _, N := range []int{1, 3} {
		N := N
		t.Run(fmt.Sprintf("num_nodes=%d", N), func(t *testing.T) {
			t.Parallel()

			net, err := NewLocalNetwork(&config)
			if err != nil {
				t.Fatalf("failed to create new local network: %v", err)
			}
			t.Cleanup(func() {
				_ = net.Shutdown()
			})

			apps := []driver.Application{}
			for i := 0; i < N; i++ {
				app, err := net.CreateApplication(&driver.ApplicationConfig{
					Name: fmt.Sprintf("T-%d", i),
				})
				if err != nil {
					t.Errorf("failed to create app: %v", err)
					continue
				}

				if got, want := app.Config().Name, fmt.Sprintf("T-%d", i); got != want {
					t.Errorf("app configurion not propagated: %v != %v", got, want)
				}

				apps = append(apps, app)
			}

			for _, app := range apps {
				if err := app.Start(); err != nil {
					t.Errorf("failed to start app: %v", err)
				}
			}

			for _, app := range apps {
				if err := app.Stop(); err != nil {
					t.Errorf("failed to stop app: %v", err)
				}
			}
		})
	}
}

func TestLocalNetwork_CanPerformNetworkShutdown(t *testing.T) {
	t.Parallel()
	N := 2
	config := driver.NetworkConfig{NumberOfValidators: 1}

	net, err := NewLocalNetwork(&config)
	if err != nil {
		t.Fatalf("failed to create new local network: %v", err)
	}
	t.Cleanup(func() {
		_ = net.Shutdown()
	})

	for i := 0; i < N; i++ {
		_, err := net.CreateNode(&driver.NodeConfig{
			Name: fmt.Sprintf("T-%d", i),
		})
		if err != nil {
			t.Errorf("failed to create node: %v", err)
		}
	}

	for i := 0; i < N; i++ {
		_, err := net.CreateApplication(&driver.ApplicationConfig{
			Name: fmt.Sprintf("T-%d", i),
		})
		if err != nil {
			t.Errorf("failed to create app: %v", err)
		}
	}

	if err := net.Shutdown(); err != nil {
		t.Errorf("failed to shut down network: %v", err)
	}
}

func TestLocalNetwork_CanRunWithMultipleValidators(t *testing.T) {
	t.Parallel()
	for _, N := range []int{1, 3} {
		N := N
		config := driver.NetworkConfig{NumberOfValidators: N}
		t.Run(fmt.Sprintf("num_validators=%d", N), func(t *testing.T) {
			t.Parallel()
			net, err := NewLocalNetwork(&config)
			if err != nil {
				t.Fatalf("failed to create new local network: %v", err)
			}
			t.Cleanup(func() {
				_ = net.Shutdown()
			})

			app, err := net.CreateApplication(&driver.ApplicationConfig{
				Name: "TestApp",
			})
			if err != nil {
				t.Fatalf("failed to create app: %v", err)
			}

			if err := app.Start(); err != nil {
				t.Errorf("failed to start app: %v", err)
			}

			if err := app.Stop(); err != nil {
				t.Errorf("failed to stop app: %v", err)
			}
		})
	}
}

func TestLocalNetwork_NotifiesListenersOnNodeStartup(t *testing.T) {
	t.Parallel()
	config := driver.NetworkConfig{NumberOfValidators: 2}
	ctrl := gomock.NewController(t)
	listener := driver.NewMockNetworkListener(ctrl)

	net, err := NewLocalNetwork(&config)
	if err != nil {
		t.Fatalf("failed to create new local network: %v", err)
	}
	t.Cleanup(func() {
		_ = net.Shutdown()
	})

	activeNodes := net.GetActiveNodes()
	if got, want := len(activeNodes), config.NumberOfValidators; got != want {
		t.Errorf("invalid number of active nodes, got %d, want %d", got, want)
	}

	net.RegisterListener(listener)
	listener.EXPECT().AfterNodeCreation(gomock.Any())

	net.CreateNode(&driver.NodeConfig{
		Name: "Test",
	})

	activeNodes = net.GetActiveNodes()
	if got, want := len(activeNodes), config.NumberOfValidators+1; got != want {
		t.Errorf("invalid number of active nodes, got %d, want %d", got, want)
	}

}

func TestLocalNetwork_NotifiesListenersOnAppStartup(t *testing.T) {
	t.Parallel()
	config := driver.NetworkConfig{NumberOfValidators: 1}
	ctrl := gomock.NewController(t)
	listener := driver.NewMockNetworkListener(ctrl)

	net, err := NewLocalNetwork(&config)
	if err != nil {
		t.Fatalf("failed to create new local network: %v", err)
	}
	t.Cleanup(func() {
		_ = net.Shutdown()
	})

	net.RegisterListener(listener)
	listener.EXPECT().AfterApplicationCreation(gomock.Any())

	_, err = net.CreateApplication(&driver.ApplicationConfig{
		Name: "TestApp",
	})
	if err != nil {
		t.Errorf("creation of app failed: %v", err)
	}
}

func TestLocalNetwork_CanRemoveNode(t *testing.T) {
	t.Parallel()
	config := driver.NetworkConfig{NumberOfValidators: 1}
	for _, N := range []int{1, 3} {
		N := N
		t.Run(fmt.Sprintf("num_nodes=%d", N), func(t *testing.T) {
			t.Parallel()
			net, err := NewLocalNetwork(&config)
			ctrl := gomock.NewController(t)
			listener := driver.NewMockNetworkListener(ctrl)
			listener.EXPECT().AfterNodeCreation(gomock.Any()).Times(N)
			listener.EXPECT().AfterNodeRemoval(gomock.Any()).Times(N)
			net.RegisterListener(listener)

			if err != nil {
				t.Fatalf("failed to create new local network: %v", err)
			}
			t.Cleanup(func() {
				_ = net.Shutdown()
			})

			nodes := make([]driver.Node, 0, N)
			for i := 0; i < N; i++ {
				node, err := net.CreateNode(&driver.NodeConfig{
					Name: fmt.Sprintf("T-%d", i),
				})
				if err != nil {
					t.Errorf("failed to create node: %s", err)
				}
				nodes = append(nodes, node)
			}

			// remove nodes one by one
			for _, node := range nodes {
				if err := net.RemoveNode(node); err != nil {
					t.Errorf("cannot remove node: %s", err)
				}

				id, err := node.GetNodeID()
				if err != nil {
					t.Errorf("cannot get node ID: %s", err)
				}

				_, exists := net.nodes[id]
				if exists {
					t.Errorf("node %s was not removed", id)
				}
			}

			// removed nodes are only detached from the network, but still running - i.e. they can be turned off
			for _, node := range nodes {
				if err := node.Stop(); err != nil {
					t.Errorf("failed to stop node: %v", err)
				}
				if err := node.Cleanup(); err != nil {
					t.Errorf("failed to cleanup node: %v", err)
				}
			}
		})
	}
}

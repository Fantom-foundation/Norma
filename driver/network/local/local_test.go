package local

import (
	"fmt"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
)

func TestLocalNetworkIsNetwork(t *testing.T) {
	var net LocalNetwork
	var _ driver.Network = &net
}

func TestLocalNetwork_CanStartNodesAndShutThemDown(t *testing.T) {
	for N := 1; N <= 3; N++ {
		t.Run(fmt.Sprintf("num_nodes=%d", N), func(t *testing.T) {

			net, err := NewLocalNetwork()
			if err != nil {
				t.Fatalf("failed to create new local network: %v", err)
			}

			nodes := []driver.Node{}
			for i := 0; i < N; i++ {
				node, err := net.CreateNode(&driver.NodeConfig{
					Name: fmt.Sprintf("T-%d", i),
				})
				if err != nil {
					t.Errorf("failed to create node: %v", err)
				}
				defer node.Cleanup()
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
	for N := 1; N <= 3; N++ {
		t.Run(fmt.Sprintf("num_nodes=%d", N), func(t *testing.T) {

			net, err := NewLocalNetwork()
			if err != nil {
				t.Fatalf("failed to create new local network: %v", err)
			}

			// We need at last one node running.
			node, err := net.CreateNode(&driver.NodeConfig{})
			if err != nil {
				t.Fatalf("failed to create node in network: %v", err)
			}
			defer node.Cleanup()

			apps := []driver.Application{}
			for i := 0; i < N; i++ {
				app, err := net.CreateApplication(&driver.ApplicationConfig{
					Name: fmt.Sprintf("T-%d", i),
				})
				if err != nil {
					t.Errorf("failed to create app: %v", err)
				}
				defer app.Stop()
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
	N := 2

	net, err := NewLocalNetwork()
	if err != nil {
		t.Fatalf("failed to create new local network: %v", err)
	}

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

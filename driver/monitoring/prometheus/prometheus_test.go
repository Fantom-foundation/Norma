package prometheusmon

import (
	"net/http"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/node"
)

func TestPrometheusCanBeRun(t *testing.T) {
	prom := startPrometheus(t, nil, nil)
	// test prometheus is running
	resp, err := http.Get(prom.GetUrl() + "/-/ready")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("prometheus is not running")
	}
}

func TestNodeCanBeAdded(t *testing.T) {
	dn := createDockerNetwork(t)
	prom := startPrometheus(t, nil, dn)
	opera := startOperaNode(t, dn)
	// test node is added
	if err := prom.AddNode(opera); err != nil {
		t.Fatalf("error: %v", err)
	}

	// TODO: CHECK NODE WAS ADDED
}

// startPrometheus starts a prometheus node and returns it.
func startPrometheus(t *testing.T, net driver.Network, dn *docker.Network) PrometheusNode {
	prometheus := PrometheusDocker{}
	node, err := prometheus.Start(net, dn)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Cleanup(func() {
		_ = node.Shutdown()
	})
	return node
}

// startOperaNode starts a opera node and returns it.
func startOperaNode(t *testing.T, dn *docker.Network) driver.Node {
	client, err := docker.NewClient()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	opera, err := node.StartOperaDockerNode(client, dn, &node.OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Cleanup(func() {
		_ = opera.Cleanup()
		_ = client.Close()
	})
	return opera
}

// createDockerNetwork creates a docker network and returns it.
func createDockerNetwork(t *testing.T) *docker.Network {
	client, err := docker.NewClient()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	dn, err := client.CreateBridgeNetwork()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Cleanup(func() {
		_ = dn.Cleanup()
		_ = client.Close()
	})
	return dn
}

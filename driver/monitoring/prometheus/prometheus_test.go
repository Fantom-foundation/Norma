package prometheusmon

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

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

func TestPrometheusCanBeShutdown(t *testing.T) {
	prom := startPrometheus(t, nil, nil)
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
	dn := createDockerNetwork(t)
	prom := startPrometheus(t, nil, dn)
	opera := startOperaNode(t, dn)
	// add node
	if err := prom.AddNode(opera); err != nil {
		t.Fatalf("error: %v", err)
	}
	// wait for prometheus to reload config
	time.Sleep(5 * time.Second)
	// verify node is added by calling prometheus API
	resp, err := http.Get(prom.GetUrl() + "/api/v1/targets")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	// check response contains the node's label
	rawResponse, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !bytes.Contains(rawResponse, []byte(opera.GetLabel())) {
		t.Fatalf("expected response to contain %s", opera.GetLabel())
	}
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

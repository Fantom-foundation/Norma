package prometheusmon

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

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
	loaded := false
	for i := 0; i < 50; i++ {
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
		_ = resp.Body.Close()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if bytes.Contains(rawResponse, []byte(node.GetLabel())) {
			loaded = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !loaded {
		t.Fatalf("node '%s' not added into prometheus", node.GetLabel())
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

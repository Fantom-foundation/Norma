package prometheusmon

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/network"
)

// PrometheusPort is the default port for the PrometheusDocker service.
const PrometheusPort = 9090

// prometheusImage is the default Docker image for the PrometheusDocker service.
const prometheusImage = "prom/prometheus:v2.44.0"

// PrometheusRunner is the interface for starting a prometheus instance.
type PrometheusRunner interface {
	Start(net driver.Network, dn *docker.Network) (Prometheus, error)
}

// Prometheus is the interface for a prometheus instance.
type Prometheus interface {
	AddNode(node driver.Node) error
	Shutdown() error
}

// PrometheusDocker is a Prometheus instance running in a Docker container.
type PrometheusDocker struct {
	container *docker.Container
	net       driver.Network
	nodesLock sync.Mutex
	exited    bool
}

// PrometheusDockerRunner is a PrometheusRunner that starts a Prometheus instance in a Docker container.
type PrometheusDockerRunner struct{}

// Start starts a Prometheus instance in a Docker container.
func (p *PrometheusDockerRunner) Start(net driver.Network, dn *docker.Network) (Prometheus, error) {
	timeout := 1 * time.Second

	client, err := docker.NewClient()
	if err != nil {
		return nil, err
	}

	ports, err := network.GetFreePorts(1)
	if err != nil {
		return nil, err
	}

	// start the container
	container, err := client.Start(&docker.ContainerConfig{
		ImageName:       prometheusImage,
		ShutdownTimeout: &timeout,
		PortForwarding: map[network.Port]network.Port{
			PrometheusPort: ports[0],
		},
		Network: dn,
	})
	if err != nil {
		return nil, err
	}

	prometheus := &PrometheusDocker{
		container: container,
		net:       net,
	}

	// initialize the config
	err = prometheus.initializeConfig()
	if err != nil {
		_ = container.Cleanup()
		return nil, err
	}

	// wait until the prometheus inside the Container is ready. (15 seconds max)
	// this is necessary for SIGHUP signal to be delivered correctly
	url := fmt.Sprintf("http://localhost:%d", ports[0])
	for i := 0; i < 15; i++ {
		// send get request to url
		resp, err := http.Get(url + "/-/ready")
		if err == nil {
			// check response status
			if resp.StatusCode != http.StatusOK {
				continue
			}
			// check response contains "Ready"
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			if !strings.Contains(string(b), "Ready") {
				continue
			}

			log.Printf("started PrometheusDocker on %s", url)

			// listen for new Nodes
			net.RegisterListener(prometheus)

			// get nodes that have been started before this instance creation
			for _, node := range prometheus.net.GetActiveNodes() {
				prometheus.AfterNodeCreation(node)
			}

			return prometheus, nil
		}
		time.Sleep(time.Second)
	}

	// if we reach this point, the prometheus instance is not ready
	_ = container.Cleanup()
	return nil, fmt.Errorf("prometheus instance is not ready")
}

// AddNode adds a new target to the PrometheusDocker configuration to be observed.
func (p *PrometheusDocker) AddNode(node driver.Node) error {
	cfg := fmt.Sprintf("%s", fmt.Sprintf(promTargetCfgTmpl, node.Hostname(), node.MetricsPort(), node.GetLabel()))
	_, err := p.container.Exec(
		[]string{"sh", "-c", fmt.Sprintf("echo '%s' > /etc/prometheus/opera-%s.json", cfg, node.Hostname())})
	if err != nil {
		return err
	}
	// we also need to reload the config
	return p.reloadConfig()
}

// Shutdown shuts down the PrometheusDocker instance.
func (p *PrometheusDocker) Shutdown() error {
	if p.exited {
		return nil
	}
	p.exited = true
	p.net.UnregisterListener(p)
	return p.container.Cleanup()
}

func (p *PrometheusDocker) AfterNodeCreation(node driver.Node) {
	p.nodesLock.Lock()
	defer p.nodesLock.Unlock()
	if err := p.AddNode(node); err != nil {
		log.Printf("failed to add node %s to PrometheusDocker: %s", node.Hostname(), err)
	}
}

func (p *PrometheusDocker) AfterApplicationCreation(driver.Application) {
	// ignored
}

// initializeConfig initializes the PrometheusDocker configuration file by echoing config content
// into container's config location.
func (p *PrometheusDocker) initializeConfig() error {
	_, err := p.container.Exec(
		[]string{"sh", "-c", fmt.Sprintf("echo '%s' > /etc/prometheus/prometheus.yml", promCfg)})
	return err
}

// reloadConfig reloads the PrometheusDocker configuration by sending "SIGHUP" signal.
func (p *PrometheusDocker) reloadConfig() error {
	return p.container.SendSignal("SIGHUP")
}

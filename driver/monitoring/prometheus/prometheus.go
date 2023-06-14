package prometheusmon

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	"github.com/Fantom-foundation/Norma/driver/network"
)

// PrometheusPort is the default port for the PrometheusDockerNode service.
const PrometheusPort = 9090

// prometheusImage is the default Docker image for the PrometheusDockerNode service.
const prometheusImage = "prom/prometheus:v2.44.0"

// Prometheus is the interface for starting a prometheus instance.
type Prometheus interface {
	Start(net driver.Network, dn *docker.Network) (PrometheusNode, error)
}

// PrometheusNode is the interface for a prometheus instance.
type PrometheusNode interface {
	// AddNode adds a new target to the PrometheusNode configuration to be observed.
	AddNode(node driver.Node) error
	// Shutdown shuts down the PrometheusNode instance.
	Shutdown() error
	// GetUrl returns the URL of the PrometheusNode instance.
	GetUrl() string
}

// PrometheusDockerNode is a PrometheusNode instance running in a Docker container.
type PrometheusDockerNode struct {
	container *docker.Container
	port      network.Port
	net       driver.Network
}

// PrometheusDocker is a Prometheus that starts a PrometheusNode instance in a Docker container.
type PrometheusDocker struct{}

// Start starts a PrometheusNode instance in a Docker container.
func (p *PrometheusDocker) Start(net driver.Network, dn *docker.Network) (PrometheusNode, error) {
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

	prometheus := &PrometheusDockerNode{
		container: container,
		net:       net,
		port:      ports[0],
	}

	// initialize the config
	err = prometheus.initializeConfig()
	if err != nil {
		_ = container.Cleanup()
		return nil, err
	}

	// wait until the prometheus inside the Container is ready. (15 seconds max)
	// this is necessary for SIGHUP signal to be delivered correctly
	for i := 0; i < 15; i++ {
		// send get request to `<url>/-/ready` which contains status
		resp, err := http.Get(prometheus.GetUrl() + "/-/ready")
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

			log.Printf("started Prometheus on %s", prometheus.GetUrl())

			// register the prometheus instance as a listener
			if net != nil {
				// listen for new Nodes
				net.RegisterListener(prometheus)

				// get nodes that have been started before this instance creation
				for _, node := range prometheus.net.GetActiveNodes() {
					prometheus.AfterNodeCreation(node)
				}
			}

			return prometheus, nil
		}
		time.Sleep(time.Second)
	}

	// if we reach this point, the prometheus instance is not ready
	_ = container.Cleanup()
	return nil, fmt.Errorf("prometheus instance is not ready")
}

// AddNode adds a new target to the PrometheusDockerNode configuration to be observed.
func (p *PrometheusDockerNode) AddNode(node driver.Node) error {
	cfg, err := renderConfigForNode(node)
	if err != nil {
		return err
	}
	_, err = p.container.Exec(
		[]string{"sh", "-c", fmt.Sprintf("echo '%s' > /etc/prometheus/opera-%s.json", cfg, node.Hostname())})
	if err != nil {
		return err
	}
	// we also need to reload the config
	return p.reloadConfig()
}

// Shutdown shuts down the PrometheusDockerNode instance.
func (p *PrometheusDockerNode) Shutdown() error {
	if p.net != nil {
		p.net.UnregisterListener(p)
	}
	return p.container.Cleanup()
}

// GetUrl returns the URL of the PrometheusDockerNode instance.
func (p *PrometheusDockerNode) GetUrl() string {
	return fmt.Sprintf("http://localhost:%d", p.port)
}

func (p *PrometheusDockerNode) AfterNodeCreation(node driver.Node) {
	if err := p.AddNode(node); err != nil {
		log.Printf("failed to add node %s to PrometheusDockerNode: %s", node.Hostname(), err)
	}
}

func (p *PrometheusDockerNode) AfterApplicationCreation(driver.Application) {
	// ignored
}

// initializeConfig initializes the PrometheusDockerNode configuration file by echoing config content
// into container's config location.
func (p *PrometheusDockerNode) initializeConfig() error {
	_, err := p.container.Exec(
		[]string{"sh", "-c", fmt.Sprintf("echo '%s' > /etc/prometheus/prometheus.yml", promCfg)})
	return err
}

// reloadConfig reloads the PrometheusDockerNode configuration by sending "SIGHUP" signal.
func (p *PrometheusDockerNode) reloadConfig() error {
	return p.container.SendSignal(docker.SigHup)
}

package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Client provides means to spawn Docker containers capable of hosting
// services like the go-opera client.
type Client struct {
	cli *client.Client
}

// Container represents a Docker Container, typically used for running a
// Fantom network Node, thus an instance of the go-opera client.
// *Container implements the driver.Host interface.
type Container struct {
	id      string
	client  *Client
	config  *ContainerConfig
	stopped bool
	cleaned bool
}

// ContainerConfig defines parameters for running Docker Containers.
type ContainerConfig struct {
	ImageName       string
	ShutdownTimeout *time.Duration
	PortForwarding  map[network.Port]network.Port // Container Port => Host Port
	Environment     map[string]string
}

// NewClient creates a new client facilitating the creation of Docker
// Containers capable of hosting services. Clients successfully created
// through this function should be Closed() eventually.
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Client{cli}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

// Start creates and runs one Container. The provided configuration allows
// to configure the Docker image to run inside the container -- and thus the
// services to be offered -- and port-forwarding specifications to make those
// services reachable from outside the Docker container (e.g. by the
// application running this code).
func (c *Client) Start(config *ContainerConfig) (*Container, error) {

	envVars := []string{}
	for key, value := range config.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	portMapping := nat.PortMap{}
	for inner, outer := range config.PortForwarding {
		portMapping[nat.Port(fmt.Sprintf("%d/tcp", inner))] = []nat.PortBinding{{
			HostIP:   "localhost",
			HostPort: fmt.Sprintf("%d/tcp", outer),
		}}
	}

	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Image: config.ImageName,
		Tty:   false,
		Env:   envVars,
	}, &container.HostConfig{
		PortBindings: portMapping,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return &Container{resp.ID, c, config, false, false}, nil
}

// IsRunning returns true if the Container has not been stopped yet and is
// expected to offer its services.
func (c *Container) IsRunning() bool {
	return !c.stopped
}

// Stop terminates this container. Services within the container will be
// signaled about the upcoming termination followed by being killed after a set
// timeout (see ContainerConfig.ShutdownTimeout).
func (c *Container) Stop() error {
	if c.stopped {
		return nil
	}
	c.stopped = true
	return c.client.cli.ContainerStop(context.Background(), c.id, c.config.ShutdownTimeout)
}

// Cleanup stops the container (unless it is already stopped) and frees any
// resources associated to it. After the operation, the Container is to be
// considered invalid.
func (c *Container) Cleanup() error {
	if c.cleaned {
		return nil
	}
	if err := c.Stop(); err != nil {
		return err
	}
	c.cleaned = true
	return c.client.cli.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{})
}

// GetAddressForService retrieves the Address of a service running in this
// Container and being exported to the Docker's host environment. If there is
// no such service (e.g., because it was not marked as to be exported during
// the Start of the Container), nil will be returned.
func (n *Container) GetAddressForService(service *network.ServiceDescription) *network.AddressPort {
	// All services inside the container are reached through port-forwarding
	// on the localhost. Non-forwarded services are not supported.
	port, ok := n.config.PortForwarding[service.Port]
	if !ok {
		return nil
	}
	res := network.AddressPort(fmt.Sprintf("%s:%d", "localhost", port))
	return &res
}

// SaveLogTo fetches the log of the container and saves it to the given directory.
func (c *Container) SaveLogTo(directory string) error {
	opt := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}

	// TODO if this proves insufficient, an alternative would be to mount certain directories from
	// the container to temp on the host and here just copy local directories
	reader, err := c.client.cli.ContainerLogs(context.Background(), c.id, opt)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s_%s.log", directory, c.config.ImageName, c.id))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listContainers() ([]types.Container, error) {
	return c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

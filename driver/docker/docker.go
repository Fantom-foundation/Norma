package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const OperaImageName = "opera"

type Port uint16

// Container represents a running instance of the client application such as go-opera
type Container struct {
	id      string
	client  *Client
	config  *ContainerConfig
	running bool
}

// Client is an initialized application that can run containers
type Client struct {
	cli *client.Client
}

// ClientConfig configures common parameters for running containers.
type ContainerConfig struct {
	ImageName       string
	ShutdownTimeout *time.Duration
	PortForwarding  map[Port]Port // Inner Port => public Port
	Environment     map[string]string
}

// NewClient creates the docker environment
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

// Start creates and runs one Container
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

	return &Container{resp.ID, c, config, true}, nil
}

func (c *Container) IsRunning() bool {
	return c.running
}

// Stop terminates the running container.
func (c *Container) Stop() error {
	c.running = false
	return c.client.cli.ContainerStop(context.Background(), c.id, c.config.ShutdownTimeout)
}

func (c *Container) Cleanup() error {
	if err := c.Stop(); err != nil {
		return err
	}
	return c.client.cli.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{})
}

func (c *Container) GetIP() driver.IP {
	return "localhost"
}

func (n *Container) GetAddressForService(service *driver.ServiceDescription) driver.AddressPort {
	port, ok := n.config.PortForwarding[Port(service.Port)]
	if !ok {
		return ""
	}
	return driver.AddressPort(fmt.Sprintf("%s:%d", n.GetIP(), port))
}

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

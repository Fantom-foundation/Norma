package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"time"
)

const imageName = "opera"

// Container represents a running instance of the client application such as go-opera
type Container struct {
	id     string
	client *Client
	config *ClientConfig
}

// Client is an initialized application that can run containers
type Client struct {
	cli *client.Client
}

// ClientConfig configures common parameters for running containers.
type ClientConfig struct {
	imageName       string
	shutdownTimeout *time.Duration
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
func (c *Client) Start(config *ClientConfig) (*Container, error) {
	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Image: config.imageName,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return &Container{resp.ID, c, config}, nil
}

// Stop terminates the running container.
func (c *Container) Stop() error {
	return c.client.cli.ContainerStop(context.Background(), c.id, c.config.shutdownTimeout)
}

func (c *Container) Cleanup() error {
	if err := c.Stop(); err != nil {
		return err
	}
	return c.client.cli.ContainerRemove(context.Background(), c.id, types.ContainerRemoveOptions{})
}

func (c *Container) GetAddress() (string, error) {
	containers, err := c.client.listContainers()
	if err != nil {
		return "", err
	}

	var existingCont *types.Container
	for _, cont := range containers {
		if cont.ID == c.id {
			existingCont = &cont
			break
		}
	}

	if existingCont == nil {
		return "", errors.New(fmt.Sprintf("container %s does not run", c.id))
	}

	var ip string
	for _, v := range existingCont.NetworkSettings.Networks {
		ip = v.IPAddress
		break // we expect only one IP address is assigned.
	}

	return ip, nil
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

	file, err := os.Create(fmt.Sprintf("%s/%s_%s.log", directory, c.config.imageName, c.id))
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

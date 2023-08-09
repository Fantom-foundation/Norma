package docker

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// projectLabel is the label used to identify objects created by norma.
const objectsLabel = "norma"

// Signal represents a signal that can be sent to a Docker container.
type Signal string

// SigHup is the SIGHUP signal.
var SigHup Signal = "SIGHUP"
var SigKill Signal = "SIGKILL"

// Client provides means to spawn Docker containers capable of hosting
// services like the go-opera client.
type Client struct {
	cli *client.Client
}

// Network represents a Docker network. It is used to connect Containers
// to each other.
type Network struct {
	id      string
	name    string
	client  *Client
	cleaned bool
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
	Entrypoint      []string // Entrypoint to run when starting the container. Optional.
	Network         *Network // Docker network to join, nil to join bridge network
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

// Purge removes all Docker objects created by norma.
func Purge() error {
	cli, err := NewClient()
	if err != nil {
		return err
	}

	// get all containers created by norma
	containers, err := cli.listContainers()
	if err != nil {
		return err
	}

	// remove all containers
	for _, c := range containers {
		// remove the container
		err = cli.cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}

	// get all networks created by norma
	networks, err := cli.listNetworks()
	if err != nil {
		return err
	}

	// remove all networks
	for _, n := range networks {
		err = cli.cli.NetworkRemove(context.Background(), n.ID)
		if err != nil {
			return err
		}
	}

	return nil
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
		Image:      config.ImageName,
		Tty:        false,
		Env:        envVars,
		Entrypoint: config.Entrypoint,
		Labels: map[string]string{
			objectsLabel: "true",
		},
	}, &container.HostConfig{
		PortBindings: portMapping,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}

	// connect to custom network if specified
	// this way the container will be connected to bridge network and
	// custom network at the same time (otherwise on network cleanup the
	// forwarded ports would be lost)
	if config.Network != nil {
		err = c.cli.NetworkConnect(context.Background(), config.Network.id, resp.ID, nil)
		if err != nil {
			return nil, err
		}
	}

	if err := network.Retry(network.DefaultRetryAttempts, 1*time.Second, func() error {
		return c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	}); err != nil {
		return nil, err
	}

	return &Container{resp.ID, c, config, false, false}, nil
}

// CreateBridgeNetwork creates a new Docker bridge network.
func (c *Client) CreateBridgeNetwork() (*Network, error) {
	// generate random name for network
	name := fmt.Sprintf("norma_network_%d", rand.Int())

	// create new network
	resp, err := c.cli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Labels: map[string]string{
			objectsLabel: "true",
		},
	})
	if err != nil {
		return nil, err
	}

	return &Network{
		id:     resp.ID,
		name:   name,
		client: c,
	}, nil
}

// Hostname returns the hostname of the Container. In this case it is the ID of the
// Docker Container.
func (c *Container) Hostname() string {
	// return the truncated container ID
	return c.id[:12]
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
	timeout := int(c.config.ShutdownTimeout.Seconds())
	return c.client.cli.ContainerStop(context.Background(), c.id, container.StopOptions{Timeout: &timeout})
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
func (c *Container) GetAddressForService(service *network.ServiceDescription) *network.AddressPort {
	// All services inside the container are reached through port-forwarding
	// on the localhost. Non-forwarded services are not supported.
	port, ok := c.config.PortForwarding[service.Port]
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

func (c *Container) StreamLog() (io.ReadCloser, error) {
	opt := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}

	reader, err := c.client.cli.ContainerLogs(context.Background(), c.id, opt)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// SendSignal sends a signal to the container.
func (c *Container) SendSignal(signal Signal) error {
	return c.client.cli.ContainerKill(context.Background(), c.id, string(signal))
}

// Exec executes a command in the container.
// This method is blocking until the command has finished.
// The output of the command is returned as a string (stdout + stderr).
// The command is required to be tokenized and interpreted in shell's exec form.
func (c *Container) Exec(cmd []string) (string, error) {
	// Create a container exec instance
	execConfig := types.ExecConfig{
		Tty:          true,
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := c.client.cli.ContainerExecCreate(context.Background(), c.id, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %s", err)
	}

	// Check if any error occurred during exec creation
	if execResp.ID == "" {
		return "", fmt.Errorf("failed to create exec instance: Empty exec ID")
	}

	// Attach to the exec instance
	resp, err := c.client.cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %s", err)
	}
	defer resp.Close()

	// Capture the output and errors
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %s", err)
	}

	// Wait for the exec command to finish
	execInspect, err := c.client.cli.ContainerExecInspect(context.Background(), execResp.ID)
	if err != nil {
		return (string)(output), fmt.Errorf("failed to inspect exec instance: %s", err)
	}

	// Check the exit code of the executed command
	if execInspect.ExitCode != 0 {
		return (string)(output), fmt.Errorf(
			"command '%s' execution failed with exit code %d", strings.Join(cmd, " "), execInspect.ExitCode)
	}

	return (string)(output), nil
}

// Cleanup removes the network from the Docker host.
func (n *Network) Cleanup() error {
	if n.cleaned {
		return nil
	}
	// remove all containers from the network, so we can remove the network
	containers, err := n.client.listContainers()
	if err != nil {
		return err
	}
	for _, c := range containers {
		for _, cn := range c.NetworkSettings.Networks {
			if cn.NetworkID == n.id {
				if err := n.client.cli.NetworkDisconnect(context.Background(), n.id, c.ID, true); err != nil {
					return err
				}
			}
		}
	}
	n.cleaned = true
	// remove the network
	return n.client.cli.NetworkRemove(context.Background(), n.id)
}

// listNetworks returns a list of all networks on the Docker host filtered by label.
func (c *Client) listNetworks() ([]types.NetworkResource, error) {
	return c.cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(getObjectsLabelFilter()),
	})
}

// listContainers returns a list of all containers on the Docker host filtered by label.
func (c *Client) listContainers() ([]types.Container, error) {
	return c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

// getObjectsLabelFilter returns a filter for the objects label.
func getObjectsLabelFilter() filters.KeyValuePair {
	return filters.Arg("label", fmt.Sprintf("%s=true", objectsLabel))
}

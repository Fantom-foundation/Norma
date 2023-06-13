package docker

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver/network"
)

func TestImplements(t *testing.T) {
	var inst Container
	var _ network.Host = &inst

}

func TestRunAndStopOneContainer(t *testing.T) {
	cli, cont := startContainer(t)

	if !containerExists(t, cli, cont.id) {
		t.Errorf("container %s is not running", cont.id)
	}

	if err := cont.Stop(); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestContainer_StopCanBeCalledMoreThanOnce(t *testing.T) {
	_, cont := startContainer(t)

	if !cont.IsRunning() {
		t.Errorf("started container is not running")
	}

	if err := cont.Stop(); err != nil {
		t.Fatalf("error stopping container: %v", err)
	}

	if cont.IsRunning() {
		t.Errorf("stopped container is still running")
	}

	if err := cont.Stop(); err != nil {
		t.Fatalf("error calling Stop() on stopped container: %v", err)
	}
}

func TestContainer_Cleanup(t *testing.T) {
	cli, cont := startContainer(t)

	if err := cont.Cleanup(); err != nil {
		t.Fatalf("error: %v", err)
	}

	if containerExists(t, cli, cont.id) {
		t.Errorf("container should not exist: %s", cont.id)
	}
}

func TestContainer_CleanupCanBeCalledMoreThanOnce(t *testing.T) {
	cli, cont := startContainer(t)

	if !cont.IsRunning() {
		t.Errorf("started container is not running")
	}

	if err := cont.Cleanup(); err != nil {
		t.Fatalf("error cleaning up container: %v", err)
	}

	if containerExists(t, cli, cont.id) {
		t.Errorf("container should no longer exist: %s", cont.id)
	}

	if cont.IsRunning() {
		t.Errorf("cleaned up container is still running")
	}

	if err := cont.Cleanup(); err != nil {
		t.Fatalf("error calling Cleanup() on cleared container: %v", err)
	}
}

func TestContainer_SaveLogTo(t *testing.T) {
	_, cont := startContainer(t)

	tmp := t.TempDir()
	if err := cont.SaveLogTo(tmp); err != nil {
		t.Fatalf("cannot save logs: %e", err)
	}

	files, err := os.ReadDir(tmp)
	if err != nil {
		log.Fatal(err)
	}

	var numOfFiles int
	for _, file := range files {
		if !file.IsDir() {
			numOfFiles++
		}
	}

	if numOfFiles == 0 {
		t.Errorf("no log files were obtained")
	}
}

func TestContainer_StreamLog(t *testing.T) {
	_, cont := startContainer(t)

	reader, err := cont.StreamLog()
	if err != nil {
		t.Fatalf("cannot read logs: %e", err)
	}

	t.Cleanup(func() {
		_ = reader.Close()
	})

	done := make(chan bool)

	go func() {
		defer close(done)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Hello from Docker!") {
				done <- true
			}
		}
	}()

	var containerStarted bool
	select {
	case containerStarted = <-done:
	case <-time.After(10 * time.Second):
	}

	if !containerStarted {
		t.Errorf("expected log not found")
	}
}

func TestContainer_StreamManyReaders(t *testing.T) {
	_, cont := startContainer(t)

	reader1, err := cont.StreamLog()
	if err != nil {
		t.Fatalf("cannot read logs: %e", err)
	}
	reader2, err := cont.StreamLog()
	if err != nil {
		t.Fatalf("cannot read logs: %e", err)
	}
	reader3, err := cont.StreamLog()
	if err != nil {
		t.Fatalf("cannot read logs: %e", err)
	}

	t.Cleanup(func() {
		_ = reader1.Close()
		_ = reader2.Close()
		_ = reader3.Close()
	})

	done := make(chan int)

	go func() {
		defer close(done)
		var count int
		for _, reader := range []io.Reader{reader1, reader2, reader3} {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "Hello from Docker!") {
					count++
				}
			}
		}
		done <- count
	}()

	var count int
	select {
	case count = <-done:
	case <-time.After(10 * time.Second):
	}

	if count != 3 {
		t.Errorf("not all readers got data: %d", count)
	}
}

func TestContainer_Exec(t *testing.T) {
	_, cont := startRunningContainer(t, nil)
	testString := "Hello world!"
	out, err := cont.Exec([]string{"sh", "-c", "echo " + testString})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(out, testString) {
		t.Errorf("expected %s, got %s", testString, out)
	}
}

func TestContainer_SendSignal(t *testing.T) {
	cli, cont := startRunningContainer(t, nil)
	if err := cont.SendSignal(SigKill); err != nil {
		t.Fatalf("error: %v", err)
	}
	// check the container is stopped
	info, err := cli.cli.ContainerInspect(context.Background(), cont.id)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if info.State.Running {
		t.Errorf("expected container to be Killed")
	}
}

func TestNetwork_Cleanup(t *testing.T) {
	cli, net := createNetwork(t)

	if err := net.Cleanup(); err != nil {
		t.Fatalf("error: %v", err)
	}

	if networkExists(t, cli, net.id) {
		t.Errorf("network should not exist: %s", net.id)
	}
}

func TestNetwork_CleanupCanBeCalledMoreThanOnce(t *testing.T) {
	cli, net := createNetwork(t)

	if err := net.Cleanup(); err != nil {
		t.Fatalf("error cleaning up network: %v", err)
	}

	if networkExists(t, cli, net.id) {
		t.Fatalf("network should no longer exist: %s", net.id)
	}

	if err := net.Cleanup(); err != nil {
		t.Fatalf("error calling Cleanup() on cleared network: %v", err)
	}
}

func TestContainerCanJoinNetwork(t *testing.T) {
	cli, net := createNetwork(t)
	_, cont := startRunningContainer(t, net)

	containers, err := cli.listContainers()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for _, c := range containers {
		if c.ID != cont.id {
			continue
		}
		for _, cn := range c.NetworkSettings.Networks {
			if cn.NetworkID == net.id {
				return
			}
		}
	}
	t.Fatalf("container is not connected to network: %s", net.id)
}

func containerExists(t *testing.T, cli *Client, id string) bool {
	// test the container exists
	var exists bool
	containers, err := cli.listContainers()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for _, c := range containers {
		if c.ID == id {
			exists = true
			break
		}
	}

	return exists
}

func networkExists(t *testing.T, cli *Client, id string) bool {
	// test the network exists
	var exists bool
	networks, err := cli.listNetworks()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for _, n := range networks {
		if n.ID == id {
			exists = true
			break
		}
	}

	return exists
}

func startContainer(t *testing.T) (*Client, *Container) {
	cli, err := NewClient()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	timeout := time.Second
	cont, err := cli.Start(&ContainerConfig{
		ImageName:       "hello-world",
		ShutdownTimeout: &timeout,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Cleanup(func() {
		_ = cont.Cleanup()
		_ = cli.Close()
	})

	return cli, cont
}

func startRunningContainer(t *testing.T, dn *Network) (*Client, *Container) {
	cli, err := NewClient()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	timeout := time.Second
	cont, err := cli.Start(&ContainerConfig{
		ImageName:       "alpine",                            // use minimal linux image
		Entrypoint:      []string{"tail", "-f", "/dev/null"}, // keep container running
		ShutdownTimeout: &timeout,
		Network:         dn,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Cleanup(func() {
		_ = cont.Cleanup()
		_ = cli.Close()
	})

	return cli, cont
}

func createNetwork(t *testing.T) (*Client, *Network) {
	cli, err := NewClient()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	net, err := cli.CreateBridgeNetwork()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Cleanup(func() {
		_ = net.Cleanup()
		_ = cli.Close()
	})

	return cli, net
}

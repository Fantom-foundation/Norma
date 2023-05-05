package docker

import (
	"log"
	"os"
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

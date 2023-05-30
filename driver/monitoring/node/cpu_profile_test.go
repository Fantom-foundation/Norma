package nodemon

import (
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
	opera "github.com/Fantom-foundation/Norma/driver/node"
)

func TestCanCollectCpuProfileDateFromOperaNode(t *testing.T) {
	docker, err := docker.NewClient()
	if err != nil {
		t.Fatalf("failed to create a docker client: %v", err)
	}
	defer docker.Close()
	node, err := opera.StartOperaDockerNode(docker, &opera.OperaNodeConfig{
		Label:         "test",
		NetworkConfig: &driver.NetworkConfig{NumberOfValidators: 1},
	})
	if err != nil {
		t.Fatalf("failed to create an Opera node on Docker: %v", err)
	}
	defer node.Cleanup()

	data, err := GetPprofData(node, time.Second)
	if err != nil {
		t.Errorf("failed to collect pprof data from node: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("fetched empty CPU profile")
	}
}

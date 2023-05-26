package netmon

import (
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
)

func TestNumNodeRetrievesNodeCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	var numNodes int
	net.EXPECT().GetActiveNodes().AnyTimes().DoAndReturn(func() []driver.Node {
		numNodes++
		nodes := make([]driver.Node, 0, numNodes)
		for i := 0; i < numNodes; i++ {
			node := driver.NewMockNode(ctrl)
			node.EXPECT().GetNodeID().Return(driver.NodeID(fmt.Sprintf("%d", i)), nil).AnyTimes()
			node.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node1TestLog)), nil)
			nodes = append(nodes, node)
		}
		return nodes
	})
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()

	writer := monitoring.NewMockWriterChain(ctrl)
	writer.EXPECT().Add(gomock.Any()).AnyTimes()

	monitor := monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer)
	numNodes = 0
	source := newNumNodesSource(monitor, 50*time.Millisecond)

	time.Sleep(200 * time.Millisecond)
	source.Shutdown()

	series, exists := source.GetData(monitoring.Network{})
	if series == nil || !exists {
		t.Fatalf("failed to obtain data from source")
	}

	data := series.GetRange(monitoring.Time(0), monitoring.Time(math.MaxInt64))
	if len(data) == 0 {
		t.Errorf("no data collected")
	}

	for i, point := range data {
		if got, want := point.Value, i+1; got != want {
			t.Errorf("invalid value recorded, got %d, wanted %d", got, want)
		}
	}
}

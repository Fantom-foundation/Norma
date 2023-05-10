package netmon

import (
	"math"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
)

func TestNumNodeRetrievesNodeCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	numNodes := 0
	net.EXPECT().GetActiveNodes().AnyTimes().DoAndReturn(func() []driver.Node {
		numNodes++
		return make([]driver.Node, numNodes)
	})

	source := NewNumNodesSource(net)

	time.Sleep(3 * time.Second)
	source.Shutdown()

	series := source.GetData(monitoring.Network{})
	if series == nil {
		t.Fatalf("failed to obtain data from source!")
	}

	data := (*series).GetRange(monitoring.Time(0), monitoring.Time(math.MaxInt64))
	if len(data) == 0 {
		t.Errorf("no data collected")
	}

	for i, point := range data {
		if got, want := point.Value, i+1; got != want {
			t.Errorf("invalid value recorded, got %d, wanted %d", got, want)
		}
	}
}

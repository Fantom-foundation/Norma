// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

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
	"github.com/uber-go/mock/gomock"
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
			node.EXPECT().GetLabel().Return(fmt.Sprintf("%d", i)).AnyTimes()
			node.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader(monitoring.Node1TestLog)), nil)
			url1 := driver.URL("node")
			node.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&url1)
			nodes = append(nodes, node)
		}
		return nodes
	})
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()

	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}
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

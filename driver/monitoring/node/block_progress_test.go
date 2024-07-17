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

package nodemon

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slices"
)

func TestNodeBlockHeightSourceRetrievesBlockHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	rpcServer, err := StartTestRpcServer()
	if err != nil {
		t.Fatalf("failed to start test RPC server: %v", rpcServer)
	}
	t.Cleanup(func() { rpcServer.Shutdown() })
	url := driver.URL(rpcServer.GetUrl())

	node1Id := driver.NodeID("A")
	node2Id := driver.NodeID("B")
	node3Id := driver.NodeID("C")

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetLabel().AnyTimes().Return(string(node1Id))
	node2.EXPECT().GetLabel().AnyTimes().Return(string(node2Id))
	node3.EXPECT().GetLabel().AnyTimes().Return(string(node3Id))

	node1.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&url)
	node2.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&url)
	node3.EXPECT().GetServiceUrl(gomock.Any()).AnyTimes().Return(&url)

	node1.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader("")), nil)
	node2.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader("")), nil)
	node3.EXPECT().StreamLog().AnyTimes().Return(io.NopCloser(strings.NewReader("")), nil)

	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().UnregisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().Return([]driver.Node{node1, node2}).AnyTimes()

	monitor, err := mon.NewMonitor(net, mon.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to initiate monitor: %v", err)
	}
	source := newNodeBlockHeightSource(monitor, 50*time.Millisecond)

	// Check that existing nodes are tracked.
	subjects := source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	want := []mon.Node{mon.Node(node1Id), mon.Node(node2Id)}
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	// Simulate the creation of a node after source initialization.
	source.(driver.NetworkListener).AfterNodeCreation(node3)

	// Check that subject list has updated.
	subjects = source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	want = append(want, mon.Node(node3Id))
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	time.Sleep(200 * time.Millisecond)
	if err := source.Shutdown(); err != nil {
		t.Errorf("erros encountered during shutdown: %v", err)
	}

	// Check that subject are still all there.
	subjects = source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	for _, subject := range subjects {
		data, exists := source.GetData(subject)
		if data == nil || !exists {
			t.Errorf("no data found for node %s", subject)
			continue
		}
		subrange := data.GetRange(mon.Time(0), mon.Time(math.MaxInt64))
		if len(subrange) == 0 {
			t.Errorf("no data collected for node %s", subject)
		}
		for _, point := range subrange {
			if got, want := point.Value, 0x12; got != want {
				t.Errorf("unexpected value collected for subject %s: wanted %d, got %d", subject, want, got)
			}
		}
	}
}

func TestLocalRpcServer_CanHandleRequests(t *testing.T) {
	server, err := StartTestRpcServer()
	if err != nil {
		t.Fatalf("failed to start the fake server: %v", err)
	}

	rpcClient, err := rpc.DialContext(context.Background(), server.GetUrl())
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}

	var result string
	err = rpcClient.Call(&result, "eth_blockNumber")
	if err != nil {
		t.Fatalf("failed to call service: %v", err)
	}

	if result != "0x12" {
		t.Errorf("invalid response: %v", result)
	}
	server.Shutdown()
}

type TestRpcServer struct {
	server *http.Server
	done   <-chan bool
}

func StartTestRpcServer() (*TestRpcServer, error) {
	port, err := network.GetFreePort()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", getBlockHeight)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	ready := make(chan error)
	done := make(chan bool)
	go func() {
		defer close(done)
		ln, err := net.Listen("tcp", server.Addr)
		if err != nil {
			ready <- err
			return
		}
		ready <- nil
		if err := server.Serve(ln); err != http.ErrServerClosed {
			fmt.Printf("server failed: %v\n", err)
		}
	}()
	if err := <-ready; err != nil {
		return nil, err
	}
	return &TestRpcServer{server, done}, nil
}

func (s *TestRpcServer) Shutdown() {
	s.server.Shutdown(context.Background())
	<-s.done
}

func (s *TestRpcServer) GetUrl() string {
	return "http://localhost" + s.server.Addr
}

func getBlockHeight(w http.ResponseWriter, r *http.Request) {
	// always returning the same
	response := `{
		"jsonrpc": "2.0",
		"id": "1234",
		"result": "0x12"
	}`
	io.WriteString(w, response)
}

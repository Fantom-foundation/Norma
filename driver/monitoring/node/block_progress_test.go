package nodemon

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang/mock/gomock"
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

	node1.EXPECT().GetNodeID().AnyTimes().Return(node1Id, nil)
	node2.EXPECT().GetNodeID().AnyTimes().Return(node2Id, nil)
	node3.EXPECT().GetNodeID().AnyTimes().Return(node3Id, nil)

	node1.EXPECT().GetHttpServiceUrl(gomock.Any()).AnyTimes().Return(&url)
	node2.EXPECT().GetHttpServiceUrl(gomock.Any()).AnyTimes().Return(&url)
	node3.EXPECT().GetHttpServiceUrl(gomock.Any()).AnyTimes().Return(&url)

	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().UnregisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().Return([]driver.Node{node1, node2})

	source := newNodeBlockHeightSource(net, 50*time.Millisecond).(*nodeBlockHeightSource)

	// Check that existing nodes are tracked.
	subjects := source.GetSubjects()
	sort.Slice(subjects, func(i, j int) bool { return subjects[i] < subjects[j] })
	want := []mon.Node{mon.Node(node1Id), mon.Node(node2Id)}
	if !slices.Equal(subjects, want) {
		t.Errorf("invalid list of subjects, wanted %v, got %v", want, subjects)
	}

	// Simulate the creation of a node after source initialization.
	source.AfterNodeCreation(node3)

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
		data := source.GetData(subject)
		if data == nil {
			t.Errorf("no data found for node %s", subject)
			continue
		}
		subrange := (*data).GetRange(mon.Time(0), mon.Time(math.MaxInt64))
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
	done := make(chan bool)
	go func() {
		server.ListenAndServe()
		close(done)
	}()
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

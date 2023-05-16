package monitoring

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestRegisterLogParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1Id := driver.NodeID("A")
	node2Id := driver.NodeID("B")
	node3Id := driver.NodeID("C")

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetNodeID().AnyTimes().Return(node1Id, nil)
	node2.EXPECT().GetNodeID().AnyTimes().Return(node2Id, nil)
	node3.EXPECT().GetNodeID().AnyTimes().Return(node3Id, nil)

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{node1, node2})

	reg := NewNodeLogDispatcher(net)
	ch := make(chan Node, 10)
	listener := &testBlockNodeListener{data: map[Node]int{}, ch: ch}
	reg.RegisterLogListener(listener)

	// simulate added node
	reg.AfterNodeCreation(node3)

	// drain 3 nodes from the channel
	for _, node := range []Node{<-ch, <-ch, <-ch} {
		blocks := listener.getBlocks(node)
		// TODO verify blocks content
		if blocks == 0 {
			t.Errorf("wrong number of collected blocks: %d for node: %v", blocks, node)
		}
	}

	if reg.getNumNodes() != 3 {
		t.Errorf("wrong size")
	}

	if len(reg.getNodes()) != 3 {
		t.Errorf("wrong number of iterations")
	}
}

type testBlockNodeListener struct {
	data     map[Node]int
	dataLock sync.Mutex
	ch       chan Node
}

func (l *testBlockNodeListener) OnBlock(node Node, b Block) {
	l.dataLock.Lock()
	defer l.dataLock.Unlock()

	// send uniq nodes
	if _, exists := l.data[node]; !exists {
		l.ch <- node
	}

	// count in only non-empty blocks
	if b.Height > 0 {
		l.data[node]++
	}
}

func (l *testBlockNodeListener) getBlocks(node Node) int {
	l.dataLock.Lock()
	defer l.dataLock.Unlock()

	return l.data[node]
}

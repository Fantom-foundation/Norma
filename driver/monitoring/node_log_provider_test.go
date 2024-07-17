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

package monitoring

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/uber-go/mock/gomock"
)

func TestLogsParsersImplements(t *testing.T) {
	var inst NodeLogDispatcher
	var _ NodeLogProvider = &inst
}

func TestRegisterLogParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)

	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)

	node1.EXPECT().GetLabel().AnyTimes().Return(string(Node1TestId))
	node2.EXPECT().GetLabel().AnyTimes().Return(string(Node2TestId))
	node3.EXPECT().GetLabel().AnyTimes().Return(string(Node3TestId))

	node1.EXPECT().StreamLog().AnyTimes().DoAndReturn(func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(Node1TestLog)), nil })
	node2.EXPECT().StreamLog().AnyTimes().DoAndReturn(func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(Node2TestLog)), nil })
	node3.EXPECT().StreamLog().AnyTimes().DoAndReturn(func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(Node3TestLog)), nil })

	// simulate existing nodes
	net.EXPECT().RegisterListener(gomock.Any())
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	dir := t.TempDir()
	reg, err := NewNodeLogDispatcher(net, dir)
	if err != nil {
		t.Fatalf("failed to create log dispatcher: %v", err)
	}
	ch := make(chan Node, 10)
	listener := &testBlockNodeListener{data: map[Node][]Block{}, ch: ch}
	reg.RegisterLogListener(listener)

	// simulate added node
	reg.AfterNodeCreation(node1)
	reg.AfterNodeCreation(node2)
	reg.AfterNodeCreation(node3)

	reg.WaitForLogsToBeConsumed()

	// drain 3 nodes from the channel
	for _, node := range []Node{<-ch, <-ch, <-ch} {
		got := listener.getBlocks(node)
		want := NodeBlockTestData[node]
		blockEqual(t, node, got, want)
	}

	if reg.getNumNodes() != 3 {
		t.Errorf("wrong size")
	}

	if len(reg.getNodes()) != 3 {
		t.Errorf("wrong number of iterations")
	}

	// Check that log got copied to output files.
	logs := []struct {
		path, content string
	}{
		{dir + "/node_logs/A.log", Node1TestLog},
		{dir + "/node_logs/B.log", Node2TestLog},
		{dir + "/node_logs/C.log", Node3TestLog},
	}
	for _, log := range logs {
		content, err := os.ReadFile(log.path)
		if err != nil {
			t.Errorf("failed to read log file: %v", err)
			continue
		}
		if got, want := log.content, string(content); got != want {
			t.Errorf("invalid log, wanted:\n%s\ngot:\n%s", want, got)
		}
	}
}

type testBlockNodeListener struct {
	data     map[Node][]Block
	dataLock sync.Mutex
	ch       chan Node
}

func blockEqual(t *testing.T, node Node, got, want []Block) {
	if len(got) != len(want) {
		t.Errorf("wrong blocks collected for Node %v: %v != %v", node, got, want)
	}

	for i, b := range got {
		if want[i].Height != b.Height || want[i].Txs != b.Txs || want[i].GasUsed != b.GasUsed {
			t.Errorf("wrong blocks collected for Node %v: %v != %v", node, want[i], b)
		}
	}
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
		l.data[node] = append(l.data[node], b)
	}
}

func (l *testBlockNodeListener) getBlocks(node Node) []Block {
	l.dataLock.Lock()
	defer l.dataLock.Unlock()

	return l.data[node]
}

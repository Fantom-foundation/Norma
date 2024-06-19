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

package checking

import (
	"strings"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/golang/mock/gomock"
)

func TestBlockHeightCheckerValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	rpc := rpc.NewMockRpcClient(ctrl)
	net.EXPECT().GetActiveNodes().MinTimes(1).Return([]driver.Node{node1, node2})
	node1.EXPECT().DialRpc().MinTimes(1).Return(rpc, nil)
	node2.EXPECT().DialRpc().MinTimes(1).Return(rpc, nil)

	blockHeight := "0x1234"
	rpc.EXPECT().Call(gomock.Any(), "eth_blockNumber").Times(2).SetArg(0, blockHeight)
	rpc.EXPECT().Close().Times(2)

	err := new(BlockHeightChecker).Check(net)
	if err != nil {
		t.Errorf("unexpected error from BlockHeightChecker: %v", err)
	}
}

func TestBlockHeightCheckerInvalid(t *testing.T) {
	tests := []struct {
		name         string
		blockHeight1 string
		blockHeight2 string
	}{
		{name: "ascending", blockHeight1: "0x42", blockHeight2: "0x1234"},
		{name: "descending", blockHeight1: "0x1234", blockHeight2: "0x42"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			net := driver.NewMockNetwork(ctrl)
			node1 := driver.NewMockNode(ctrl)
			node2 := driver.NewMockNode(ctrl)
			rpc1 := rpc.NewMockRpcClient(ctrl)
			rpc2 := rpc.NewMockRpcClient(ctrl)
			net.EXPECT().GetActiveNodes().MinTimes(1).Return([]driver.Node{node1, node2})
			node1.EXPECT().DialRpc().MinTimes(1).Return(rpc1, nil)
			node2.EXPECT().DialRpc().MinTimes(1).Return(rpc2, nil)
			node1.EXPECT().GetLabel().AnyTimes().Return("node1")
			node2.EXPECT().GetLabel().AnyTimes().Return("node2")

			rpc1.EXPECT().Call(gomock.Any(), "eth_blockNumber").SetArg(0, test.blockHeight1)
			rpc2.EXPECT().Call(gomock.Any(), "eth_blockNumber").SetArg(0, test.blockHeight2)
			rpc1.EXPECT().Close()
			rpc2.EXPECT().Close()

			err := new(BlockHeightChecker).Check(net)
			if err == nil || !strings.Contains(err.Error(), "reports too old block") {
				t.Errorf("unexpected error from BlockHeightChecker: %v", err)
			}
		})
	}
}

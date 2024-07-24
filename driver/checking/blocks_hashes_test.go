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
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"
)

func TestBlockHashesCheckerValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	rpc := rpc.NewMockRpcClient(ctrl)
	net.EXPECT().GetActiveNodes().MinTimes(1).Return([]driver.Node{node1, node2})
	node1.EXPECT().DialRpc().MinTimes(1).Return(rpc, nil)
	node2.EXPECT().DialRpc().MinTimes(1).Return(rpc, nil)
	result := blockHashes{
		Hash:         common.Hash{0x11},
		StateRoot:    common.Hash{0x22},
		ReceiptsRoot: common.Hash{0x33},
	}
	gomock.InOrder(
		rpc.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).Times(6).SetArg(0, &result),
		rpc.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).AnyTimes(),
		rpc.EXPECT().Close().Times(2),
	)
	err := new(BlocksHashesChecker).Check(net)
	if err != nil {
		t.Errorf("unexpected error from BlocksHashesChecker: %v", err)
	}
}

func TestBlockHashesCheckerInvalidStateRoot(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	rpc1 := rpc.NewMockRpcClient(ctrl)
	rpc2 := rpc.NewMockRpcClient(ctrl)
	net.EXPECT().GetActiveNodes().MinTimes(1).Return([]driver.Node{node1, node2})
	node1.EXPECT().DialRpc().MinTimes(1).Return(rpc1, nil)
	node2.EXPECT().DialRpc().MinTimes(1).Return(rpc2, nil)
	result1 := blockHashes{
		Hash:         common.Hash{0x11},
		StateRoot:    common.Hash{0x22},
		ReceiptsRoot: common.Hash{0x33},
	}
	result2 := blockHashes{
		Hash:         common.Hash{0x11},
		StateRoot:    common.Hash{0xFF}, // different
		ReceiptsRoot: common.Hash{0x33},
	}

	rpc1.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).AnyTimes().SetArg(0, &result1)
	rpc2.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).Times(3).SetArg(0, &result1)
	rpc2.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).SetArg(0, &result2)
	rpc1.EXPECT().Close()
	rpc2.EXPECT().Close()

	err := new(BlocksHashesChecker).Check(net)
	if err.Error() != "stateRoot of the block 3 does not match" {
		t.Errorf("unexpected error from BlocksHashesChecker: %v", err)
	}
}

func TestBlockHashesCheckerInvalidLastBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	net := driver.NewMockNetwork(ctrl)
	node1 := driver.NewMockNode(ctrl)
	node2 := driver.NewMockNode(ctrl)
	node3 := driver.NewMockNode(ctrl)
	rpc1 := rpc.NewMockRpcClient(ctrl)
	rpc2 := rpc.NewMockRpcClient(ctrl)
	rpc3 := rpc.NewMockRpcClient(ctrl)
	net.EXPECT().GetActiveNodes().MinTimes(1).Return([]driver.Node{node1, node2, node3})
	node1.EXPECT().DialRpc().MinTimes(1).Return(rpc1, nil)
	node2.EXPECT().DialRpc().MinTimes(1).Return(rpc2, nil)
	node3.EXPECT().DialRpc().MinTimes(1).Return(rpc3, nil)
	result1 := blockHashes{
		Hash:         common.Hash{0x11},
		StateRoot:    common.Hash{0x22},
		ReceiptsRoot: common.Hash{0x33},
	}
	result2 := blockHashes{
		Hash:         common.Hash{0x11},
		StateRoot:    common.Hash{0x22},
		ReceiptsRoot: common.Hash{0xFF}, // different
	}

	rpc1.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).Times(4).SetArg(0, &result1)
	rpc1.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).AnyTimes()

	rpc2.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).Times(3).SetArg(0, &result1)
	rpc2.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).AnyTimes() // does not have block 3 (should be ignored)

	rpc3.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).Times(3).SetArg(0, &result1)
	rpc3.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", "0x3", false).SetArg(0, &result2) // different block 3
	rpc3.EXPECT().Call(gomock.Any(), "eth_getBlockByNumber", gomock.Any(), false).AnyTimes()

	rpc1.EXPECT().Close()
	rpc2.EXPECT().Close()
	rpc3.EXPECT().Close()

	err := new(BlocksHashesChecker).Check(net)
	if err.Error() != "receiptsRoot of the block 3 does not match" {
		t.Errorf("unexpected error from BlocksHashesChecker: %v", err)
	}
}

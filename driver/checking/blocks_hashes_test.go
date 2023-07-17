package checking

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"testing"
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

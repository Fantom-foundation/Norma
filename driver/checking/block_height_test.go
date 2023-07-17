package checking

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
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
		t.Errorf("unexpected error from BlocksHashesChecker: %v", err)
	}
}

func TestBlockHeightCheckerInvalid(t *testing.T) {
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

	blockHeight1 := "0x1234"
	blockHeight2 := "0x42"
	rpc1.EXPECT().Call(gomock.Any(), "eth_blockNumber").SetArg(0, blockHeight1)
	rpc2.EXPECT().Call(gomock.Any(), "eth_blockNumber").SetArg(0, blockHeight2)
	rpc1.EXPECT().Close()
	rpc2.EXPECT().Close()

	err := new(BlockHeightChecker).Check(net)
	if err == nil || !strings.Contains(err.Error(), "reports too old block") {
		t.Errorf("unexpected error from BlocksHashesChecker: %v", err)
	}
}

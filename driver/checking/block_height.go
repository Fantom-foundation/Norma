package checking

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"strconv"
	"strings"
)

// BlockHeightChecker is a Checker checking if all Opera nodes achieved the same block height.
type BlockHeightChecker struct {
}

func (*BlockHeightChecker) Check(net driver.Network) error {
	nodes := net.GetActiveNodes()
	heights := make([]int64, len(nodes))
	maxHeight := int64(0)
	for i, n := range nodes {
		height, err := getBlockHeight(n)
		if err != nil {
			return fmt.Errorf("failed to get block height of node %s; %v", n.GetLabel(), err)
		}
		if height == 1 {
			return fmt.Errorf("node %s reports it is at block 1 (only genesis is applied)", n.GetLabel())
		}
		if height < 1 {
			return fmt.Errorf("node %s reports it is at invalid block %d", n.GetLabel(), height)
		}
		if maxHeight < height {
			maxHeight = height
		}
		heights[i] = height
	}
	for i, n := range nodes {
		if heights[i] < maxHeight-1 {
			return fmt.Errorf("node %s reports too old block %d (max block is %d)", n.GetLabel(), heights[i], maxHeight)
		}
	}
	return nil
}

func getBlockHeight(n driver.Node) (int64, error) {
	rpcClient, err := n.DialRpc()
	if err != nil {
		return 0, fmt.Errorf("failed to dial node RPC; %v", err)
	}
	defer rpcClient.Close()
	var blockNumber string
	err = rpcClient.Call(&blockNumber, "eth_blockNumber")
	if err != nil {
		return 0, fmt.Errorf("failed to get block number from RPC; %v", err)
	}
	blockNumber = strings.TrimPrefix(blockNumber, "0x")
	return strconv.ParseInt(blockNumber, 16, 64)
}

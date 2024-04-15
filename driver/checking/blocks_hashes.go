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
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlocksHashesChecker is a Checker checking if all Opera nodes provides the same hashes for all blocks/stateRoots.
type BlocksHashesChecker struct {
}

func (*BlocksHashesChecker) Check(net driver.Network) (err error) {
	nodes := net.GetActiveNodes()
	rpcClients := make([]rpc.RpcClient, len(nodes))
	for i, n := range nodes {
		rpcClients[i], err = n.DialRpc()
		if err != nil {
			return fmt.Errorf("failed to dial RPC for node %s; %v", n.GetLabel(), err)
		}
	}
	defer func() {
		for _, rpcClient := range rpcClients {
			rpcClient.Close()
		}
	}()

	for blockNumber := uint64(0); ; blockNumber++ {
		var referenceHashes *blockHashes
		var nodesLackingTheBlock = 0
		for i, n := range nodes {
			block, err := getBlockHashes(rpcClients[i], blockNumber)
			if err != nil {
				return fmt.Errorf("failed to get block %d detail at node %s; %v", blockNumber, n.GetLabel(), err)
			}
			if block == nil { // block does not exist on the node
				if blockNumber <= 2 {
					return fmt.Errorf("unable to check block hashes - block %d does not exists at node %s", blockNumber, n.GetLabel())
				}
				nodesLackingTheBlock++
				continue
			}
			if referenceHashes == nil {
				referenceHashes = block
			} else {
				if referenceHashes.StateRoot != block.StateRoot {
					return fmt.Errorf("stateRoot of the block %d does not match", blockNumber)
				}
				if referenceHashes.ReceiptsRoot != block.ReceiptsRoot {
					return fmt.Errorf("receiptsRoot of the block %d does not match", blockNumber)
				}
				if referenceHashes.Hash != block.Hash {
					return fmt.Errorf("hash of the block %d does not match", blockNumber)
				}
			}
		}
		if nodesLackingTheBlock == len(nodes) { // no node has the last block
			return nil // finish successfully
		}
	}
}

type blockHashes struct {
	Hash         common.Hash
	StateRoot    common.Hash
	ReceiptsRoot common.Hash
}

func getBlockHashes(rpcClient rpc.RpcClient, blockNumber uint64) (*blockHashes, error) {
	var block *blockHashes
	err := rpcClient.Call(&block, "eth_getBlockByNumber", hexutil.EncodeUint64(blockNumber), false)
	if err != nil {
		return nil, fmt.Errorf("failed to get block state root from RPC; %v", err)
	}
	return block, nil
}

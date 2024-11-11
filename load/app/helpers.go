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

package app

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/ethereum/go-ethereum/common"
)

// WaitUntilAccountNonceIs blocks until the account nonce at the latest block on the chain is given value
func WaitUntilAccountNonceIs(account common.Address, awaitedNonce uint64, rpcClient rpc.RpcClient) error {
	var nonce uint64
	var err error
	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		nonce, err = rpcClient.NonceAt(context.Background(), account, nil) // nonce at latest block
		if err != nil {
			return fmt.Errorf("failed to check address nonce; %v", err)
		}
		if nonce == awaitedNonce {
			return nil
		}
	}
	return fmt.Errorf("nonce not achieved before timeout (awaited %d, current %d)", awaitedNonce, nonce)
}

// GetGasPrice obtains optimal gasPrice for regular transactions
func GetGasPrice(rpcClient rpc.RpcClient) (*big.Int, error) {
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price; %v", err)
	}
	var regularPrice big.Int
	regularPrice.Mul(gasPrice, big.NewInt(2)) // lower gas price for regular txs (but more than suggested by Opera)
	return &regularPrice, nil
}

func getPriorityGasPrice(regularGasPrice *big.Int) *big.Int {
	var priorityPrice big.Int
	priorityPrice.Mul(regularGasPrice, big.NewInt(2)) // greater gas price for init
	return &priorityPrice
}

func reverseAddresses(in []common.Address) []common.Address {
	out := make([]common.Address, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[len(in)-1-i]
	}
	return out
}

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
	"github.com/ethereum/go-ethereum/core/types"
)

// transferValue transfer a financial value from account identified by given privateKey, to given toAddress.
// It returns when the value is already available on the target account.
func transferValue(rpcClient rpc.RpcClient, from *Account, toAddress common.Address, value *big.Int, gasPrice *big.Int) (err error) {
	signedTx, err := createTx(from, toAddress, value, nil, gasPrice, 21000)
	if err != nil {
		return err
	}
	return rpcClient.SendTransaction(context.Background(), signedTx)
}

func createTx(from *Account, toAddress common.Address, value *big.Int, data []byte, gasPrice *big.Int, gasLimit uint64) (*types.Transaction, error) {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    from.getNextNonce(),
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})
	return types.SignTx(tx, types.NewEIP155Signer(from.chainID), from.privateKey)
}

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

// waitUntilAllSentTxsAreOnChain blocks until all txs sent from given accounts are on the chain (by account nonce)
func waitUntilAllSentTxsAreOnChain(accounts []*Account, rpcClient rpc.RpcClient) error {
	for i := 0; i < len(accounts); i++ {
		err := WaitUntilAccountNonceIs(accounts[i].address, accounts[i].getCurrentNonce(), rpcClient)
		if err != nil {
			return err
		}
	}
	return nil
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

func generateStartingAccounts(rpcClient rpc.RpcClient, primaryAccount *Account, factory *AccountFactory, numAccounts int, regularGasPrice *big.Int) ([]*Account, error) {
	var err error
	startingAccounts := make([]*Account, numAccounts/500+1)
	for i := 0; i < len(startingAccounts); i++ {
		startingAccounts[i], err = factory.CreateAccount(rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create starting account %d; %v", i, err)
		}
		err = startingAccounts[i].Fund(primaryAccount, rpcClient, regularGasPrice, 1_000_000)
		if err != nil {
			return nil, fmt.Errorf("failed to fund starting account %d; %v", i, err)
		}
	}
	return startingAccounts, nil
}

func reverseAddresses(in []common.Address) []common.Address {
	out := make([]common.Address, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[len(in)-1-i]
	}
	return out
}

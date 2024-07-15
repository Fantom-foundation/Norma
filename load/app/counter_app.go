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
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/Fantom-foundation/Norma/driver/rpc"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewCounterApplication deploys a Counter contract to the chain.
// The Counter contract is a simple contract sustaining an integer value, to be incremented by sent txs.
// It allows to easily test the tx generating, as reading the contract field provides the amount of applied contract calls.
func NewCounterApplication(rpcClient rpc.RpcClient, primaryAccount *Account, numUsers int, feederId, appId uint32) (Application, error) {
	// get price of gas from the network
	regularGasPrice, err := getGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// Deploy the Counter contract to be used by tx generators
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	contractAddress, _, _, err := contract.DeployCounter(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Counter contract; %v", err)
	}

	accountFactory, err := NewAccountFactory(primaryAccount.chainID, feederId, appId)
	if err != nil {
		return nil, err
	}

	// deploying too many generators from one account leads to excessive gasPrice growth - we
	// need to spread the initialization in between multiple startingAccounts
	startingAccounts, err := generateStartingAccounts(rpcClient, primaryAccount, accountFactory, numUsers, regularGasPrice)
	if err != nil {
		return nil, err
	}

	// parse ABI for generating txs data
	parsedAbi, err := contract.CounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	err = WaitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Counter contract is deployed; %v", err)
	}

	return &CounterApplication{
		abi:              parsedAbi,
		startingAccounts: startingAccounts,
		contractAddress:  contractAddress,
		accountFactory:   accountFactory,
	}, nil
}

// CounterApplication represents a simple on-chain Counter incremented by sent transactions.
// A factory represents one deployed Counter contract, incremented by all its generators.
// While the factory is thread-safe, each created generator should be used in a single thread only.
type CounterApplication struct {
	abi              *abi.ABI
	startingAccounts []*Account
	contractAddress  common.Address
	accountFactory   *AccountFactory
}

// CreateUser creates a new user for the app.
func (f *CounterApplication) CreateUser(rpcClient rpc.RpcClient) (User, error) {

	// get price of gas from the network
	regularGasPrice, err := getGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// Generate a new account for each worker - avoid account nonces related bottlenecks
	workerAccount, err := f.accountFactory.CreateAccount(rpcClient)
	if err != nil {
		return nil, err
	}
	startingAccount := f.startingAccounts[workerAccount.id%len(f.startingAccounts)]
	err = workerAccount.Fund(startingAccount, rpcClient, regularGasPrice, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to fund worker account %d; %v", workerAccount.id, err)
	}

	gen := &CounterUser{
		abi:      f.abi,
		sender:   workerAccount,
		gasPrice: regularGasPrice,
		contract: f.contractAddress,
	}
	return gen, nil
}

func (f *CounterApplication) WaitUntilApplicationIsDeployed(rpcClient rpc.RpcClient) error {
	return waitUntilAllSentTxsAreOnChain(f.startingAccounts, rpcClient)
}

func (f *CounterApplication) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	counterContract, err := contract.NewCounter(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get Counter contract representation; %v", err)
	}
	count, err := counterContract.GetCount(nil)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

// CounterUser represents a user sending txs to increment a trivial Counter contract value.
// A generator is supposed to be used in a single thread.
type CounterUser struct {
	abi      *abi.ABI
	sender   *Account
	gasPrice *big.Int
	contract common.Address
	sentTxs  uint64
}

func (g *CounterUser) GenerateTx() (*types.Transaction, error) {
	// prepare tx data
	data, err := g.abi.Pack("incrementCounter")
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	const gasLimit = 50000 // IncrementCounter method call takes 43426 of gas
	tx, err := createTx(g.sender, g.contract, big.NewInt(0), data, g.gasPrice, gasLimit)
	if err == nil {
		atomic.AddUint64(&g.sentTxs, 1)
	}
	return tx, err
}

func (g *CounterUser) GetSentTransactions() uint64 {
	return atomic.LoadUint64(&g.sentTxs)
}

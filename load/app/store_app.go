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

// NewStoreApplication deploys a Store contract to the chain.
// The Store contract is a simple contract managing a user-private key/value store.
// It is intended to produce state-heavy transactions.
func NewStoreApplication(rpcClient rpc.RpcClient, primaryAccount *Account, numUsers int, feederId, appId uint32) (Application, error) {
	// get price of gas from the network
	regularGasPrice, err := GetGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// Deploy the Store contract to be used by tx generators
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	contractAddress, _, _, err := contract.DeployStore(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Store contract; %v", err)
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
	parsedAbi, err := contract.StoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	err = WaitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Store contract is deployed; %v", err)
	}

	return &StoreApplication{
		abi:              parsedAbi,
		startingAccounts: startingAccounts,
		contractAddress:  contractAddress,
		accountFactory:   accountFactory,
	}, nil
}

// StoreApplication represents a simple on-chain user-private Key/Value store.
// A instance represents one deployed Store contract as well as a set of users.
type StoreApplication struct {
	abi              *abi.ABI
	startingAccounts []*Account
	contractAddress  common.Address
	accountFactory   *AccountFactory
}

// CreateUser creates a new user for the app.
func (f *StoreApplication) CreateUser(rpcClient rpc.RpcClient) (User, error) {

	// get price of gas from the network
	regularGasPrice, err := GetGasPrice(rpcClient)
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

	gen := &StoreUser{
		abi:      f.abi,
		sender:   workerAccount,
		contract: f.contractAddress,
	}
	return gen, nil
}

func (f *StoreApplication) WaitUntilApplicationIsDeployed(rpcClient rpc.RpcClient) error {
	return waitUntilAllSentTxsAreOnChain(f.startingAccounts, rpcClient)
}

func (f *StoreApplication) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	storeContract, err := contract.NewStore(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get Store contract representation; %v", err)
	}
	count, err := storeContract.GetCount(nil)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

// StoreUser represents a user sending txs to manipulate a user-private key/value store.
// Instances are not thread safe.
type StoreUser struct {
	abi      *abi.ABI
	sender   *Account
	contract common.Address
	sentTxs  atomic.Uint64
}

func (g *StoreUser) GenerateTx(currentGasPrice *big.Int) (*types.Transaction, error) {
	const updateSize = 260 // ~ 1 GB/minute new netto data at 1000 Tx/s

	// prepare tx data -- since as single put is rather cheap, we use the 'fill' operation
	// to perform a number of updates at once. Each transaction is allocating updateSize
	// extra slots, which correspond to ~(32 byte key + 32 byte value) extra storage.
	val := int64(g.sentTxs.Load()) + 1
	from := val * updateSize
	to := from + updateSize
	data, err := g.abi.Pack("fill", big.NewInt(from), big.NewInt(to), big.NewInt(val))
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	const gasLimit = 52000 + 25000*updateSize // wild guess ...
	tx, err := createTx(g.sender, g.contract, big.NewInt(0), data, currentGasPrice, gasLimit)
	if err == nil {
		g.sentTxs.Add(1)
	}
	return tx, err
}

func (g *StoreUser) GetSentTransactions() uint64 {
	return g.sentTxs.Load()
}

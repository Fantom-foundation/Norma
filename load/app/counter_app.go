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
	"sync"
	"sync/atomic"

	"github.com/Fantom-foundation/Norma/driver/rpc"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// NewCounterApplication deploys a Counter contract to the chain.
// The Counter contract is a simple contract sustaining an integer value, to be incremented by sent txs.
// It allows to easily test the tx generating, as reading the contract field provides the amount of applied contract calls.
func NewCounterApplication(ctxt AppContext, feederId, appId uint32) (Application, error) {
	client := ctxt.GetClient()
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID; %w", err)
	}

	// Deploy the Counter contract to be used by this application.
	_, receipt, err := DeployContract(ctxt, contract.DeployCounter)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Counter contract; %w", err)
	}

	accountFactory, err := NewAccountFactory(chainId, feederId, appId)
	if err != nil {
		return nil, err
	}

	return &CounterApplication{
		contractAddress: receipt.ContractAddress,
		accountFactory:  accountFactory,
	}, nil
}

// CounterApplication represents a simple on-chain Counter incremented by sent transactions.
// A factory represents one deployed Counter contract, incremented by all its generators.
// While the factory is thread-safe, each created generator should be used in a single thread only.
type CounterApplication struct {
	contractAddress common.Address
	accountFactory  *AccountFactory
}

// CreateUsers creates a list of new users for the app.
func (f *CounterApplication) CreateUsers(appContext AppContext, numUsers int) ([]User, error) {

	users := make([]User, numUsers)
	addresses := make([]common.Address, numUsers)
	for i := 0; i < numUsers; i++ {
		// Generate a new account for each worker - avoid account nonces related bottlenecks
		workerAccount, err := f.accountFactory.CreateAccount(appContext.GetClient())
		if err != nil {
			return nil, err
		}
		users[i] = &CounterUser{
			sender:   workerAccount,
			contract: f.contractAddress,
		}
		addresses[i] = workerAccount.address
	}

	fundsPerUser := big.NewInt(1_000)
	fundsPerUser = new(big.Int).Mul(fundsPerUser, big.NewInt(1_000_000_000_000_000_000)) // to wei
	err := appContext.FundAccounts(addresses, fundsPerUser)

	return users, err
}

func (f *CounterApplication) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	counterContract, err := contract.NewCounter(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get Counter contract representation; %w", err)
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
	mu       sync.Mutex
	sender   *Account
	contract common.Address
	sentTxs  atomic.Uint64
	opts     *bind.TransactOpts
}

func (g *CounterUser) SendTransaction(rpcClient rpc.RpcClient) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	contract, err := contract.NewCounter(g.contract, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get Counter contract proxy; %w", err)
	}

	if g.opts == nil {
		g.opts, err = GetTransactOptions(rpcClient, g.sender)
		if err != nil {
			return fmt.Errorf("failed to get transaction options; %w", err)
		}
		g.opts.GasLimit = 28036
	}

	_, err = contract.IncrementCounter(g.opts)
	if err != nil {
		return fmt.Errorf("failed to create transaction; %w", err)
	}
	g.opts.Nonce.Add(g.opts.Nonce, big.NewInt(1))
	g.sentTxs.Add(1)
	return nil
}

func (g *CounterUser) GetTotalNumberOfSentTransactions() uint64 {
	return g.sentTxs.Load()
}

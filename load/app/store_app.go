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
	"sync/atomic"

	"github.com/Fantom-foundation/Norma/driver/rpc"

	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewStoreApplication deploys a Store contract to the chain.
// The Store contract is a simple contract managing a user-private key/value store.
// It is intended to produce state-heavy transactions.
func NewStoreApplication(ctxt AppContext, feederId, appId uint32) (Application, error) {

	client := ctxt.GetClient()
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID; %w", err)
	}

	// Deploy the Store contract to be used by this application.
	_, receipt, err := DeployContract(ctxt, contract.DeployStore)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Store contract; %w", err)
	}

	accountFactory, err := NewAccountFactory(chainId, feederId, appId)
	if err != nil {
		return nil, err
	}

	return &StoreApplication{
		contractAddress: receipt.ContractAddress,
		accountFactory:  accountFactory,
	}, nil
}

// StoreApplication represents a simple on-chain user-private Key/Value store.
// A instance represents one deployed Store contract as well as a set of users.
type StoreApplication struct {
	contractAddress common.Address
	accountFactory  *AccountFactory
}

// CreateUsers creates a list of new users for the app.
func (f *StoreApplication) CreateUsers(appContext AppContext, numUsers int) ([]User, error) {

	users := make([]User, numUsers)
	addresses := make([]common.Address, numUsers)
	for i := 0; i < numUsers; i++ {
		// Generate a new account for each worker - avoid account nonces related bottlenecks
		workerAccount, err := f.accountFactory.CreateAccount(appContext.GetClient())
		if err != nil {
			return nil, err
		}
		users[i] = &StoreUser{
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

func (f *StoreApplication) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	storeContract, err := contract.NewStore(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get Store contract representation; %w", err)
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
	sender   *Account
	contract common.Address
	sentTxs  atomic.Uint64
}

func (g *StoreUser) SendTransaction(rpcClient rpc.RpcClient) error {
	const updateSize = 260 // ~ 1 GB/minute new net data at 1000 Tx/s

	contract, err := contract.NewStore(g.contract, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get Store contract proxy; %w", err)
	}

	// prepare tx data -- since as single put is rather cheap, we use the 'fill' operation
	// to perform a number of updates at once. Each transaction is allocating updateSize
	// extra slots, which correspond to ~(32 byte key + 32 byte value) extra storage.
	val := int64(g.sentTxs.Load()) + 1
	from := val * updateSize
	to := from + updateSize
	receipt, err := Run(rpcClient, g.sender, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Fill(opts, big.NewInt(from), big.NewInt(to), big.NewInt(val))
	})
	if err != nil {
		return fmt.Errorf("failed to execute a transaction: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction reverted")
	}
	g.sentTxs.Add(1)
	return nil
}

func (g *StoreUser) GetTotalNumberOfSentTransactions() uint64 {
	return g.sentTxs.Load()
}

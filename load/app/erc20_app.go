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
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"

	"github.com/Fantom-foundation/Norma/driver/rpc"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewERC20Application deploys a new ERC-20 dapp to the chain.
// The ERC20 contract is a contract sustaining balances of the token for individual owner addresses.
func NewERC20Application(ctxt AppContext, feederId, appId uint32) (Application, error) {
	rpcClient := ctxt.GetClient()
	primaryAccount := ctxt.GetTreasure()

	// Deploy the ERC20 contract to be used by generators created using the factory
	txOpts, err := ctxt.GetTransactOptions(primaryAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %w", err)
	}
	contractAddress, transaction, _, err := contract.DeployERC20(txOpts, rpcClient, "Testing Token", "TOK")
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ERC20 contract; %w", err)
	}
	recipients, err := generateRecipientsAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipients addresses; %w", err)
	}

	accountFactory, err := NewAccountFactory(primaryAccount.chainID, feederId, appId)
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	_, err = ctxt.GetReceipt(transaction.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the ERC20 contract is deployed; %w", err)
	}

	return &ERC20Application{
		contractAddress: contractAddress,
		recipients:      recipients,
		accountFactory:  accountFactory,
	}, nil
}

func generateRecipientsAddresses() ([]common.Address, error) {
	recipients := make([]common.Address, 100)
	for i := 0; i < 100; i++ {
		_, err := crand.Read(recipients[i][:])
		if err != nil {
			return nil, err
		}
	}
	return recipients, nil
}

// ERC20Application represents one application deployed to the network - an ERC-20 contract.
// Each created app should be used in a single thread only.
type ERC20Application struct {
	contractAddress common.Address
	recipients      []common.Address
	accountFactory  *AccountFactory
}

// CreateUsers creates a list of new users for the app.
func (f *ERC20Application) CreateUsers(appContext AppContext, numUsers int) ([]User, error) {

	// Create a list of users.
	users := make([]User, numUsers)
	addresses := make([]common.Address, numUsers)
	for i := 0; i < numUsers; i++ {
		// Generate a new account for each worker - avoid account nonces related bottlenecks
		workerAccount, err := f.accountFactory.CreateAccount(appContext.GetClient())
		if err != nil {
			return nil, err
		}
		users[i] = &ERC20User{
			sender:     workerAccount,
			contract:   f.contractAddress,
			recipients: f.recipients,
		}
		addresses[i] = workerAccount.address
	}

	// Provide native currency to each user.
	fundsPerUser := big.NewInt(1_000)
	fundsPerUser = new(big.Int).Mul(fundsPerUser, big.NewInt(1_000_000_000_000_000_000)) // to wei
	err := appContext.FundAccounts(addresses, fundsPerUser)
	if err != nil {
		return nil, fmt.Errorf("failed to fund accounts; %w", err)
	}

	// Provide ERC-20 tokens to each user.
	erc20Contract, err := contract.NewERC20(f.contractAddress, appContext.GetClient())
	if err != nil {
		return nil, fmt.Errorf("failed to get ERC20 contract representation; %w", err)
	}
	receipt, err := appContext.Run(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return erc20Contract.MintForAll(opts, addresses, big.NewInt(1_000000000000000000))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mint ERC-20 for all users; %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("failed to mint ERC-20 for all users; transaction reverted")
	}
	return users, nil
}

func (f *ERC20Application) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	ERC20Contract, err := contract.NewERC20(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get ERC20 contract representation; %w", err)
	}
	totalReceived := uint64(0)
	for _, recipient := range f.recipients {
		recipientBalance, err := ERC20Contract.BalanceOf(nil, recipient)
		if err != nil {
			return 0, err
		}
		totalReceived += recipientBalance.Uint64()
	}
	return totalReceived, nil
}

// ERC20User represents a user sending txs to transfer ERC20 tokens.
// A generator is supposed to be used in a single thread.
type ERC20User struct {
	sender     *Account
	contract   common.Address
	recipients []common.Address
	sentTxs    atomic.Uint64
}

func (g *ERC20User) SendTransaction(rpcClient rpc.RpcClient) error {
	contract, err := contract.NewERC20(g.contract, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get ERC20 contract proxy; %w", err)
	}

	recipient := g.recipients[rand.Intn(len(g.recipients))]
	receipt, err := Run(rpcClient, g.sender, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Transfer(opts, recipient, big.NewInt(1))
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

func (g *ERC20User) GetTotalNumberOfSentTransactions() uint64 {
	return g.sentTxs.Load()
}

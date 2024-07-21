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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewERC20Application deploys a new ERC-20 dapp to the chain.
// The ERC20 contract is a contract sustaining balances of the token for individual owner addresses.
func NewERC20Application(rpcClient rpc.RpcClient, primaryAccount *Account, numUsers int, feederId, appId uint32) (Application, error) {
	// get price of gas from the network
	regularGasPrice, err := GetGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// Deploy the ERC20 contract to be used by generators created using the factory
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	contractAddress, _, _, err := contract.DeployERC20(txOpts, rpcClient, "Testing Token", "TOK")
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ERC20 contract; %v", err)
	}
	recipients, err := generateRecipientsAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipients addresses; %v", err)
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
	parsedAbi, err := contract.ERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	err = WaitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the ERC20 contract is deployed; %v", err)
	}

	return &ERC20Application{
		abi:              parsedAbi,
		startingAccounts: startingAccounts,
		contractAddress:  contractAddress,
		recipients:       recipients,
		accountFactory:   accountFactory,
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
	abi              *abi.ABI
	startingAccounts []*Account
	contractAddress  common.Address
	recipients       []common.Address
	accountFactory   *AccountFactory
}

// CreateUser creates a new user for the app.
func (f *ERC20Application) CreateUser(rpcClient rpc.RpcClient) (User, error) {
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

	// mint ERC-20 tokens for the worker account - tokens to be transferred in the transactions
	erc20Contract, err := contract.NewERC20(f.contractAddress, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get ERC20 contract representation; %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(startingAccount.privateKey, startingAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(startingAccount.getNextNonce()))
	_, err = erc20Contract.Mint(txOpts, workerAccount.address, big.NewInt(1_000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("failed to mint ERC-20; %v", err)
	}

	return &ERC20User{
		abi:        f.abi,
		sender:     workerAccount,
		contract:   f.contractAddress,
		recipients: f.recipients,
	}, nil
}

func (f *ERC20Application) WaitUntilApplicationIsDeployed(rpcClient rpc.RpcClient) error {
	return waitUntilAllSentTxsAreOnChain(f.startingAccounts, rpcClient)
}

func (f *ERC20Application) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	ERC20Contract, err := contract.NewERC20(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get ERC20 contract representation; %v", err)
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
	abi        *abi.ABI
	sender     *Account
	contract   common.Address
	recipients []common.Address
	sentTxs    uint64
}

func (g *ERC20User) GenerateTx(currentGasPrice *big.Int) (*types.Transaction, error) {
	// choose random recipient
	recipient := g.recipients[rand.Intn(len(g.recipients))]

	// prepare tx data
	data, err := g.abi.Pack("transfer", recipient, big.NewInt(1))
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	const gasLimit = 52000 // Transfer method call takes 51349 of gas
	tx, err := createTx(g.sender, g.contract, big.NewInt(0), data, currentGasPrice, gasLimit)
	if err == nil {
		atomic.AddUint64(&g.sentTxs, 1)
	}
	return tx, err
}

func (g *ERC20User) GetSentTransactions() uint64 {
	return atomic.LoadUint64(&g.sentTxs)
}

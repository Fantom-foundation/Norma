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
	"bytes"
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

const TokensInChain = 4
const PairsInChain = TokensInChain - 1

var AmountSwapped = big.NewInt(100) // swapped in one tx
var WorkerInitialBalance = big.NewInt(0).Mul(big.NewInt(1_000_000_000), big.NewInt(1_000000000000000000))
var PairLiquidity = big.NewInt(0).Mul(big.NewInt(1_000_000_000_000_000), big.NewInt(1_000000000000000000))

// NewUniswapApplication deploys a new Uniswap dapp to the chain.
// Created Uniswap pairs allows to swap first ERC-20 token for second, second for third etc.
// This app swaps first token for the last one, using all intermediate tokens.
func NewUniswapApplication(rpcClient rpc.RpcClient, primaryAccount *Account, numUsers int, feederId, appId uint32) (Application, error) {
	// get price of gas from the network
	regularGasPrice, err := GetGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	tokenAddresses := make([]common.Address, TokensInChain)
	tokenContracts := make([]*contract.ERC20, TokensInChain)
	pairsAddresses := make([]common.Address, PairsInChain)
	pairsContracts := make([]*contract.UniswapV2Pair, PairsInChain)

	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)

	// Deploy router
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	routerAddress, _, _, err := contract.DeployUniswapRouter(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy UniswapRouter; %v", err)
	}

	// Deploy tokens
	for i := 0; i < TokensInChain; i++ {
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		name := fmt.Sprintf("Testing token %d", i)
		symbol := fmt.Sprintf("TOK%d", i)
		tokenAddresses[i], _, tokenContracts[i], err = contract.DeployERC20(txOpts, rpcClient, name, symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy ERC-20 token %d; %v", i, err)
		}
	}

	// Deploy pairs
	for i := 0; i < PairsInChain; i++ {
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		pairsAddresses[i], _, pairsContracts[i], err = contract.DeployUniswapV2Pair(txOpts, rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy Uniswap pair; %v", err)
		}
	}

	// wait until contracts are available on the chain
	err = WaitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Uniswap contract is deployed; %v", err)
	}

	// Mint tokens into pairs
	for i := 0; i < PairsInChain; i++ {
		tokenA, tokenB := tokenContracts[i], tokenContracts[i+1]
		tokenAAddress, tokenBAddress := tokenAddresses[i], tokenAddresses[i+1]
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		_, err = tokenA.Mint(txOpts, pairsAddresses[i], PairLiquidity)
		if err != nil {
			return nil, fmt.Errorf("failed to fund Uniswap pair; %v", err)
		}
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		_, err = tokenB.Mint(txOpts, pairsAddresses[i], PairLiquidity)
		if err != nil {
			return nil, fmt.Errorf("failed to fund Uniswap pair; %v", err)
		}

		// tokens addresses must be passed in ascending order into initializing method
		if bytes.Compare(tokenAAddress[:], tokenBAddress[:]) > 0 {
			tokenAAddress, tokenBAddress = tokenBAddress, tokenAAddress
		}
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		_, err = pairsContracts[i].Initialize(txOpts, tokenAAddress, tokenBAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Uniswap pair; %v", err)
		}
	}

	// Whitelist Uniswap router in the token (skip setting allowance by every user)
	for i := 0; i < TokensInChain; i++ {
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		_, err = tokenContracts[i].WhitelistSpender(txOpts, routerAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to whitelist Uniswap contract in the ERC-20 token %d; %v", i, err)
		}
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
	routerAbi, err := contract.UniswapRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the starting accounts will be available on the chain (and will be possible to call CreateUser)
	err = WaitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Uniswap contract is deployed; %v", err)
	}

	return &UniswapApplication{
		routerAbi:        routerAbi,
		startingAccounts: startingAccounts,
		routerAddress:    routerAddress,
		tokensAddresses:  tokenAddresses,
		pairsAddresses:   pairsAddresses,
		accountFactory:   accountFactory,
	}, nil
}

// UniswapApplication represents one application deployed to the network - an ERC-20 contract.
// Each created app should be used in a single thread only.
type UniswapApplication struct {
	routerAbi        *abi.ABI
	startingAccounts []*Account
	routerAddress    common.Address
	tokensAddresses  []common.Address
	pairsAddresses   []common.Address
	accountFactory   *AccountFactory
}

// CreateUser creates a new user for the app.
func (f *UniswapApplication) CreateUser(rpcClient rpc.RpcClient) (User, error) {
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
	token0Contract, err := contract.NewERC20(f.tokensAddresses[0], rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get token representation; %v", err)
	}
	tokenNContract, err := contract.NewERC20(f.tokensAddresses[TokensInChain-1], rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get token representation; %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(startingAccount.privateKey, startingAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(startingAccount.getNextNonce()))
	_, err = token0Contract.Mint(txOpts, workerAccount.address, WorkerInitialBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to mint ERC-20; %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(startingAccount.getNextNonce()))
	_, err = tokenNContract.Mint(txOpts, workerAccount.address, WorkerInitialBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to mint ERC-20; %v", err)
	}

	return &UniswapUser{
		routerAbi:               f.routerAbi,
		sender:                  workerAccount,
		routerAddress:           f.routerAddress,
		tokensAddresses:         f.tokensAddresses,
		pairsAddresses:          f.pairsAddresses,
		tokensAddressesReversed: reverseAddresses(f.tokensAddresses),
		pairsAddressesReversed:  reverseAddresses(f.pairsAddresses),
	}, nil
}

func (f *UniswapApplication) WaitUntilApplicationIsDeployed(rpcClient rpc.RpcClient) error {
	return waitUntilAllSentTxsAreOnChain(f.startingAccounts, rpcClient)
}

func (f *UniswapApplication) GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	routerContract, err := contract.NewUniswapRouter(f.routerAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get UniswapRouter representation; %v", err)
	}
	count, err := routerContract.GetCount(nil)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

// UniswapUser represents a user sending txs to swap ERC-20 tokens using Uniswap.
// A generator is supposed to be used in a single thread.
type UniswapUser struct {
	routerAbi               *abi.ABI
	sender                  *Account
	routerAddress           common.Address
	tokensAddresses         []common.Address
	pairsAddresses          []common.Address
	tokensAddressesReversed []common.Address
	pairsAddressesReversed  []common.Address
	sentTxs                 uint64
}

func (g *UniswapUser) GenerateTx(currentGasPrice *big.Int) (*types.Transaction, error) {
	var data []byte
	var err error

	// prepare tx data
	if rand.Intn(2) == 0 {
		// swap token1 for tokenN (forward)
		data, err = g.routerAbi.Pack("swapExactTokensForTokens", AmountSwapped, g.tokensAddresses, g.pairsAddresses)
	} else {
		// swap tokenN for token1 (backward)
		data, err = g.routerAbi.Pack("swapExactTokensForTokens", AmountSwapped, g.tokensAddressesReversed, g.pairsAddressesReversed)
	}
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	// swapExactTokensForTokens consumes 157571 for 2 tokens + cca 94314 for each additional token
	const gasLimit = 160_000 + (TokensInChain-2)*95000
	tx, err := createTx(g.sender, g.routerAddress, big.NewInt(0), data, currentGasPrice, gasLimit)
	if err == nil {
		atomic.AddUint64(&g.sentTxs, 1)
	}
	return tx, err
}

func (g *UniswapUser) GetSentTransactions() uint64 {
	return atomic.LoadUint64(&g.sentTxs)
}

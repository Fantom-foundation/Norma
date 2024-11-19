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
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Fantom-foundation/Norma/driver/rpc"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockgen -source context.go -destination context_mock.go -package app

// AppContext provides a context for the application to interact with the network.
// It includes the network client, the account paying for management tasks, and a helper
// contract used for on-chain operations. It also provides utility functions to interact
// with the network, such as deploying contracts, sending transactions, and waiting for
// receipts.
type AppContext interface {
	GetClient() rpc.RpcClient
	GetTreasure() *Account
	GetTransactOptions(account *Account) (*bind.TransactOpts, error)
	GetReceipt(txHash common.Hash) (*types.Receipt, error)
	Run(operation func(*bind.TransactOpts) (*types.Transaction, error)) (*types.Receipt, error)
	FundAccounts(accounts []common.Address, value *big.Int) error
	Close()
}

type RpcClientFactory interface {
	DialRandomRpc() (rpc.RpcClient, error)
}

func NewContext(factory RpcClientFactory, treasury *Account) (*appContext, error) {
	rpcClient, err := factory.DialRandomRpc()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to network: %w", err)
	}

	// Install the helper contract used by this contract for its operations.
	// Create a context to interact with the network.
	res := &appContext{
		rpcClient: rpcClient,
		treasury:  treasury,
	}

	// Install a helper contract on the network to perform operations.
	helper, receipt, err := DeployContract(res, contract.DeployHelper)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy helper contract: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("failed to deploy helper contract: transaction reverted")
	}

	res.helper = helper
	return res, nil
}

type appContext struct {
	rpcClient rpc.RpcClient    // < access to the network
	treasury  *Account         // < the account paying for management tasks
	helper    *contract.Helper // < a contract used for on-chain operations
}

func (c *appContext) Close() {
	c.rpcClient.Close()
}

func (c *appContext) GetClient() rpc.RpcClient {
	return c.rpcClient
}

func (c *appContext) GetTreasure() *Account {
	return c.treasury
}

// GetTransactOptions provides transaction options to be used to send a transaction
// with the given account. The options include the chain ID, a suggested gas price,
// the next free nonce of the given account, and a hard-coded gas limit of 1e6.
// The main purpose of this function is to provide a convenient way to collect all
// the necessary information required to create a transaction in one place.
func (c *appContext) GetTransactOptions(account *Account) (*bind.TransactOpts, error) {
	client := c.rpcClient

	ctxt := context.Background()
	chainId, err := client.ChainID(ctxt)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctxt)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price suggestion: %w", err)
	}

	nonce, err := client.NonceAt(ctxt, account.address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(account.privateKey, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction options: %w", err)
	}
	txOpts.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(2))
	txOpts.Nonce = big.NewInt(int64(nonce))
	return txOpts, nil
}

// GetReceipt waits for the receipt of the given transaction hash to be available.
// The function times out after 10 seconds.
func (c *appContext) GetReceipt(txHash common.Hash) (*types.Receipt, error) {
	client := c.rpcClient

	// Wait for the response with some exponential backoff.
	const maxDelay = 100 * time.Millisecond
	begin := time.Now()
	delay := time.Millisecond
	for time.Since(begin) < 10*time.Second {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if errors.Is(err, ethereum.NotFound) {
			time.Sleep(delay)
			delay = 2 * delay
			if delay > maxDelay {
				delay = maxDelay
			}
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
		}
		return receipt, nil
	}
	return nil, fmt.Errorf("failed to get transaction receipt: timeout")
}

// Apply sends a transaction to the network using the network's validator account
// and waits for the transaction to be processed. The resulting receipt is returned.
func (c *appContext) Run(
	operation func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	txOpts, err := c.GetTransactOptions(c.treasury)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %w", err)
	}
	transaction, err := operation(txOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return c.GetReceipt(transaction.Hash())
}

// FundAccounts transfers the given amount of funds from the treasure to each of the
// given accounts.
func (c *appContext) FundAccounts(accounts []common.Address, value *big.Int) error {
	// Group funding requests in batches to avoid making individual transactions
	// to big fo a single block.
	const batchSize = 128
	batches := make([][]common.Address, 0)
	for i := 0; i < len(accounts); i += batchSize {
		batches = append(batches, accounts[i:min(i+batchSize, len(accounts))])
	}

	// Send one transaction per batch of accounts.
	opts, err := c.GetTransactOptions(c.GetTreasure())
	if err != nil {
		return fmt.Errorf("failed to get transaction options: %w", err)
	}
	txs := make([]*types.Transaction, 0, len(batches))
	for _, batch := range batches {
		opts.Value = new(big.Int).Mul(value, big.NewInt(int64(len(batch))))
		tx, err := c.helper.Distribute(opts, batch)
		if err != nil {
			return fmt.Errorf("failed to distribute funds: %w", err)
		}
		txs = append(txs, tx)
		opts.Nonce.Add(opts.Nonce, big.NewInt(1))
	}

	// Wait for all the transactions to be completed.
	for _, tx := range txs {
		receipt, err := c.GetReceipt(tx.Hash())
		if err != nil {
			return fmt.Errorf("failed to get receipt: %w", err)
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("failed to distribute funds: transaction reverted")
		}
	}
	return nil
}

// DeployContract is a utility function handling the deployment of a contract on the network.
// The contract is deployed with by the network's treasure account. The function returns the
// deployed contract instance and the transaction receipt.
func DeployContract[T any](c AppContext, deploy contractDeployer[T]) (*T, *types.Receipt, error) {
	client := c.GetClient()

	transactOptions, err := c.GetTransactOptions(c.GetTreasure())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transaction options: %w", err)
	}

	_, transaction, contract, err := deploy(transactOptions, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy contract: %w", err)
	}

	receipt, err := c.GetReceipt(transaction.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt: %w", err)
	}
	return contract, receipt, nil
}

// contractDeployer is the type of the deployment functions generated by abigen.
type contractDeployer[T any] func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *T, error)

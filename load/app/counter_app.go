package app

import (
	"context"
	"fmt"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sync/atomic"
)

// NewCounterApplication deploys a Counter contract to the chain.
// The Counter contract is a simple contract sustaining an integer value, to be incremented by sent txs.
// It allows to easily test the tx generating, as reading the contract field provides the amount of applied contract calls.
func NewCounterApplication(rpcClient RpcClient, primaryAccount *Account) (*CounterApplication, error) {
	// get price of gas from the network
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price; %v", err)
	}
	// use greater gas price
	gasPrice.Mul(gasPrice, big.NewInt(4)) // higher coefficient for the contract deploy

	// Deploy the Counter contract to be used by tx generators
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = gasPrice
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	contractAddress, _, _, err := contract.DeployCounter(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Counter contract; %v", err)
	}

	// parse ABI for generating txs data
	parsedAbi, err := contract.CounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	err = waitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Counter contract is deployed; %v", err)
	}

	return &CounterApplication{
		abi:             parsedAbi,
		primaryAccount:  primaryAccount,
		contractAddress: contractAddress,
	}, nil
}

// CounterApplication represents a simple on-chain Counter incremented by sent transactions.
// A factory represents one deployed Counter contract, incremented by all its generators.
// While the factory is thread-safe, each created generator should be used in a single thread only.
type CounterApplication struct {
	abi             *abi.ABI
	primaryAccount  *Account
	contractAddress common.Address
	sentTxs         uint64
}

// CreateGenerator creates a new transaction generator for the app.
func (f *CounterApplication) CreateGenerator(rpcClient RpcClient) (TransactionGenerator, error) {

	// generate a new account for each worker - avoid account nonces related bottlenecks
	workerAccount, err := GenerateAccount(f.primaryAccount.chainID.Int64())
	if err != nil {
		return nil, err
	}

	// get price of gas from the network
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price; %v", err)
	}
	priorityGasPrice := big.NewInt(0)
	regularGasPrice := big.NewInt(0)
	priorityGasPrice.Mul(gasPrice, big.NewInt(4)) // greater gas price for init
	regularGasPrice.Mul(gasPrice, big.NewInt(2))  // lower gas price for regular txs

	// transfer budget (10 FTM) to worker's account - finances to cover transaction fees
	workerBudget := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1_000000000000000000))
	err = transferValue(rpcClient, f.primaryAccount, workerAccount.address, workerBudget, priorityGasPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to tranfer from primary account to app account: %v", err)
	}

	gen := &CounterGenerator{
		abi:      f.abi,
		sender:   workerAccount,
		gasPrice: regularGasPrice,
		sentTxs:  &f.sentTxs,
		contract: f.contractAddress,
	}
	return gen, nil
}

func (f *CounterApplication) WaitUntilApplicationIsDeployed(rpcClient RpcClient) error {
	return waitUntilAccountNonceIs(f.primaryAccount.address, f.primaryAccount.getCurrentNonce(), rpcClient)
}

func (f *CounterApplication) GetTransactionCounts(rpcClient RpcClient) (TransactionCounts, error) {
	// get a representation of the deployed contract
	counterContract, err := contract.NewCounter(f.contractAddress, rpcClient)
	if err != nil {
		return TransactionCounts{}, fmt.Errorf("failed to get Counter contract representation; %v", err)
	}
	count, err := counterContract.GetCount(nil)
	if err != nil {
		return TransactionCounts{}, err
	}
	return TransactionCounts{
		ReceivedTxs: count.Uint64(),
		SentTxs:     atomic.LoadUint64(&f.sentTxs),
	}, nil
}

// CounterGenerator is a txs generator incrementing trivial Counter contract.
// A generator is supposed to be used in a single thread.
type CounterGenerator struct {
	abi      *abi.ABI
	sender   *Account
	gasPrice *big.Int
	sentTxs  *uint64
	contract common.Address
}

func (g *CounterGenerator) GenerateTx() (*types.Transaction, error) {
	// prepare tx data
	data, err := g.abi.Pack("incrementCounter")
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	const gasLimit = 50000 // IncrementCounter method call takes 43426 of gas
	tx, err := createTx(g.sender, g.contract, big.NewInt(0), data, g.gasPrice, gasLimit)
	if err == nil {
		atomic.AddUint64(g.sentTxs, 1)
	}
	return tx, err
}

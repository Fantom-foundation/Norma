package app

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"sync/atomic"

	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewERC20Application deploys a new ERC-20 dapp to the chain.
// The ERC20 contract is a contract sustaining balances of the token for individual owner addresses.
func NewERC20Application(rpcClient RpcClient, primaryAccount *Account) (*ERC20Application, error) {
	// get price of gas from the network
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price; %v", err)
	}
	// use greater gas price
	gasPrice.Mul(gasPrice, big.NewInt(4)) // higher coefficient for the contract deploy

	// Deploy the ERC20 contract to be used by generators created using the factory
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = gasPrice
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	contractAddress, _, _, err := contract.DeployERC20(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ERC20 contract; %v", err)
	}
	recipients, err := generateRecipientsAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipients addresses; %v", err)
	}

	// parse ABI for generating txs data
	parsedAbi, err := contract.ERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	err = waitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the ERC20 contract is deployed; %v", err)
	}

	return &ERC20Application{
		abi:             parsedAbi,
		primaryAccount:  primaryAccount,
		contractAddress: contractAddress,
		recipients:      recipients,
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
	abi             *abi.ABI
	primaryAccount  *Account
	contractAddress common.Address
	recipients      []common.Address
	numAccounts     int
}

// CreateGenerator creates a new transaction generator for the app.
func (f *ERC20Application) CreateGenerator(rpcClient RpcClient) (TransactionGenerator, error) {
	// get price of gas from the network
	regularGasPrice, err := getGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// generate a new account for each worker - avoid account nonces related bottlenecks
	id := f.numAccounts
	f.numAccounts++
	workerAccount, err := GenerateAndFundAccount(f.primaryAccount, rpcClient, regularGasPrice, id)
	if err != nil {
		return nil, err
	}

	// mint ERC-20 tokens for the worker account - tokens to be transferred in the transactions
	erc20Contract, err := contract.NewERC20(f.contractAddress, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get ERC20 contract representation; %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(f.primaryAccount.privateKey, f.primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(f.primaryAccount.getNextNonce()))
	_, err = erc20Contract.Mint(txOpts, workerAccount.address, big.NewInt(1_000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("failed to mint ERC-20; %v", err)
	}

	return &ERC20Generator{
		abi:        f.abi,
		sender:     workerAccount,
		gasPrice:   regularGasPrice,
		contract:   f.contractAddress,
		recipients: f.recipients,
	}, nil
}

func (f *ERC20Application) WaitUntilApplicationIsDeployed(rpcClient RpcClient) error {
	return waitUntilAccountNonceIs(f.primaryAccount.address, f.primaryAccount.getCurrentNonce(), rpcClient)
}

func (f *ERC20Application) GetReceivedTransations(rpcClient RpcClient) (uint64, error) {
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

// ERC20Generator is a txs app transferring ERC20 tokens.
// A generator is supposed to be used in a single thread.
type ERC20Generator struct {
	abi        *abi.ABI
	sender     *Account
	gasPrice   *big.Int
	contract   common.Address
	recipients []common.Address
	sentTxs    uint64
}

func (g *ERC20Generator) GenerateTx() (*types.Transaction, error) {
	// choose random recipient
	recipient := g.recipients[rand.Intn(len(g.recipients))]

	// prepare tx data
	data, err := g.abi.Pack("transfer", recipient, big.NewInt(1))
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	const gasLimit = 52000 // Transfer method call takes 51349 of gas
	tx, err := createTx(g.sender, g.contract, big.NewInt(0), data, g.gasPrice, gasLimit)
	if err == nil {
		atomic.AddUint64(&g.sentTxs, 1)
	}
	return tx, err
}

func (g *ERC20Generator) GetSentTransactions() uint64 {
	return atomic.LoadUint64(&g.sentTxs)
}

package generator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sync/atomic"
)

// NewCounterGeneratorFactory provides a factory of tx generators incrementing one deployed Counter contract.
// The Counter contract is a simple contract sustaining an integer value, to be incremented by sent txs.
// It allows to easily test the tx generating, as reading the contract field provides the amount of applied contract calls.
func NewCounterGeneratorFactory(rpcUrl URL, primaryPrivateKey *ecdsa.PrivateKey, chainID *big.Int) (*CounterGeneratorFactory, error) {
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryPrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	rpcClient, err := ethclient.Dial(string(rpcUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC to initialize the Counter; %v", err)
	}
	defer rpcClient.Close()
	// Deploy the Counter contract to be used by generators created using the factory
	contractAddress, _, _, err := abi.DeployCounter(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Counter contract; %v", err)
	}
	err = waitUntilContractStartExisting(contractAddress, rpcClient)
	if err != nil {
		return nil, err
	}
	return &CounterGeneratorFactory{
		rpcUrl:            rpcUrl,
		primaryPrivateKey: primaryPrivateKey,
		chainID:           chainID,
		contractAddress:   contractAddress,
	}, nil
}

// CounterGeneratorFactory is a factory of tx generators incrementing one deployed Counter contract.
// A factory represents one deployed Counter contract, incremented by all its generators.
// While the factory is thread-safe, each created generator should be used in a single thread only.
type CounterGeneratorFactory struct {
	rpcUrl            URL
	primaryPrivateKey *ecdsa.PrivateKey
	chainID           *big.Int
	contractAddress   common.Address
	sentTxs           uint64
}

// Create a new generator to be used by one worker thread.
func (f *CounterGeneratorFactory) Create() (TransactionGenerator, error) {
	// generate a new account for each worker - avoid account nonces related bottlenecks
	address, privateKey, err := generateAddress()
	if err != nil {
		return nil, err
	}

	rpcClient, err := ethclient.Dial(string(f.rpcUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC to create a tx generator; %v", err)
	}

	err = f.primeGeneratorAccount(rpcClient, address)
	if err != nil {
		return nil, fmt.Errorf("account priming failed; %v", err)
	}

	// get a representation of the deployed contract, bound to worker's rpcClient
	counterContract, err := abi.NewCounter(f.contractAddress, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get Counter contract representation; %v", err)
	}

	// get nonce of the worker account
	nonce, err := rpcClient.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce; %v", err)
	}

	// get price of gas
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price; %v", err)
	}

	// prepare options to generate transactions with
	txOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, f.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts; %v", err)
	}
	// adjust txOpts for the runtime to avoid slow loading of gasPrice/nonce from RPC for each tx
	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.GasLimit = 50000 // IncrementCounter method call takes 43426 of gas
	// use static gasPrice - multiply recommended gas price by two to have some levy
	// for the cases the min gas price is increased by the network,
	// and we do not have to query the gas price before every transaction.
	txOpts.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(int64(2)))

	gen := &CounterGenerator{
		rpcClient:       rpcClient,
		txOpts:          txOpts,
		counterContract: counterContract,
		sentTxs:         &f.sentTxs,
	}
	return gen, nil
}

func (f *CounterGeneratorFactory) primeGeneratorAccount(rpcClient *ethclient.Client, address common.Address) error {
	primaryAddress := crypto.PubkeyToAddress(f.primaryPrivateKey.PublicKey)
	primaryNonce, err := rpcClient.PendingNonceAt(context.Background(), primaryAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce of primary account: %s", err)
	}

	// transfer budget (10 FTM) to worker's account - finances to cover transaction fees
	workerBudget := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1_000000000000000000))
	err = transferValue(rpcClient, f.chainID, f.primaryPrivateKey, address, workerBudget, primaryNonce)
	if err != nil {
		return fmt.Errorf("failed to tranfer from primary account to generator account: %v", err)
	}

	// wait until required updates are on the chain
	err = waitUntilAccountNonceIsAtLeast(primaryAddress, primaryNonce+1, rpcClient)
	if err != nil {
		return fmt.Errorf("waiting for chain changes failed; %v", err)
	}
	return nil
}

// GetAmountOfSentTxs provides the amount of txs send from all generators of the factory
func (f *CounterGeneratorFactory) GetAmountOfSentTxs() uint64 {
	return atomic.LoadUint64(&f.sentTxs)
}

// GetAmountOfReceivedTxs provides the amount of relevant txs applied to the chain state
// This is obtained as the counter value in the Counter contract.
func (f *CounterGeneratorFactory) GetAmountOfReceivedTxs() (uint64, error) {
	rpcClient, err := ethclient.Dial(string(f.rpcUrl))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to RPC to get amount of txs on chain; %v", err)
	}
	defer rpcClient.Close()
	// get a representation of the deployed contract, bound to worker's rpcClient
	counterContract, err := abi.NewCounter(f.contractAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get Counter contract representation; %v", err)
	}
	count, err := counterContract.GetCount(nil)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

// CounterGenerator is a txs generator incrementing trivial Counter contract.
// A generator is supposed to be used in a single thread.
type CounterGenerator struct {
	rpcClient       *ethclient.Client
	txOpts          *bind.TransactOpts
	counterContract *abi.Counter
	sentTxs         *uint64
}

func (g *CounterGenerator) SendTx() error {
	_, err := g.counterContract.IncrementCounter(g.txOpts)
	if err == nil {
		g.txOpts.Nonce.Add(g.txOpts.Nonce, big.NewInt(1))
		atomic.AddUint64(g.sentTxs, 1)
	}
	return err
}

func (g *CounterGenerator) Close() error {
	g.rpcClient.Close()
	return nil
}

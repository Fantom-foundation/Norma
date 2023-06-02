package generator

import (
	"context"
	"crypto/ecdsa"
	crand "crypto/rand"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"math/rand"
	"sync/atomic"
)

// NewERC20GeneratorFactory provides a factory of tx generators transferring ERC20 tokens.
// The ERC20 contract is a contract sustaining balances of the token for individual owner addresses.
func NewERC20GeneratorFactory(rpcUrl driver.URL, primaryPrivateKey *ecdsa.PrivateKey, chainID *big.Int) (*ERC20GeneratorFactory, error) {
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryPrivateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	rpcClient, err := ethclient.Dial(string(rpcUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC to initialize the ERC20; %v", err)
	}
	defer rpcClient.Close()
	// Deploy the ERC20 contract to be used by generators created using the factory
	contractAddress, _, _, err := abi.DeployERC20(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ERC20 contract; %v", err)
	}
	recipients, err := generateRecipientsAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to generate recipients addresses; %v", err)
	}
	err = waitUntilContractStartExisting(contractAddress, rpcClient)
	if err != nil {
		return nil, err
	}
	return &ERC20GeneratorFactory{
		rpcUrl:            rpcUrl,
		primaryPrivateKey: primaryPrivateKey,
		chainID:           chainID,
		contractAddress:   contractAddress,
		recipients:        recipients,
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

// ERC20GeneratorFactory is a factory of tx generators transferring tokens of one ERC20 contract.
// While the factory is thread-safe, each created generator should be used in a single thread only.
type ERC20GeneratorFactory struct {
	rpcUrl            driver.URL
	primaryPrivateKey *ecdsa.PrivateKey
	chainID           *big.Int
	contractAddress   common.Address
	sentTxs           uint64
	recipients        []common.Address
}

// Create a new generator to be used by one worker thread.
func (f *ERC20GeneratorFactory) Create() (TransactionGenerator, error) {
	// generate a new account for each worker - avoid account nonces related bottlenecks
	address, privateKey, err := generateAddress()
	if err != nil {
		return nil, err
	}

	rpcClient, err := ethclient.Dial(string(f.rpcUrl))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC to create a tx generator; %v", err)
	}

	// get a representation of the deployed contract, bound to worker's rpcClient
	erc20Contract, err := abi.NewERC20(f.contractAddress, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get ERC20 contract representation; %v", err)
	}

	err = f.primeGeneratorAccount(rpcClient, address, erc20Contract)
	if err != nil {
		return nil, fmt.Errorf("account priming failed; %v", err)
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
	txOpts.GasLimit = 52000    // Transfer method call takes 51349 of gas
	txOpts.GasPrice = gasPrice // use static gasPrice

	gen := &ERC20Generator{
		rpcClient:     rpcClient,
		txOpts:        txOpts,
		erc20Contract: erc20Contract,
		sentTxs:       &f.sentTxs,
		recipients:    f.recipients,
	}
	return gen, nil
}

func (f *ERC20GeneratorFactory) primeGeneratorAccount(rpcClient *ethclient.Client, workerAddress common.Address, erc20Contract *abi.ERC20) error {
	primaryAddress := crypto.PubkeyToAddress(f.primaryPrivateKey.PublicKey)
	primaryNonce, err := rpcClient.PendingNonceAt(context.Background(), primaryAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce of primary account: %s", err)
	}

	// transfer budget (10 FTM) to worker's account - finances to cover transaction fees
	workerBudget := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1_000000000000000000))
	err = transferValue(rpcClient, f.chainID, f.primaryPrivateKey, workerAddress, workerBudget, primaryNonce)
	if err != nil {
		return fmt.Errorf("failed to tranfer from primary account to generator account: %v", err)
	}

	// prepare options to generate transactions with
	txOpts, err := bind.NewKeyedTransactorWithChainID(f.primaryPrivateKey, f.chainID)
	if err != nil {
		return fmt.Errorf("failed to create txOpts; %v", err)
	}

	// mint ERC-20 tokens for the worker account - tokens to be transferred in the transactions
	txOpts.Nonce = big.NewInt(int64(primaryNonce + 1))
	_, err = erc20Contract.Mint(txOpts, workerAddress, big.NewInt(1_000000000000000000))
	if err != nil {
		return fmt.Errorf("failed to mint ERC-20; %v", err)
	}

	// wait until required updates are on the chain
	err = waitUntilAccountNonceIsAtLeast(primaryAddress, primaryNonce+2, rpcClient)
	if err != nil {
		return fmt.Errorf("waiting for chain changes failed; %v", err)
	}
	return nil
}

// GetAmountOfSentTxs provides the amount of txs send from all generators of the factory
func (f *ERC20GeneratorFactory) GetAmountOfSentTxs() uint64 {
	return atomic.LoadUint64(&f.sentTxs)
}

// GetAmountOfReceivedTxs provides the amount of relevant txs applied to the chain state
// This is obtained as the total ERC20 balance of all recipient accounts.
func (f *ERC20GeneratorFactory) GetAmountOfReceivedTxs() (uint64, error) {
	rpcClient, err := ethclient.Dial(string(f.rpcUrl))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to RPC to get amount of txs on chain; %v", err)
	}
	defer rpcClient.Close()
	// get a representation of the deployed contract, bound to worker's rpcClient
	ERC20Contract, err := abi.NewERC20(f.contractAddress, rpcClient)
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

// ERC20Generator is a txs generator transferring ERC20 tokens.
// A generator is supposed to be used in a single thread.
type ERC20Generator struct {
	rpcClient     *ethclient.Client
	txOpts        *bind.TransactOpts
	erc20Contract *abi.ERC20
	sentTxs       *uint64
	recipients    []common.Address
}

func (g *ERC20Generator) SendTx() error {
	// choose random recipient
	recipient := g.recipients[rand.Intn(len(g.recipients))]
	// transfer tokens to the recipient
	_, err := g.erc20Contract.Transfer(g.txOpts, recipient, big.NewInt(1))
	if err == nil {
		g.txOpts.Nonce.Add(g.txOpts.Nonce, big.NewInt(1))
		atomic.AddUint64(g.sentTxs, 1)
	}
	return err
}

func (g *ERC20Generator) Close() error {
	g.rpcClient.Close()
	return nil
}

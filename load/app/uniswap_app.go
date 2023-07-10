package app

import (
	"fmt"
	"math/big"
	"sync/atomic"

	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NewUniswapApplication deploys a new Uniswap dapp to the chain.
func NewUniswapApplication(rpcClient RpcClient, primaryAccount *Account, numUsers int) (*UniswapApplication, error) {
	// get price of gas from the network
	regularGasPrice, err := getGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// Deploy the Uniswap contract to be used by generators created using the factory
	txOpts, err := bind.NewKeyedTransactorWithChainID(primaryAccount.privateKey, primaryAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts for contract deploy; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	token0Address, _, token0, err := contract.DeployERC20(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy first ERC-20 contract; %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	token1Address, _, token1, err := contract.DeployERC20(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy second ERC-20 contract; %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	pairAddress, _, pair, err := contract.DeployUniswapV2Pair(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Uniswap pair contract; %v", err)
	}

	// wait until contracts will be available on the chain
	err = waitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Uniswap contract is deployed; %v", err)
	}

	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	_, err = token0.Mint(txOpts, pairAddress, big.NewInt(1_000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("failed to fund Uniswap pair with token0; %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	_, err = token1.Mint(txOpts, pairAddress, big.NewInt(1_000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("failed to fund Uniswap pair with token1; %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	_, err = pair.Initialize(txOpts, token0Address, token1Address)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Uniswap pair contract; %v", err)
	}
	txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
	_, err = pair.Sync(txOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to sync Uniswap pair contract; %v", err)
	}

	// deploying too many generators from one account leads to excessive gasPrice growth - we
	// need to spread the initialization in between multiple startingAccounts
	startingAccounts, err := generateStartingAccounts(rpcClient, primaryAccount, numUsers, regularGasPrice)
	if err != nil {
		return nil, err
	}

	// parse ABI for generating txs data
	pairAbi, err := contract.UniswapV2PairMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the contract will be available on the chain (and will be possible to call CreateGenerator)
	err = waitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Uniswap contract is deployed; %v", err)
	}

	return &UniswapApplication{
		pairAbi:          pairAbi,
		startingAccounts: startingAccounts,
		pairAddress:      pairAddress,
		token0Address:    token0Address,
		token1Address:    token1Address,
	}, nil
}

// UniswapApplication represents one application deployed to the network - an ERC-20 contract.
// Each created app should be used in a single thread only.
type UniswapApplication struct {
	pairAbi          *abi.ABI
	startingAccounts []*Account
	pairAddress      common.Address
	token0Address    common.Address
	token1Address    common.Address
	recipients       []common.Address
	numAccounts      int64
}

// CreateUser creates a new user for the app.
func (f *UniswapApplication) CreateUser(rpcClient RpcClient) (User, error) {
	// get price of gas from the network
	regularGasPrice, err := getGasPrice(rpcClient)
	if err != nil {
		return nil, err
	}

	// generate a new account for each worker - avoid account nonces related bottlenecks
	id := atomic.AddInt64(&f.numAccounts, 1)
	startingAccount := f.startingAccounts[id%int64(len(f.startingAccounts))]
	workerAccount, err := GenerateAndFundAccount(startingAccount, rpcClient, getPriorityGasPrice(regularGasPrice), int(id), 1000)
	if err != nil {
		return nil, err
	}

	// mint ERC-20 tokens for the worker account - tokens to be transferred in the transactions
	token0Contract, err := contract.NewERC20(f.token0Address, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get Uniswap contract representation; %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(startingAccount.privateKey, startingAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(startingAccount.getNextNonce()))
	_, err = token0Contract.Mint(txOpts, workerAccount.address, big.NewInt(1_000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("failed to mint ERC-20; %v", err)
	}

	return &UniswapUser{
		abi:      f.pairAbi,
		sender:   workerAccount,
		gasPrice: regularGasPrice,
		contract: f.pairAddress,
	}, nil
}

func (f *UniswapApplication) WaitUntilApplicationIsDeployed(rpcClient RpcClient) error {
	return waitUntilAllSentTxsAreOnChain(f.startingAccounts, rpcClient)
}

func (f *UniswapApplication) GetReceivedTransactions(rpcClient RpcClient) (uint64, error) {
	return 0, nil // TODO
}

// UniswapUser represents a user sending txs to swap ERC-20 tokens using Uniswap.
// A generator is supposed to be used in a single thread.
type UniswapUser struct {
	abi      *abi.ABI
	sender   *Account
	gasPrice *big.Int
	contract common.Address
	sentTxs  uint64
}

func (g *UniswapUser) GenerateTx() (*types.Transaction, error) {

	// prepare tx data
	amount0Out := big.NewInt(0)
	amount1Out := big.NewInt(10)
	data, err := g.abi.Pack("swap", amount0Out, amount1Out, g.sender.address, []byte{0xAA})
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

func (g *UniswapUser) GetSentTransactions() uint64 {
	return atomic.LoadUint64(&g.sentTxs)
}

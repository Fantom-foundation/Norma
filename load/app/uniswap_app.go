package app

import (
	"bytes"
	"fmt"
	contract "github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"math/rand"
	"sync/atomic"
)

const TokensInChain = 3
const PairsInChain = TokensInChain - 1

var AmountSwapped = big.NewInt(100) // swapped in one tx
var WorkerInitialBalance = big.NewInt(1_000000000000000000)
var PairLiquidity = big.NewInt(0).Mul(big.NewInt(10_000_000_000), big.NewInt(1_000000000000000000))

// NewUniswapApplication deploys a new Uniswap dapp to the chain.
func NewUniswapApplication(rpcClient RpcClient, primaryAccount *Account, numUsers int) (*UniswapApplication, error) {
	// get price of gas from the network
	regularGasPrice, err := getGasPrice(rpcClient)
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
	routerAddress, _, _, err := contract.DeploySimpleRouter(txOpts, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy SimpleRouter; %v", err)
	}

	// Deploy tokens
	for i := 0; i < TokensInChain; i++ {
		txOpts.Nonce = big.NewInt(int64(primaryAccount.getNextNonce()))
		tokenAddresses[i], _, tokenContracts[i], err = contract.DeployERC20(txOpts, rpcClient)
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
	err = waitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
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
		fmt.Printf("initilize pair %x to connect %x with %x\n", pairsAddresses[i], tokenAAddress, tokenBAddress)
		_, err = pairsContracts[i].Initialize(txOpts, tokenAAddress, tokenBAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Uniswap pair; %v", err)
		}
	}

	// deploying too many generators from one account leads to excessive gasPrice growth - we
	// need to spread the initialization in between multiple startingAccounts
	startingAccounts, err := generateStartingAccounts(rpcClient, primaryAccount, numUsers, regularGasPrice)
	if err != nil {
		return nil, err
	}

	// parse ABI for generating txs data
	routerAbi, err := contract.SimpleRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// wait until the starting accounts will be available on the chain (and will be possible to call CreateUser)
	err = waitUntilAccountNonceIs(primaryAccount.address, primaryAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until the Uniswap contract is deployed; %v", err)
	}

	return &UniswapApplication{
		routerAbi:        routerAbi,
		startingAccounts: startingAccounts,
		routerAddress:    routerAddress,
		tokensAddresses:  tokenAddresses,
		pairsAddresses:   pairsAddresses,
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

	// TEMPORARY TEST
	// wait until funds are available
	err = waitUntilAccountNonceIs(startingAccount.address, startingAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until funded; %v", err)
	}

	for i := 0; i < TokensInChain; i++ {
		tokenContract, err := contract.NewERC20(f.tokensAddresses[i], rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to get token representation; %v", err)
		}
		balance, err := tokenContract.BalanceOf(nil, workerAccount.address)
		fmt.Printf("token %d (%x) balance: %s, %s\n", i, f.tokensAddresses[i], balance.String(), err)
	}

	routerContract, err := contract.NewSimpleRouter(f.routerAddress, rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to crete SimpleRouter; %v", err)
	}
	txOpts, err = bind.NewKeyedTransactorWithChainID(workerAccount.privateKey, workerAccount.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = getPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(workerAccount.getNextNonce()))

	fmt.Printf("tokens param: %v\n", f.tokensAddresses)
	fmt.Printf("pairs param: %v\n", f.pairsAddresses)

	tx, err := routerContract.SwapExactTokensForTokens(txOpts, AmountSwapped, f.tokensAddresses, f.pairsAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to SwapExactTokensForTokens; %v", err)
	}
	fmt.Printf("SwapExactTokensForTokens successful, gas: %d\n", tx.Gas())

	err = waitUntilAccountNonceIs(workerAccount.address, workerAccount.getCurrentNonce(), rpcClient)
	if err != nil {
		return nil, fmt.Errorf("failed to wait until funded; %v", err)
	}
	for i := 0; i < TokensInChain; i++ {
		tokenContract, err := contract.NewERC20(f.tokensAddresses[i], rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to get token representation; %v", err)
		}
		balance, err := tokenContract.BalanceOf(nil, workerAccount.address)
		fmt.Printf("token %d (%x) balance: %s, %s\n", i, f.tokensAddresses[i], balance.String(), err)
	}

	return &UniswapUser{
		routerAbi:               f.routerAbi,
		sender:                  workerAccount,
		gasPrice:                regularGasPrice,
		routerAddress:           f.routerAddress,
		tokensAddresses:         f.tokensAddresses,
		pairsAddresses:          f.pairsAddresses,
		tokensAddressesReversed: reverseAddresses(f.tokensAddresses),
		pairsAddressesReversed:  reverseAddresses(f.pairsAddresses),
	}, nil
}

func (f *UniswapApplication) WaitUntilApplicationIsDeployed(rpcClient RpcClient) error {
	return waitUntilAllSentTxsAreOnChain(f.startingAccounts, rpcClient)
}

func (f *UniswapApplication) GetReceivedTransactions(rpcClient RpcClient) (uint64, error) {
	// get a representation of the deployed contract
	routerContract, err := contract.NewSimpleRouter(f.routerAddress, rpcClient)
	if err != nil {
		return 0, fmt.Errorf("failed to get SimpleRouter representation; %v", err)
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
	gasPrice                *big.Int
	routerAddress           common.Address
	tokensAddresses         []common.Address
	pairsAddresses          []common.Address
	tokensAddressesReversed []common.Address
	pairsAddressesReversed  []common.Address
	sentTxs                 uint64
}

func (g *UniswapUser) GenerateTx() (*types.Transaction, error) {
	var data []byte
	var err error

	// prepare tx data
	if rand.Intn(2) == 0 {
		// swap from token1 to tokenN
		data, err = g.routerAbi.Pack("swapExactTokensForTokens", AmountSwapped, g.tokensAddresses, g.pairsAddresses)
	} else {
		// swap from tokenN to token1
		data, err = g.routerAbi.Pack("swapExactTokensForTokens", AmountSwapped, g.tokensAddressesReversed, g.pairsAddressesReversed)
	}
	if err != nil || data == nil {
		return nil, fmt.Errorf("failed to prepare tx data; %v", err)
	}

	// prepare tx
	const gasLimit = 500000 // swapExactTokensForTokens consumes 157571 with 1 pair, 251884 with 2 pairs
	tx, err := createTx(g.sender, g.routerAddress, big.NewInt(0), data, g.gasPrice, gasLimit)
	if err == nil {
		atomic.AddUint64(&g.sentTxs, 1)
	}
	return tx, err
}

func (g *UniswapUser) GetSentTransactions() uint64 {
	return atomic.LoadUint64(&g.sentTxs)
}

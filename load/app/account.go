package app

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents an account from which we can send transactions.
// It sustains the nonce value - it allows multiple generators which use one Account
// to produce multiple txs in one block.
type Account struct {
	id         int
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
	nonce      uint64
}

// NewAccount creates an Account instance from the provided private key
func NewAccount(id int, privateKeyHex string, chainID int64) (*Account, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Account{
		id:         id,
		privateKey: privateKey,
		address:    address,
		chainID:    big.NewInt(chainID),
		nonce:      0,
	}, nil
}

// GenerateAccount creates a new Account from a random private key
func GenerateAccount(id int, chainID int64) (*Account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Account{
		id:         id,
		privateKey: privateKey,
		address:    address,
		chainID:    big.NewInt(chainID),
		nonce:      0,
	}, nil
}

// GenerateAndFundAccount creates a new Account with a random private key and transfer finances to cover txs fees
func GenerateAndFundAccount(sourceAccount *Account, rpcClient rpc.RpcClient, regularGasPrice *big.Int, accountId int, endowment int64) (*Account, error) {
	priorityGasPrice := getPriorityGasPrice(regularGasPrice)
	account, err := GenerateAccount(accountId, sourceAccount.chainID.Int64())
	if err != nil {
		return nil, fmt.Errorf("failed to generate account; %v", err)
	}
	// transfers (tokens) FTM to the new account
	value := big.NewInt(0).Mul(big.NewInt(endowment), big.NewInt(1_000_000_000_000_000_000)) // FTM to wei
	if err := transferValue(rpcClient, sourceAccount, account.address, value, priorityGasPrice); err != nil {
		return nil, fmt.Errorf("failed to transfer (value: %s, gasPrice: %s): %v", value, priorityGasPrice, err)
	}
	return account, nil
}

// getNextNonce provides a nonce to be used for next transactions sent using this account
func (a *Account) getNextNonce() uint64 {
	current := atomic.AddUint64(&a.nonce, 1)
	return current - 1
}

func (a *Account) getCurrentNonce() uint64 {
	return atomic.LoadUint64(&a.nonce)
}

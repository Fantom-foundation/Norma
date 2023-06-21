package app

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"sync/atomic"
)

// Account represents an account from which we can send transactions.
// It sustains the nonce value - it allows multiple generators which use one Account
// to produce multiple txs in one block.
type Account struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
	nonce      uint64
}

// NewAccount creates an Account instance from the provided private key
func NewAccount(privateKeyHex string, chainID int64) (*Account, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Account{
		privateKey: privateKey,
		address:    address,
		chainID:    big.NewInt(chainID),
		nonce:      0,
	}, nil
}

// GenerateAccount creates a new Account from a random private key
func GenerateAccount(chainID int64) (*Account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Account{
		privateKey: privateKey,
		address:    address,
		chainID:    big.NewInt(chainID),
		nonce:      0,
	}, nil
}

// GenerateAndFundAccount creates a new Account with a random private key and transfer finances to cover txs fees
func GenerateAndFundAccount(sourceAccount *Account, rpcClient RpcClient, gasPrice *big.Int) (*Account, error) {
	account, err := GenerateAccount(sourceAccount.chainID.Int64())
	if err != nil {
		return nil, fmt.Errorf("failed to generate account; %v", err)
	}
	// transfers 1000 FTM to the new account - finances to cover transaction fees
	workerBudget := big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(1_000000000000000000))
	if err := transferValue(rpcClient, sourceAccount, account.address, workerBudget, gasPrice); err != nil {
		return nil, fmt.Errorf("failed to fund account: %v", err)
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

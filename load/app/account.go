package app

import (
	"crypto/ecdsa"
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

// getNextNonce provides a nonce to be used for next transactions sent using this account
func (a *Account) getNextNonce() uint64 {
	current := atomic.AddUint64(&a.nonce, 1)
	return current - 1
}

func (a *Account) getCurrentNonce() uint64 {
	return atomic.LoadUint64(&a.nonce)
}

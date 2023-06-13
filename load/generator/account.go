package generator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
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

// NewAccount provides a new Account instance
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

// getNextNonce provides a nonce to be used for next transactions sent using this account
func (a *Account) getNextNonce(rpcClient *ethclient.Client) (uint64, error) {
	current := a.nonce
	if current == 0 {
		var err error
		current, err = rpcClient.PendingNonceAt(context.Background(), a.address)
		if err != nil {
			return 0, fmt.Errorf("failed to get nonce of account %x; %v", a.address, err)
		}
	}
	a.nonce = current + 1
	return current, nil
}

// WaitUntilAllTxsApplied blocks until all txs with nonces from getNextNonce are in the chain
func (a *Account) WaitUntilAllTxsApplied(rpcClient *ethclient.Client) error {
	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		nonce, err := rpcClient.NonceAt(context.Background(), a.address, nil) // nonce at latest block
		if err != nil {
			return fmt.Errorf("failed to check address nonce; %v", err)
		}
		if nonce == a.nonce {
			return nil
		}
	}
	return fmt.Errorf("nonce not increased before timeout")
}

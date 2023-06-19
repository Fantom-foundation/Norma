package app

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

// transferValue transfer a financial value from account identified by given privateKey, to given toAddress.
// It returns when the value is already available on the target account.
func transferValue(rpcClient RpcClient, from *Account, toAddress common.Address, value *big.Int) (err error) {
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	signedTx, err := createTx(from, toAddress, value, nil, gasPrice, 21000)
	if err != nil {
		return err
	}
	return rpcClient.SendTransaction(context.Background(), signedTx)
}

func createTx(from *Account, toAddress common.Address, value *big.Int, data []byte, gasPrice *big.Int, gasLimit uint64) (*types.Transaction, error) {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    from.getNextNonce(),
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})
	return types.SignTx(tx, types.NewEIP155Signer(from.chainID), from.privateKey)
}

// waitUntilAccountNonceIs blocks until the account nonce at the latest block on the chain is given value
func waitUntilAccountNonceIs(account common.Address, awaitedNonce uint64, rpcClient RpcClient) error {
	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		nonce, err := rpcClient.NonceAt(context.Background(), account, nil) // nonce at latest block
		if err != nil {
			return fmt.Errorf("failed to check address nonce; %v", err)
		}
		if nonce == awaitedNonce {
			return nil
		}
	}
	return fmt.Errorf("nonce not achieved before timeout")
}

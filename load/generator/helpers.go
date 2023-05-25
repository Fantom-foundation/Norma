package generator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

// generateAddress generate a new pair of private key and the account address
func generateAddress() (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return address, privateKey, nil
}

// transferValue transfer a financial value from account identified by given privateKey, to given toAddress.
// It returns when the value is already available on the target account.
func transferValue(rpcClient *ethclient.Client, chainID *big.Int, privateKey *ecdsa.PrivateKey, toAddress common.Address, value *big.Int) error {
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := rpcClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      21000, // standard amount of gas for plain transfer
		To:       &toAddress,
		Value:    value,
	})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	err = rpcClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	return waitUntilAccountHasBalance(toAddress, rpcClient)
}

// waitUntilAccountHasBalance allows to wait until the given address has non-zero balance.
func waitUntilAccountHasBalance(address common.Address, rpcClient *ethclient.Client) error {
	for i := 0; i < 150; i++ {
		time.Sleep(time.Second)
		balance, err := rpcClient.BalanceAt(context.Background(), address, nil)
		if err != nil {
			return fmt.Errorf("failed to check contract existence; %v", err)
		}
		if balance.Uint64() != 0 {
			return nil
		}
	}
	return fmt.Errorf("transfered balance available before timeout")
}

// waitUntilContractStartExisting allows to wait until the given contract is available on the chain.
func waitUntilContractStartExisting(contractAddress common.Address, rpcClient *ethclient.Client) error {
	for i := 0; i < 150; i++ {
		time.Sleep(time.Second)
		code, err := rpcClient.CodeAt(context.Background(), contractAddress, nil)
		if err != nil {
			return fmt.Errorf("failed to check contract existence; %v", err)
		}
		if len(code) != 0 {
			return nil
		}
	}
	return fmt.Errorf("deployed contract not available before timeout")
}

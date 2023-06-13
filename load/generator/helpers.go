package generator

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
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
func transferValue(rpcClient *ethclient.Client, from *Account, toAddress common.Address, value *big.Int) (err error) {
	gasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	nonce, err := from.getNextNonce(rpcClient)
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
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(from.chainID), from.privateKey)
	if err != nil {
		return err
	}
	return rpcClient.SendTransaction(context.Background(), signedTx)
}

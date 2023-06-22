package rules

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/Norma/common/transact"
	"github.com/Fantom-foundation/Norma/driver/network/rules/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var nodeDriverAuthAddress = common.HexToAddress("0xd100ae0000000000000000000000000000000000")

func SetNetworkRules(rpcClient transact.RpcClient, ownerAccount *transact.Account, rulesDiff string) error {
	// get price of gas from the network
	regularGasPrice, err := rpcClient.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to suggest gas price; %v", err)
	}

	authContract, err := abi.NewNodeDriverAuth(nodeDriverAuthAddress, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get NodeDriverAuth contract representation; %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(ownerAccount.PrivateKey, ownerAccount.ChainID)
	if err != nil {
		return fmt.Errorf("failed to create txOpts; %v", err)
	}
	txOpts.GasPrice = transact.GetPriorityGasPrice(regularGasPrice)
	txOpts.Nonce = big.NewInt(int64(ownerAccount.GetNextNonce()))

	_, err = authContract.UpdateNetworkRules(txOpts, []byte(rulesDiff))
	if err != nil {
		return fmt.Errorf("failed to update network rules; %v", err)
	}
	return nil
}

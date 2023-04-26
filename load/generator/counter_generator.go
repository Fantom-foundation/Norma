package generator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

//go:generate solc -o ../contracts/abi --overwrite --pretty-json --optimize --abi --bin ../contracts/Counter.sol
//go:generate abigen --type Counter --pkg abi --abi ../contracts/abi/Counter.abi --bin ../contracts/abi/Counter.bin --out ../contracts/abi/Counter.go

// CounterTransactionGenerator is a txs generator incrementing trivial Counter contract
type CounterTransactionGenerator struct {
	auth            *bind.TransactOpts
	contractAddress common.Address
	counterContract *abi.Counter
}

func NewCounterTransactionGenerator(privateKey *ecdsa.PrivateKey, chainID *big.Int) (*CounterTransactionGenerator, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	return &CounterTransactionGenerator{
		auth: auth,
	}, nil
}

func (cg *CounterTransactionGenerator) Init(rpcClient *ethclient.Client) (err error) {

	// deploy testing contract
	cg.contractAddress, _, cg.counterContract, err = abi.DeployCounter(cg.auth, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to deploy Counter contract; %v", err)
	}
	// counterContract can be obtained after the deployment as NewCounter(cg.contractAddress, rpcClient)

	return waitUntilContractStartExisting(cg.contractAddress, rpcClient)
}

func waitUntilContractStartExisting(contractAddress common.Address, rpcClient *ethclient.Client) error {
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
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

func (cg *CounterTransactionGenerator) SendTx() error {
	_, err := cg.counterContract.IncrementCounter(cg.auth)
	return err
}

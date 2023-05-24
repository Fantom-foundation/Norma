package generator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/Fantom-foundation/Norma/load/contracts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

func (cg *CounterTransactionGenerator) SendTx() error {
	_, err := cg.counterContract.IncrementCounter(cg.auth)
	return err
}

func (cg *CounterTransactionGenerator) GetAmountOfReceivedTxs() (uint64, error) {
	count, err := cg.counterContract.GetCount(nil)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

package app

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockgen -source app.go -destination app_mock.go -package app

type Application interface {
	// CreateGenerator creates a new transaction generator for given application
	CreateGenerator(rpcClient RpcClient) (TransactionGenerator, error)

	WaitUntilApplicationIsDeployed(rpcClient RpcClient) error

	GetReceivedTransations(rpcClient RpcClient) (uint64, error)
}

type RpcClient interface {
	bind.ContractBackend
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	Close()
}

// TransactionGenerator produces a stream of transactions to generate traffic on the chain.
// Generators are not thread-safe.
type TransactionGenerator interface {
	GenerateTx() (*types.Transaction, error)
	GetSentTransactions() uint64
}

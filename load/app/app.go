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
	// CreateUser creates a new user generating transactions for this application.
	CreateUser(rpcClient RpcClient) (User, error)

	WaitUntilApplicationIsDeployed(rpcClient RpcClient) error

	GetReceivedTransations(rpcClient RpcClient) (uint64, error)
}

type RpcClient interface {
	bind.ContractBackend
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	Close()
}

// User produces a stream of transactions to generate traffic on the chain.
// Implementations are not required to be thread-safe.
type User interface {
	GenerateTx() (*types.Transaction, error)
	GetSentTransactions() uint64
}

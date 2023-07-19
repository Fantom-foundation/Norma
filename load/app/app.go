package app

import (
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockgen -source app.go -destination app_mock.go -package app

type Application interface {
	// CreateUser creates a new user generating transactions for this application.
	CreateUser(rpcClient rpc.RpcClient) (User, error)

	WaitUntilApplicationIsDeployed(rpcClient rpc.RpcClient) error

	GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error)
}

// User produces a stream of transactions to generate traffic on the chain.
// Implementations are not required to be thread-safe.
type User interface {
	GenerateTx() (*types.Transaction, error)
	GetSentTransactions() uint64
}

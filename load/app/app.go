package app

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

//go:generate mockgen -source app.go -destination app_mock.go -package app

type Application interface {
	// CreateGenerator creates a new transaction generator for given application
	CreateGenerator(rpcClient RpcClient) (TransactionGenerator, error)

	WaitUntilApplicationIsDeployed(rpcClient RpcClient) error
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
}

type ApplicationProvidingTxCount interface {
	Application
	GetTransactionCounts(rpcClient RpcClient) (TransactionCounts, error)
}

// TransactionCounts should be implemented by an instance that can provide the number of received
// and expected transactions.
type TransactionCounts struct {
	// SentTxs represents the number of transactions originally sent to an application.
	// This number of transactions was not necessarily received by the application as the transactions
	// could be filtered out by any layers between the RPC endpoint and actual block processing,
	// or the client was not able to process requested amount of transactions and the transactions could not reach
	// the block processing.
	SentTxs uint64

	// ReceivedTxs represents the number of transactions received by an application.
	// This number of transactions may be smaller than the number of actually sent transactions
	// as the transactions could be filtered out by any layers between the RPC endpoint and actual block processing,
	// or the client was not able to process requested amount of transactions and the transactions could not reach
	// the block processing.
	ReceivedTxs uint64
}

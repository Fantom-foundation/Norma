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
	// SentTxs represents the number of transactions originally sent to an application per account. The
	// slice is indexed by the account IDs such that SentTxs[14] is the number of transactions send by
	// account 14 of the appication that has returned this struct. The length may by differ from the total
	// number of accounts, in which case all non-covered account lengths should be considered to have sent
	// zero transactions.
	// The number of transactions reported here are not necessarily received by the application as the transactions
	// could be filtered out by any layers between the RPC endpoint and actual block processing,
	// or the client was not able to process requested amount of transactions and the transactions could not reach
	// the block processing.
	SentTxs []uint64

	// ReceivedTxs represents the number of transactions received by an application.
	// This number of transactions may be smaller than the number of actually sent transactions
	// as the transactions could be filtered out by any layers between the RPC endpoint and actual block processing,
	// or the client was not able to process requested amount of transactions and the transactions could not reach
	// the block processing.
	ReceivedTxs uint64
}

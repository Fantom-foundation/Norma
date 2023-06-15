package app

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:generate mockgen -source app.go -destination app_mock.go -package app

type Application interface {
	// CreateGenerator creates a new transaction generator for given application
	CreateGenerator(rpcClient *ethclient.Client) (TransactionGenerator, error)

	WaitUntilGeneratorsCreated(rpcClient *ethclient.Client) error
}

// TransactionGenerator produces a stream of transactions to generate traffic on the chain.
// Generators are not thread-safe.
type TransactionGenerator interface {
	GenerateTx() (*types.Transaction, error)
}

type ApplicationProvidingTxCount interface {
	Application
	TransactionCountsProvider
}

// TransactionCountsProvider should be implemented by an Application that can provide the number of received txs.
type TransactionCountsProvider interface {
	// GetAmountOfSentTxs returns the number of transactions originally sent to an application.
	// This number of transactions was not necessarily received by the application as the transactions
	// could be filtered out by any layers between the RPC endpoint and actual block processing,
	// or the client was not able to process requested amount of transactions and the transactions could not reach
	// the block processing.
	GetAmountOfSentTxs() uint64

	// GetAmountOfReceivedTxs returns the number of transactions received by an application.
	// This number of transactions may be smaller than the number of actually sent transactions
	// as the transactions could be filtered out by any layers between the RPC endpoint and actual block processing,
	// or the client was not able to process requested amount of transactions and the transactions could not reach
	// the block processing.
	GetAmountOfReceivedTxs(rpcClient *ethclient.Client) (uint64, error)
}

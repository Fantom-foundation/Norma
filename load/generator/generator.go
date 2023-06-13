package generator

import (
	"io"
)

//go:generate mockgen -source generator.go -destination generator_mock.go -package generator

// URL is a mere alias type for a string supposed to encode a URL.
type URL string

// TransactionGenerator produces a stream of transactions to generate traffic on the chain.
// Generators are not thread-safe.
type TransactionGenerator interface {
	// SendTx generates a new tx and send it to the RPC
	SendTx() error

	io.Closer
}

type TransactionGeneratorFactory interface {
	// Create and return a new generator instance
	Create() (TransactionGenerator, error)

	// WaitForInit blocks until generators initialization is finished in the latest block of the chain
	// Should be called after a batch of Create calls, before the generators will start to be used.
	WaitForInit() error
}

type TransactionGeneratorFactoryWithStats interface {
	TransactionGeneratorFactory
	TransactionCounts
}

// TransactionCounts should be implemented by an instance that can provide the number of received
// and expected transactions.
type TransactionCounts interface {
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
	GetAmountOfReceivedTxs() (uint64, error)
}

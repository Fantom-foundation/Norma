package generator

import (
	"io"
)

//go:generate mockgen -source generator.go -destination generator_mock.go -package generator

// TransactionGenerator produces a stream of transactions to generate a traffic on the chain.
type TransactionGenerator interface {
	// SendTx generates a new tx and send it to the RPC
	SendTx() error

	io.Closer
}

type TransactionGeneratorFactory interface {
	Create() (TransactionGenerator, error)
}

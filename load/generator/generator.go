package generator

import "github.com/ethereum/go-ethereum/ethclient"

//go:generate mockgen -source generator.go -destination generator_mock.go -package generator

// TransactionGenerator produces a stream of transactions to generate a traffic on the chain.
type TransactionGenerator interface {
	// Init prepares the txs generating - deploys contracts, transfer tokens to senders addresses etc.
	// Provided client will be used also for following txs
	Init(client *ethclient.Client) error

	// SendTx generates a new tx and send it to the RPC
	SendTx() error
}

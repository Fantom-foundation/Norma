package monitoring

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

// Throughput is a monitoring interface to analyze speed related properties of nodes
type Throughput interface {

	// GetTransactions returns the number of transactions processed in the input block
	// ErrNotFound is returned when the block does not exist
	GetTransactions(block int) (int, error)

	// GetGas returns the amount of Gas consumed in the input block
	// ErrNotFound is returned when the block does not exist
	GetGas(block int) (int, error)

	// GetBlockTime returns timestamp of the input block
	// ErrNotFound is returned when the block does not exist
	GetBlockTime(block int) (time.Time, error)

	// GetBlockDelay returns time that passed between start of this block and end of previous block,
	// i.e. delay between blocks
	// ErrNotFound is returned when the block does not exist
	GetBlockDelay(block int) (time.Duration, error)

	// GetBlockHeight returns number of already processed blocks
	GetBlockHeight() (int, error)
}

// System is a monitoring interface to analyze system utilization
type System interface {

	// GetCPULoad returns current cpu load range [0,1]
	GetCPULoad() float32

	// GetMemoryLoad returns memory utilization in MB
	GetMemoryLoad() int

	// GetMemoryCapacity returns available memory
	GetMemoryCapacity() int

	// GetStorageLoad returns storage space utilization in MB
	GetStorageLoad() int
}

// Opera is a monitoring interface to analyze Opera client specific properties
type Opera interface {

	// GetBlockProcessingTime returns time spent specifically for block processing,
	// i.e. excluding the concensus, network, transaction ordering, validation, etc.
	GetBlockProcessingTime(block int) time.Duration

	// GetBlockCommitmentTime returns time spent in committing a block after it has been processed.
	GetBlockCommitmentTime(block int) time.Duration

	// GetTransactionPoolSize return current number of transactions in the transaction pool
	GetTransactionPoolSize() int
}

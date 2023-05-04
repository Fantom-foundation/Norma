package monitoring

import "time"

// Throughput is a monitoring interface to analyze speed related properties of nodes
type Throughput interface {

	// GetTransactions returns the number of transactions processed in the input block
	GetTransactions(block int) int

	// GetGas returns the amount of Gas consumed in the input block
	GetGas(block int) int

	// GetBlockTime returns total time to process the input block
	GetBlockTime(block int) time.Duration

	// GetBlockDelay returns time that passed between start of this block and end of previous block,
	// i.e. delay between blocks
	GetBlockDelay(block int) time.Duration
}

// System is a monitoring interface to analyze system utilization
type System interface {

	// GetCPULoad returns current cpu load in percentage
	GetCPULoad() int

	// GetMemoryLoad returns memory utilization in MB
	GetMemoryLoad() int

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

	// GetTransactionDropRate returns percentage of dropped transactions from the pool
	GetTransactionDropRate() int
}

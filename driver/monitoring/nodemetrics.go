package monitoring

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNotFound = errors.New("not found")

// NodeMetrics allows for getting metrics of a node and block such as block height, transactions in a block, gas used in a block, etc.
type NodeMetrics struct {
	blockStream <-chan Block // allows for collecting stream throughput related data of a block

	data      *sync.Map // map block number -> block, for now it stores all data in-memory
	lastBlock *atomic.Int64
	done      chan bool
}

// CreateNodeMetrics creates a new instance.
func CreateNodeMetrics(blockStream <-chan Block) *NodeMetrics {
	m := &NodeMetrics{
		blockStream: blockStream,
		data:        new(sync.Map),
		lastBlock:   new(atomic.Int64),
		done:        make(chan bool),
	}

	m.start()
	return m
}

func (n *NodeMetrics) start() {
	go func() {
		defer close(n.done)
		for b := range n.blockStream {
			n.lastBlock.Store(int64(b.height))
			n.data.Store(b.height, b)
		}
	}()
}

func (n *NodeMetrics) drain() {
	<-n.done
}

func (n *NodeMetrics) GetNumberOfTransactions(block int) (int, error) {
	b, err := n.getBlock(block)
	if err != nil {
		return 0, err
	}

	return b.txs, nil
}

func (n *NodeMetrics) GetGas(block int) (int, error) {
	b, err := n.getBlock(block)
	if err != nil {
		return 0, err
	}

	return b.gasUsed, nil
}

func (n *NodeMetrics) GetBlockTime(block int) (time.Time, error) {
	b, err := n.getBlock(block)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return b.time, nil
}

func (n *NodeMetrics) GetBlockDelay(block int) (time.Duration, error) {
	b, err := n.getBlock(block)
	if err != nil {
		return 0, err
	}

	prev, err := n.getBlock(block - 1)
	if err != nil {
		return 0, err
	}

	return b.time.Sub(prev.time), nil
}

func (n *NodeMetrics) GetBlockHeight() (int, error) {
	return int(n.lastBlock.Load()), nil
}

// getBlock reads block information from internal cache. If the block is not there, it returns ErrNotFound
func (n *NodeMetrics) getBlock(blockNum int) (Block, error) {
	if b, exists := n.data.Load(blockNum); exists {
		return b.(Block), nil
	} else {
		var empty Block
		return empty, ErrNotFound
	}
}

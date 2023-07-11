package rpc

import (
	"github.com/ethereum/go-ethereum/core/types"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryRpcReturnGracefully(t *testing.T) {
	t.Parallel()

	pool := NewRpcWorkerPool()
	start := time.Now()
	err := pool.runRpcSenderLoop("wrong", 6)

	// we expect error, but should not panic
	if err == nil {
		t.Errorf("method was expected to fail")
	}

	if got, want := time.Now().Sub(start).Seconds(), float64(5); got < want {
		t.Errorf("RPC should be attempted around 6s, was: %f < %f", got, want)
	}
}

func TestClosePool(t *testing.T) {
	t.Parallel()

	pool := NewRpcWorkerPool()
	counter := &atomic.Int32{}
	wg := &sync.WaitGroup{}

	go func() {
		for range pool.txs {
			counter.Add(1)
			wg.Done()
		}
		wg.Done() // will get here when the channel is closed
	}()

	var tx types.Transaction
	for i := 0; i < 10; i++ {
		wg.Add(1)
		pool.SendTransaction(&tx)
	}

	wg.Wait()
	wg.Add(1) // extra count to check the go routine ended

	if got, want := counter.Load(), int32(10); got != want {
		t.Errorf("not all data read from the channel: %d != %d", got, want)
	}

	if err := pool.Close(); err != nil {
		t.Fatalf("error: %s", err)
	}

	wg.Wait()

	if got, want := counter.Load(), int32(10); got > want {
		t.Errorf("not all data read from the channel: %d != %d", got, want)
	}
}

// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

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

	start := time.Now()
	txs := make(chan *types.Transaction)
	w := newWorker("wrong", txs)

	time.Sleep(6 * time.Second)
	w.close()

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

func TestCloseWorkerStartStop(t *testing.T) {
	txs := make(chan *types.Transaction)
	w := newWorker("wrong", txs)
	w.close()
}

func TestCloseWorkerGroupStartStop(t *testing.T) {
	txs := make(chan *types.Transaction)
	wg := workerGroup{}
	for i := 0; i < 150; i++ {
		wg.add("wrong", txs)
	}
	wg.close()
}

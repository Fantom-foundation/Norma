package rpc

import (
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
	"time"
)

func TestRetryRpcReturnGracefully(t *testing.T) {
	t.Parallel()

	txs := make(chan *types.Transaction)

	start := time.Now()
	err := runRpcSenderLoop("wrong", 6, txs)

	// we expect error, but should not panic
	if err == nil {
		t.Errorf("method was expected to fail")
	}

	if got, want := time.Now().Sub(start).Seconds(), float64(5); got < want {
		t.Errorf("RPC should be attempted around 6s, was: %f < %f", got, want)
	}
}

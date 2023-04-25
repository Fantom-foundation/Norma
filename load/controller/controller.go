package controller

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

// AppController emits transactions to the testing network into a blockchain app to generate a load.
// The Shaper passed into the driver controls the frequency of emitting transactions.
// The Generator passed into the driver constructs the transactions.
// The RPC Client is used to send the transactions into the network.
type AppController struct {
	generator generator.TransactionGenerator
	shaper    shaper.Shaper
	rpcClient *ethclient.Client
}

func NewAppController(generator generator.TransactionGenerator, shaper shaper.Shaper, rpcClient *ethclient.Client) *AppController {
	return &AppController{
		generator: generator,
		shaper:    shaper,
		rpcClient: rpcClient,
	}
}

func (sd *AppController) Init() error {
	// initialize generator
	return sd.generator.Init(sd.rpcClient)
}

func (sd *AppController) Run(ctx context.Context) error {

	for {
		startTime := time.Now()
		err := sd.generator.SendTx()
		if err != nil {
			return fmt.Errorf("failed to send tx; %v", err)
		}

		waitTime := sd.shaper.GetNextWaitTime()
		waitTime -= time.Since(startTime) // subtract time consumed by generating
		if waitTime > 0 {
			time.Sleep(waitTime)
		}

		// interrupt the loop if the context have been cancelled
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err == context.DeadlineExceeded || err == context.Canceled {
				return nil // terminated gracefully
			}
			return err
		default:
			// no interruption - continue
		}
	}
}

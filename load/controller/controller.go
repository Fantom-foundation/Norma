package controller

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"log"
	"time"
)

// AppController emits transactions to the testing network into a blockchain app to generate a load.
// The Shaper passed into the driver controls the frequency of emitting transactions.
// The Generator passed into the driver constructs the transactions.
// The RPC Client is used to send the transactions into the network.
type AppController struct {
	generatorFactory generator.TransactionGeneratorFactory
	shaper           shaper.Shaper
	workers          int
	trigger          chan struct{}
}

func NewAppController(generatorFactory generator.TransactionGeneratorFactory, shaper shaper.Shaper, workers int) *AppController {
	return &AppController{
		generatorFactory: generatorFactory,
		shaper:           shaper,
		workers:          workers,
		trigger:          make(chan struct{}),
	}
}

func (ac *AppController) Init() error {
	// initialize workers
	for i := 0; i < ac.workers; i++ {
		gen, err := ac.generatorFactory.Create()
		if err != nil {
			return fmt.Errorf("failed to create load generator; %s", err)
		}

		worker := Worker{
			generator: gen,
			trigger:   ac.trigger,
		}
		go worker.Run()
	}

	return nil
}

func (ac *AppController) Run(ctx context.Context) error {
	defer close(ac.trigger)
	for {
		select {
		case <-ctx.Done():
			// interrupt the loop if the context has been cancelled
			err := ctx.Err()
			if err == context.DeadlineExceeded || err == context.Canceled {
				return nil // terminated gracefully
			}
			return err
		default:
			// trigger a worker to send a tx
			select {
			case ac.trigger <- struct{}{}:
			default:
				log.Print("sending not fast enough for the required frequency")
			}

			// wait for time determined by the shaper
			time.Sleep(ac.shaper.GetNextWaitTime())
		}
	}
}

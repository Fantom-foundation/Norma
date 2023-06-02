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
	shaper  shaper.Shaper
	trigger chan struct{}
}

func NewAppController(generatorFactory generator.TransactionGeneratorFactory, shaper shaper.Shaper, accounts int) (*AppController, error) {
	trigger := make(chan struct{})

	// initialize workers for individual accounts
	for i := 0; i < accounts; i++ {
		gen, err := generatorFactory.Create()
		if err != nil {
			return nil, fmt.Errorf("failed to create load generator; %s", err)
		}

		worker := Worker{
			generator: gen,
			trigger:   trigger,
		}
		go worker.Run()
	}

	return &AppController{
		shaper:  shaper,
		trigger: trigger,
	}, nil
}

func (ac *AppController) Run(ctx context.Context) error {
	defer close(ac.trigger)
	missed := 0
	for {
		select {
		case <-ctx.Done():
			// interrupt the loop if the context has been cancelled
			if missed != 0 {
				log.Printf("sending not fast enough for the required frequency: %d times\n", missed)
			}
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
				missed++
				log.Print("sending not fast enough for the required frequency")
			}

			// wait for time determined by the shaper
			time.Sleep(ac.shaper.GetNextWaitTime())
		}
	}
}

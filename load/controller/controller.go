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

func NewAppController(generatorFactory generator.TransactionGeneratorFactory, shaper shaper.Shaper, workers int) (*AppController, error) {
	trigger := make(chan struct{})

	// initialize workers
	for i := 0; i < workers; i++ {
		generator, err := generatorFactory.Create()
		if err != nil {
			return nil, fmt.Errorf("failed to create load generator; %s", err)
		}

		worker := Worker{
			generator: generator,
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

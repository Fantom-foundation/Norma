package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/shaper"
)

// AppController emits transactions to the testing network into a blockchain app to generate a load.
// The Shaper passed into the driver controls the frequency of emitting transactions.
// The Generator passed into the driver constructs the transactions.
// The RPC Client is used to send the transactions into the network.
type AppController struct {
	shaper         shaper.Shaper
	countsProvider app.ApplicationProvidingTxCount
	network        driver.Network
	trigger        chan struct{}
}

func NewAppController(application app.Application, shaper shaper.Shaper, generators int, network driver.Network) (*AppController, error) {
	trigger := make(chan struct{})

	rpcClient, err := network.DialRandomRpc()
	if err != nil {
		return nil, fmt.Errorf("failed to dial ranom RPC; %v", err)
	}
	defer rpcClient.Close()

	// initialize workers for individual generators
	for i := 0; i < generators; i++ {
		gen, err := application.CreateGenerator(rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create load app; %s", err)
		}

		go runGeneratorLoop(gen, trigger, network)
	}

	// wait until all changes are on the chain
	if err := application.WaitUntilApplicationIsDeployed(rpcClient); err != nil {
		return nil, fmt.Errorf("failed to wait for app on-chain init; %s", err)
	}

	countsProvider, _ := application.(app.ApplicationProvidingTxCount) // nil if not TxCount provider
	return &AppController{
		shaper:         shaper,
		countsProvider: countsProvider,
		network:        network,
		trigger:        trigger,
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
			waitTime, shouldSend := ac.shaper.GetNextWaitTime()

			// send only if the shaper says so
			if shouldSend {
				// trigger a worker to send a tx
				select {
				case ac.trigger <- struct{}{}:
				default:
					missed++
				}
			}

			// wait for time determined by the shaper
			time.Sleep(waitTime)
		}
	}
}

// GetTransactionCounts returns the object that provides the number of send and received transactions
// of application managed by this application controller.
// If this application controller is not capable of providing such an information, this method returns
// ErrDoesNotProvideTxCounts in its error return argument.
func (ac *AppController) GetTransactionCounts() (app.TransactionCounts, error) {
	if ac.countsProvider == nil {
		return app.TransactionCounts{}, ErrDoesNotProvideTxCounts
	}

	rpcClient, err := ac.network.DialRandomRpc()
	if err != nil {
		return app.TransactionCounts{}, fmt.Errorf("failed to dial ranom RPC; %v", err)
	}
	defer rpcClient.Close()

	return ac.countsProvider.GetTransactionCounts(rpcClient)
}

var ErrDoesNotProvideTxCounts = errors.New("app does not provide tx counts")

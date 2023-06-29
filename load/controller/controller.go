package controller

import (
	"context"
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
	shaper      shaper.Shaper
	application app.Application
	network     driver.Network
	trigger     chan struct{}
	accounts    []app.TransactionGenerator
	rpcClient   app.RpcClient
}

func NewAppController(application app.Application, shaper shaper.Shaper, generators int, network driver.Network) (*AppController, error) {
	trigger := make(chan struct{})

	rpcClient, err := network.DialRandomRpc()
	if err != nil {
		return nil, fmt.Errorf("failed to dial ranom RPC; %v", err)
	}

	// initialize workers for individual generators
	accounts := make([]app.TransactionGenerator, 0, generators)
	for i := 0; i < generators; i++ {
		gen, err := application.CreateGenerator(rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create load app; %s", err)
		}

		go runGeneratorLoop(gen, trigger, network)
		accounts = append(accounts, gen)
	}

	// wait until all changes are on the chain
	if err := application.WaitUntilApplicationIsDeployed(rpcClient); err != nil {
		return nil, fmt.Errorf("failed to wait for app on-chain init; %s", err)
	}

	return &AppController{
		shaper:      shaper,
		application: application,
		network:     network,
		trigger:     trigger,
		accounts:    accounts,
		rpcClient:   rpcClient,
	}, nil
}

func (ac *AppController) Run(ctx context.Context) error {
	defer close(ac.trigger)
	defer ac.rpcClient.Close()
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

func (ac *AppController) GetNumberOfAccounts() int {
	return len(ac.accounts)
}

func (ac *AppController) GetSentTransactions(account int) (uint64, error) {
	if account < 0 || account >= len(ac.accounts) {
		return 0, nil
	}
	return ac.accounts[account].GetSentTransactions(), nil
}

func (ac *AppController) GetReceivedTransactions() (uint64, error) {
	for retry := 0; ; retry++ {
		// fetch transaction data from the network
		res, err := ac.application.GetReceivedTransations(ac.rpcClient)
		if err == nil {
			return res, nil
		}
		if retry >= 5 {
			return 0, err
		}

		// attempt a re-connect
		ac.rpcClient.Close()
		ac.rpcClient, err = ac.network.DialRandomRpc()
		if err != nil {
			return 0, fmt.Errorf("failed to dial random RPC; %v", err)
		}
	}
}

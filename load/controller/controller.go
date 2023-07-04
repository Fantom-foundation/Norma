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
	users       []app.User
	rpcClient   app.RpcClient
}

func NewAppController(application app.Application, shaper shaper.Shaper, numUsers int, network driver.Network) (*AppController, error) {
	trigger := make(chan struct{})

	rpcClient, err := network.DialRandomRpc()
	if err != nil {
		return nil, fmt.Errorf("failed to dial ranom RPC; %v", err)
	}

	// initialize workers for individual generators
	users := make([]app.User, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		gen, err := application.CreateUser(rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create load app; %s", err)
		}

		go runGeneratorLoop(gen, trigger, network)
		users = append(users, gen)
		if i%100 == 0 {
			log.Printf("initialized %d of %d users ...\n", i+1, numUsers)
		}
	}

	// wait until all changes are on the chain
	log.Printf("waiting until all app generators are deployed...\n")
	if err := application.WaitUntilApplicationIsDeployed(rpcClient); err != nil {
		return nil, fmt.Errorf("failed to wait for app on-chain init; %s", err)
	}
	log.Printf("the app is deployed\n")

	return &AppController{
		shaper:      shaper,
		application: application,
		network:     network,
		trigger:     trigger,
		users:       users,
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

func (ac *AppController) GetNumberOfUsers() int {
	return len(ac.users)
}

func (ac *AppController) GetSentTransactions(user int) (uint64, error) {
	if user < 0 || user >= len(ac.users) {
		return 0, nil
	}
	return ac.users[user].GetSentTransactions(), nil
}

func (ac *AppController) GetReceivedTransactions() (uint64, error) {
	for retry := 0; ; retry++ {
		// fetch transaction data from the network
		res, err := ac.application.GetReceivedTransactions(ac.rpcClient)
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

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

package controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
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
	rpcClient   rpc.RpcClient
}

func NewAppController(application app.Application, shaper shaper.Shaper, numUsers int, rpcClient rpc.RpcClient, network driver.Network) (*AppController, error) {
	trigger := make(chan struct{}, 100)

	if rpcClient == nil {
		var err error
		rpcClient, err = network.DialRandomRpc()
		if err != nil {
			return nil, fmt.Errorf("failed to dial ranom RPC; %v", err)
		}
	}

	// initialize workers for individual generators
	users := make([]app.User, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		gen, err := application.CreateUser(rpcClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create load app; %s", err)
		}
		users = append(users, gen)
		//if i%100 == 0 {
		log.Printf("initialized %d of %d users ...\n", i+1, numUsers)
		//}
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
	defer ac.rpcClient.Close()

	// start generators for each user
	var done sync.WaitGroup
	for _, user := range ac.users {
		user := user
		done.Add(1)
		go func() {
			defer done.Done()
			runGeneratorLoop(user, ac.trigger, ac.network)
		}()
	}

	var pending float64
	lastUpdate := time.Now()
	ac.shaper.Start(lastUpdate, ac)

	for {
		// re-plenish the number of pending messages
		now := time.Now()
		pending += ac.shaper.GetNumMessagesInInterval(lastUpdate, now.Sub(lastUpdate))
		lastUpdate = now

		for pending > 0 {
			ac.trigger <- struct{}{}
			pending -= 1
		}

		select {
		case <-time.After(time.Millisecond):
			// just waiting for next time to send messages.
		case <-ctx.Done():
			close(ac.trigger)
			done.Wait()
			err := ctx.Err()
			if err == context.DeadlineExceeded || err == context.Canceled {
				return nil // terminated gracefully
			}
			return err
		}
	}
}

func (ac *AppController) GetNumberOfUsers() int {
	return len(ac.users)
}

func (ac *AppController) GetTransactionsSentBy(user int) (uint64, error) {
	if user < 0 || user >= len(ac.users) {
		return 0, nil
	}
	return ac.users[user].GetSentTransactions(), nil
}

func (ac *AppController) GetSentTransactions() (uint64, error) {
	sum := uint64(0)
	for i := 0; i < ac.GetNumberOfUsers(); i++ {
		cur, _ := ac.GetTransactionsSentBy(i)
		sum += cur
	}
	return sum, nil
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

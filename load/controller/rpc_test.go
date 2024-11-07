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

package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/shaper"
)

const PrivateKey = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7" // Fakenet validator 1
const FakeNetworkID = 0xfa3

func TestTrafficGenerating(t *testing.T) {
	// run local network of one node
	net, err := local.NewLocalNetwork(&driver.NetworkConfig{NumberOfValidators: 1})
	if err != nil {
		t.Fatalf("failed to create new local network: %v", err)
	}
	t.Cleanup(func() { net.Shutdown() })

	primaryAccount, err := app.NewAccount(0, PrivateKey, nil, FakeNetworkID)
	if err != nil {
		t.Fatal(err)
	}

	appContext, err := app.NewContext(net, primaryAccount)
	if err != nil {
		t.Fatal(err)
	}

	application, err := app.NewERC20Application(appContext, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	constantShaper := shaper.NewConstantShaper(30.0) // 30 txs/sec

	numGenerators := 5 // 5 parallel workers
	app, err := controller.NewAppController(application, constantShaper, numGenerators, appContext, net)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := app.GetNumberOfUsers(), numGenerators; got != want {
		t.Errorf("unexpected number of accounts, wanted %d, got %d", want, got)
	}

	// let the app run for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// run the app in the same thread, will be interrupted by the context timeout
	err = app.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second) // wait for txs in TxPool

	// get amount of txs applied to the chain
	sum, err := app.GetSentTransactions()
	if err != nil {
		t.Fatalf("failed to fetch sent transactions: %v", err)
	}

	if received, err := app.GetReceivedTransactions(); err != nil || received != sum {
		t.Errorf("invalid number of received transactions on the network, wanted %d, got %d, err %v", sum, received, err)
	}

	// in optimal case should be generated 30 txs per second
	// as a tolerance for slow CI we require at least 20 txs
	// and at most 40, since timing is not perfect.
	if sum < 20 || sum > 40 {
		t.Errorf("unexpected amount of generated txs: %d", sum)
	}
}

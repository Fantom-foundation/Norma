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

package app_test

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/network"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/load/app"
)

const PrivateKey = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7" // Fakenet validator 1
const FakeNetworkID = 0xfa3

func TestGenerators(t *testing.T) {
	// run local network of one node
	net, err := local.NewLocalNetwork(&driver.NetworkConfig{NumberOfValidators: 1})
	if err != nil {
		t.Fatalf("failed to create new local network: %v", err)
	}
	t.Cleanup(func() { net.Shutdown() })

	rpcClient, err := net.DialRandomRpc()
	if err != nil {
		t.Fatal("unable to connect the the rpc")
	}

	primaryAccount, err := app.NewAccount(0, PrivateKey, FakeNetworkID)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Counter", func(t *testing.T) {
		counterApp, err := app.NewCounterApplication(rpcClient, primaryAccount, 1, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, counterApp, rpcClient)
	})
	t.Run("ERC20", func(t *testing.T) {
		erc20app, err := app.NewERC20Application(rpcClient, primaryAccount, 1, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, erc20app, rpcClient)
	})
	t.Run("Store", func(t *testing.T) {
		storeApp, err := app.NewStoreApplication(rpcClient, primaryAccount, 1, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, storeApp, rpcClient)
	})
	t.Run("Uniswap", func(t *testing.T) {
		uniswapApp, err := app.NewUniswapApplication(rpcClient, primaryAccount, 1, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, uniswapApp, rpcClient)
	})
}

func testGenerator(t *testing.T, app app.Application, rpcClient rpc.RpcClient) {
	gen, err := app.CreateUser(rpcClient)
	if err != nil {
		t.Fatal(err)
	}
	err = app.WaitUntilApplicationIsDeployed(rpcClient)
	if err != nil {
		t.Fatal(err)
	}

	numTransactions := 10
	for i := 0; i < numTransactions; i++ {
		tx, err := gen.GenerateTx()
		if err != nil {
			t.Fatal(err)
		}
		if err := rpcClient.SendTransaction(context.Background(), tx); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second) // wait for txs in TxPool

	if got, want := gen.GetSentTransactions(), numTransactions; got != uint64(want) {
		t.Errorf("invalid number of sent transactions reported, wanted %d, got %d", want, got)
	}

	err = network.Retry(network.DefaultRetryAttempts, 1*time.Second, func() error {
		received, err := app.GetReceivedTransactions(rpcClient)
		if err != nil {
			return fmt.Errorf("unable to get amount of received txs; %v", err)
		}
		if received != 10 {
			return fmt.Errorf("unexpected amount of txs in chain (%d)", received)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

package controller_test

import (
	"context"
	"github.com/Fantom-foundation/Norma/driver/network"
	"sync"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/ethclient"
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

	rpcUrl := net.GetActiveNodes()[0].GetServiceUrl(&node.OperaWsService)
	if rpcUrl == nil {
		t.Fatal("websocket service is not available")
	}

	rpcClient, err := network.RetryReturn(network.DefaultRetryAttempts, 1*time.Second, func() (*ethclient.Client, error) {
		return ethclient.Dial(string(*rpcUrl))
	})
	if err != nil {
		t.Fatal("unable to connect the the rpc")
	}

	primaryAccount, err := app.NewAccount(0, PrivateKey, FakeNetworkID)
	if err != nil {
		t.Fatal(err)
	}

	application, err := app.NewERC20Application(rpcClient, primaryAccount, 1)
	if err != nil {
		t.Fatal(err)
	}

	constantShaper := shaper.NewConstantShaper(30.0) // 30 txs/sec

	numGenerators := 5 // 5 parallel workers
	wg := &sync.WaitGroup{}
	app, err := controller.NewAppController(application, constantShaper, numGenerators, net, wg)
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

package controller_test

import (
	"context"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
	"time"
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

	rpcClient, err := ethclient.Dial(string(*rpcUrl))
	if err != nil {
		t.Fatal("unable to connect the the rpc")
	}

	primaryAccount, err := app.NewAccount(PrivateKey, FakeNetworkID)
	if err != nil {
		t.Fatal(err)
	}

	generatorFactory, err := app.NewERC20Application(rpcClient, primaryAccount)
	if err != nil {
		t.Fatal(err)
	}

	constantShaper := shaper.NewConstantShaper(30.0) // 30 txs/sec

	app, err := controller.NewAppController(generatorFactory, constantShaper, 5, net.GetTxsChannel(), rpcClient) // 5 parallel workers
	if err != nil {
		t.Fatal(err)
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
	countInChain, err := generatorFactory.GetAmountOfReceivedTxs(rpcClient)
	if err != nil {
		t.Fatal(err)
	}

	countSent := generatorFactory.GetAmountOfSentTxs()
	if countInChain != countSent {
		t.Errorf("amount of txs in chain (%d) does not match the sent amount (%d)", countInChain, countSent)
	}

	// in optimal case should be generated 30 txs per second
	// as a tolerance for slow CI we require at least 20 txs
	if countInChain < 20 || countInChain > 30 {
		t.Errorf("unexpected amount of generated txs: %d", countInChain)
	}

	txsCounter, ok := app.GetTransactionCounts()
	if !ok {
		t.Errorf("cannot get txs counter")
	}

	if got, err := txsCounter.GetAmountOfReceivedTxs(rpcClient); err != nil || got != countInChain {
		t.Errorf("number of transactions do not match: %d != %d", got, countInChain)
	}

	if got := txsCounter.GetAmountOfSentTxs(); err != nil || got != countSent {
		t.Errorf("number of transactions do not match: %d != %d", got, countSent)
	}

}

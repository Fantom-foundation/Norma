package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/ethereum/go-ethereum/ethclient"
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

	rpcUrl := net.GetActiveNodes()[0].GetServiceUrl(&node.OperaWsService)
	if rpcUrl == nil {
		t.Fatal("websocket service is not available")
	}

	rpcClient, err := ethclient.Dial(string(*rpcUrl))
	if err != nil {
		t.Fatal("unable to connect the the rpc")
	}

	primaryAccount, err := app.NewAccount(0, PrivateKey, FakeNetworkID)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Counter", func(t *testing.T) {
		counterApp, err := app.NewCounterApplication(rpcClient, primaryAccount, 1)
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, counterApp, rpcClient)
	})
	t.Run("ERC20", func(t *testing.T) {
		erc20app, err := app.NewERC20Application(rpcClient, primaryAccount, 1)
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, erc20app, rpcClient)
	})
}

func testGenerator(t *testing.T, app app.Application, rpcClient *ethclient.Client) {
	gen, err := app.CreateGenerator(rpcClient)
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
}

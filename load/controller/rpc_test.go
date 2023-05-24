package controller_test

import (
	"context"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/controller"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
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

	rpcUrl := net.GetActiveNodes()[0].GetHttpServiceUrl(&node.OperaRpcService)

	rpcClient, err := ethclient.Dial(string(*rpcUrl))
	if err != nil {
		t.Fatalf("failed to connecting testing Opera %s: %s", *rpcUrl, err)
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	counterGenerator, err := generator.NewCounterTransactionGenerator(privateKey, big.NewInt(FakeNetworkID))
	if err != nil {
		t.Fatalf("failed to create generator: %s", err)
	}

	constantShaper := shaper.NewConstantShaper(5.0) // 5 txs/sec

	app := controller.NewAppController(counterGenerator, constantShaper, rpcClient)
	err = app.Init()
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

	// get amount of txs applied to the chain
	count, err := counterGenerator.GetAmountOfReceivedTxs()
	if err != nil {
		t.Fatal(err)
	}

	// in optimal case should be generated 5 txs per second
	// as a tolerance for slow CI we require at least 2 txs
	if count < 2 || count > 5 {
		t.Errorf("unexpected amount of generated txs: %d (expected 2-5)", count)
	}
}

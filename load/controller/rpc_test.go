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

	primaryPrivateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	rpcClientForInit, err := ethclient.Dial(string(*rpcUrl))
	if err != nil {
		t.Fatal(err)
	}
	defer rpcClientForInit.Close()

	generatorFactory, err := generator.NewCounterGeneratorFactory(rpcClientForInit, primaryPrivateKey, big.NewInt(FakeNetworkID))
	if err != nil {
		t.Fatal(err)
	}

	constantShaper := shaper.NewConstantShaper(30.0) // 30 txs/sec

	app := controller.NewAppController(generatorFactory, constantShaper, func() (*ethclient.Client, error) {
		// create an ethclient for each worker, prevent bottlenecks in one connection used by multiple workers
		return ethclient.Dial(string(*rpcUrl))
	}, 5) // 5 parallel workers
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

	time.Sleep(2 * time.Second) // wait for txs in TxPool

	// get amount of txs applied to the chain
	countInChain, err := generatorFactory.GetAmountOfReceivedTxs(rpcClientForInit)
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
}

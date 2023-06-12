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

	rpcUrl := net.GetActiveNodes()[0].GetServiceUrl(&node.OperaWsService)
	if rpcUrl == nil {
		t.Fatal("websocket service is not available")
	}

	primaryPrivateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	generatorFactory, err := generator.NewCounterGeneratorFactory(generator.URL(*rpcUrl), primaryPrivateKey, big.NewInt(FakeNetworkID))
	if err != nil {
		t.Fatal(err)
	}

	constantShaper := shaper.NewConstantShaper(30.0) // 30 txs/sec

	app, err := controller.NewAppController(generatorFactory, constantShaper, 5) // 5 parallel workers
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

	// add an RPC node, which has been added when chain already contained txs
	_, err = net.CreateNode(&driver.NodeConfig{
		Name: "RPC-Later-Added",
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second) // wait for txs in TxPool

	// get amount of txs applied to the chain
	countInChain, err := generatorFactory.GetAmountOfReceivedTxs()
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

	if got, err := txsCounter.GetAmountOfReceivedTxs(); err != nil || got != countInChain {
		t.Errorf("number of transactions do not match: %d != %d", got, countInChain)
	}

	if got := txsCounter.GetAmountOfSentTxs(); err != nil || got != countSent {
		t.Errorf("number of transactions do not match: %d != %d", got, countSent)
	}

	checkLatestBlockMatches(t, net.GetActiveNodes())
}

// checkLatestBlockMatches compares hashes of the latest block and the state root across given nodes
func checkLatestBlockMatches(t *testing.T, nodes []driver.Node) {
	var latestBlock *node.BlockHeader
	var err error
	for _, n := range nodes {
		operaNode := n.(*node.OperaNode)
		if latestBlock == nil {
			latestBlock, err = operaNode.GetBlock("latest")
			if err != nil {
				t.Fatalf("failed to get block; %v", err)
			}
		} else {
			block, err := operaNode.GetBlock(latestBlock.Number.String())
			if err != nil {
				t.Fatalf("failed to get block %s; %v", latestBlock.Number.ToInt().String(), err)
			}
			if block.Number.ToInt().Cmp(latestBlock.Number.ToInt()) != 0 {
				t.Errorf("unexpected block hash on other node: %s != %s", block.Number, latestBlock.Number)
			}
			if block.Hash != latestBlock.Hash {
				t.Errorf("unexpected block hash on other node: %s != %s", block.Hash, latestBlock.Hash)
			}
			if block.StateRoot != latestBlock.StateRoot {
				t.Errorf("unexpected state root on other node: %s != %s", block.StateRoot, latestBlock.StateRoot)
			}
		}
	}
}

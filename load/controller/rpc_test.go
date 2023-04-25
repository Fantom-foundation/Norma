package controller

import (
	"context"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
	"time"
)

const TestingRpcUrl = "http://localhost:18545"
const PrivateKey = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7" // Fakenet validator 1
const FakeNetworkID = 0xfa3

func TestTrafficGenerating(t *testing.T) {
	t.Skip("Test requires locally running Opera - skipping until docker will be available")

	rpcClient, err := ethclient.Dial(TestingRpcUrl)
	if err != nil {
		t.Fatalf("failed to connecting testing Opera %s: %s", TestingRpcUrl, err)
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

	sourceDriver := NewAppController(counterGenerator, constantShaper, rpcClient)
	err = sourceDriver.Init()
	if err != nil {
		t.Fatal(err)
	}

	// let the sourceDriver run for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// note: Run is supposed to run in a new thread
	err = sourceDriver.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

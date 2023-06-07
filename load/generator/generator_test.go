package generator_test

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
	"time"
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

	primaryPrivateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Counter", func(t *testing.T) {
		counterGeneratorFactory, err := generator.NewCounterGeneratorFactory(generator.URL(*rpcUrl), primaryPrivateKey, big.NewInt(FakeNetworkID))
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, counterGeneratorFactory)
	})

	t.Run("ERC20", func(t *testing.T) {
		erc20GeneratorFactory, err := generator.NewERC20GeneratorFactory(generator.URL(*rpcUrl), primaryPrivateKey, big.NewInt(FakeNetworkID))
		if err != nil {
			t.Fatal(err)
		}
		testGenerator(t, erc20GeneratorFactory)
	})
}

func testGenerator(t *testing.T, factory generator.TransactionGeneratorFactoryWithStats) {
	gen, err := factory.Create()
	if err != nil {
		t.Fatal(err)
	}
	defer gen.Close()

	for i := 0; i < 10; i++ {
		err = gen.SendTx()
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second) // wait for txs in TxPool

	countSent := factory.GetAmountOfSentTxs()
	if countSent != 10 {
		t.Errorf("unexpected amount of txs sent (%d)", countSent)
	}

	countInChain, err := factory.GetAmountOfReceivedTxs()
	if err != nil {
		t.Fatal(err)
	}
	if countInChain != 10 {
		t.Errorf("unexpected amount of txs in chain (%d)", countInChain)
	}
}

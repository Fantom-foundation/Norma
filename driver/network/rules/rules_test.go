package rules_test

import (
	"context"
	"flag"
	"fmt"
	"github.com/Fantom-foundation/Norma/common/transact"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/network/local"
	"github.com/Fantom-foundation/Norma/driver/network/rules"
	"github.com/Fantom-foundation/Norma/driver/network/rules/abi"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"testing"
	"time"
)

// run with --deepNetworkRulesTest=true flag to check if the rules are applied to the chain (takes more than 10 minutes)
var deepTestingFlag = flag.String("deepNetworkRulesTest", "false", "run more time consuming tests")

const PrivateKey = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7" // Fakenet validator 1
const FakeNetworkID = 0xfa3
const TestingPatchJson = "{\"Dag\":{\"MaxFreeParents\":123}}"

var nodeDriverAddress = common.HexToAddress("0xd100a01e00000000000000000000000000000000")

var newNetworkRules = rules.NetworkRules{
	Dag: &rules.DagRules{
		MaxFreeParents: 123,
	},
}

func TestSettingNetworkRules(t *testing.T) {
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
	t.Cleanup(func() { rpcClient.Close() })

	initialBlockNumber, err := rpcClient.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ownerAccount, err := transact.NewAccount(PrivateKey, FakeNetworkID)
	if err != nil {
		t.Fatal(err)
	}

	err = rules.SetNetworkRules(rpcClient, ownerAccount, newNetworkRules)
	if err != nil {
		t.Fatal(err)
	}

	if err := waitUntilEventOnChain(rpcClient, initialBlockNumber, TestingPatchJson); err != nil {
		t.Fatalf("failed to wait for UpdateNetworkRules event; %v", err)
	}

	if *deepTestingFlag == "true" { // run with --deepTest=true
		// make sure the change is applied - consumes more than 10 minutes, until a new epoch
		if err := waitUntilRulesChanged(string(*rpcUrl), newNetworkRules); err != nil {
			t.Fatalf("failed to wait for NetworkRules change; %v", err)
		}
	}
}

func waitUntilEventOnChain(rpcClient transact.RpcClient, startBlock uint64, awaitedPatchJson string) error {
	driverContract, err := abi.NewNodeDriver(nodeDriverAddress, rpcClient)
	if err != nil {
		return fmt.Errorf("failed to get NodeDriver contract representation; %v", err)
	}
	for i := 0; i < 100; i++ {
		iterator, err := driverContract.FilterUpdateNetworkRules(&bind.FilterOpts{Start: startBlock})
		if err != nil {
			return fmt.Errorf("failed to filter UpdateNetworkRules events; %v", err)
		}
		for iterator.Next() {
			if string(iterator.Event.Diff) == awaitedPatchJson {
				return nil // succeed
			}
		}
		if err := iterator.Error(); err != nil {
			return fmt.Errorf("failed to iterate events; %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("expected event not on chain before timeout")
}

func waitUntilRulesChanged(rpcUrl string, expected rules.NetworkRules) error {
	connection, err := rpc.DialContext(context.Background(), rpcUrl)
	if err != nil {
		return err
	}
	defer connection.Close()

	var result rules.NetworkRules
	for i := 0; i < 1000; i++ {
		err = connection.CallContext(context.Background(), &result, "ftm_getRules", "latest")
		if err != nil {
			return fmt.Errorf("failed to obtain network rules from the network; %v", err)
		}
		if result.Dag.MaxFreeParents == expected.Dag.MaxFreeParents {
			return nil // succeed
		}
		fmt.Printf("waiting until new rules are applied...\n")
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("network rules not set before timeout; result = %v", result)
}

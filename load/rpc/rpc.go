package rpc

import (
	"context"
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func StartRpcWorkers(newNode driver.Node, txs <-chan *types.Transaction) error {
	rpcUrl := newNode.GetServiceUrl(&node.OperaWsService)
	if rpcUrl == nil {
		return fmt.Errorf("websocket service is not available for node %s", newNode.GetLabel())
	}

	for i := 0; i < 10; i++ {
		go runRpcSenderLoop(*rpcUrl, txs)
	}
	return nil
}

func runRpcSenderLoop(rpcUrl driver.URL, txs <-chan *types.Transaction) {
	rpcClient, err := ethclient.Dial(string(rpcUrl))
	if err != nil {
		log.Printf("failed to open RPC connection; %v", err)
	}
	defer rpcClient.Close()

	for tx := range txs {
		err := rpcClient.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Printf("failed to send tx; %v", err)
		}
	}
}

package rpc

import (
	"context"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type RpcWorkerPool struct {
	txs chan *types.Transaction
}

func NewRpcWorkerPool() *RpcWorkerPool {
	return &RpcWorkerPool{
		txs: make(chan *types.Transaction),
	}
}

func (p RpcWorkerPool) SendTransaction(tx *types.Transaction) {
	p.txs <- tx
}

func (p RpcWorkerPool) AfterNodeCreation(newNode driver.Node) {
	rpcUrl := newNode.GetServiceUrl(&node.OperaWsService)
	if rpcUrl == nil {
		return
	}
	for i := 0; i < 10; i++ {
		go runRpcSenderLoop(*rpcUrl, p.txs)
	}
}

func (p RpcWorkerPool) AfterApplicationCreation(application driver.Application) {
	// ignored
}

func (p RpcWorkerPool) Close() error {
	close(p.txs)
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

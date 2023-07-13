package rpc

import (
	"context"
	"github.com/Fantom-foundation/Norma/driver/network"
	"log"
	"sync"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/node"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RpcWorkerPool struct {
	txs    chan *types.Transaction
	done   *sync.WaitGroup
	closed bool
}

func NewRpcWorkerPool() *RpcWorkerPool {
	return &RpcWorkerPool{
		txs:  make(chan *types.Transaction),
		done: &sync.WaitGroup{},
	}
}

func (p *RpcWorkerPool) SendTransaction(tx *types.Transaction) {
	p.txs <- tx
}

func (p *RpcWorkerPool) AfterNodeCreation(newNode driver.Node) {
	if p.closed {
		return
	}

	rpcUrl := newNode.GetServiceUrl(&node.OperaWsService)
	if rpcUrl == nil {
		return
	}
	for i := 0; i < 150; i++ {
		p.done.Add(1)
		go func() {
			defer p.done.Done()
			if err := p.runRpcSenderLoop(*rpcUrl, network.DefaultRetryAttempts); err != nil {
				log.Printf("failed to open RPC connection; %v", err)
				return
			}
		}()
	}
}

func (p *RpcWorkerPool) AfterApplicationCreation(application driver.Application) {
	// ignored
}

func (p *RpcWorkerPool) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	close(p.txs)
	log.Printf("waiting for worker pool to close")
	p.done.Wait()
	log.Printf("worker pool has closed")
	return nil
}

func (p *RpcWorkerPool) runRpcSenderLoop(rpcUrl driver.URL, connectAttempts int) error {
	rpcClient, err := network.RetryReturn(connectAttempts, 1*time.Second, func() (*ethclient.Client, error) {
		return ethclient.Dial(string(rpcUrl))
	})

	if err != nil {
		return err
	}

	defer rpcClient.Close()

	for tx := range p.txs {
		err := rpcClient.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Printf("failed to send tx; %v", err)
		}
	}

	return nil
}

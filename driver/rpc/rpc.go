package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

//go:generate mockgen -source rpc.go -destination rpc_mock.go -package rpc

type RpcClient interface {
	bind.ContractBackend
	Call(result interface{}, method string, args ...interface{}) error
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	Close()
}

func WrapRpcClient(rpcClient *rpc.Client) *RpcClientImpl {
	return &RpcClientImpl{
		Client:    ethclient.NewClient(rpcClient),
		RpcClient: rpcClient,
	}
}

type RpcClientImpl struct {
	*ethclient.Client
	RpcClient *rpc.Client
}

func (r RpcClientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return r.RpcClient.Call(result, method, args...)
}

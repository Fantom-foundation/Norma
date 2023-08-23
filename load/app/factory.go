package app

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"strings"
)

type FactoryFunc func(rpcClient rpc.RpcClient, primaryAccount *Account, numUsers int, feederId, appId uint32) (Application, error)

func NewApplication(appType string, rpcClient rpc.RpcClient, primaryAccount *Account, numUsers int, feederId, appId uint32) (Application, error) {
	if factory := getFactory(appType); factory != nil {
		return factory(rpcClient, primaryAccount, numUsers, feederId, appId)
	}
	return nil, fmt.Errorf("unknown application type '%s'", appType)
}

func IsSupportedApplicationType(appType string) bool {
	return getFactory(appType) != nil
}

func getFactory(appType string) FactoryFunc {
	switch strings.ToLower(appType) {
	case "erc20":
		return NewERC20Application
	case "counter", "":
		return NewCounterApplication
	case "store":
		return NewStoreApplication
	case "uniswap":
		return NewUniswapApplication
	}
	return nil
}

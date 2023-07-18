package app

import (
	"fmt"
	"strings"
)

func NewApplication(appType string, rpcClient RpcClient, primaryAccount *Account, numUsers int) (Application, error) {
	if factory := getFactory(appType); factory != nil {
		return factory(rpcClient, primaryAccount, numUsers)
	}
	return nil, fmt.Errorf("unknown application type '%s'", appType)
}

func IsSupportedApplicationType(appType string) bool {
	return getFactory(appType) != nil
}

func getFactory(appType string) func(RpcClient, *Account, int) (Application, error) {
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

package app

import (
	"fmt"
	"strings"
)

func NewApplication(name string, rpcClient RpcClient, primaryAccount *Account, numUsers int) (Application, error) {
	switch strings.ToLower(name) {
	case "erc20":
		return NewERC20Application(rpcClient, primaryAccount, numUsers)
	case "counter", "":
		return NewCounterApplication(rpcClient, primaryAccount, numUsers)
	}
	return nil, fmt.Errorf("unknown application '%s'", name)
}

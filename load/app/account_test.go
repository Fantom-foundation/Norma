package app

import (
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"
	"math/big"
	"testing"
)

func TestAccount_CreateAccount_AccountsUniq(t *testing.T) {
	chainId := big.NewInt(0xFA)
	const loops = 100

	ctrl := gomock.NewController(t)
	rpcClient := rpc.NewMockRpcClient(ctrl)
	rpcClient.EXPECT().NonceAt(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(uint64(0), nil)

	accounts := make(map[common.Address]struct{}, loops)
	for i := 0; i < loops; i++ {
		gen, err := NewAccountFactory(chainId, 0, uint32(i))
		if err != nil {
			t.Fatalf("cannot create account factory: %v", err)
		}

		for j := 0; j < loops; j++ {
			account, err := gen.CreateAccount(rpcClient)
			if err != nil {
				t.Fatalf("cannot create account: %v", err)
			}

			if _, ok := accounts[account.address]; ok {
				t.Errorf("account address %v is not unique", account.address)
			}

			accounts[account.address] = struct{}{}
		}
	}

}

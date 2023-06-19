package controller

import (
	"context"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestMockedTrafficGenerating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var demoTx types.Transaction

	workers := 2
	mockedGenerator := app.NewMockTransactionGenerator(mockCtrl)

	mockedRpcClient := app.NewMockRpcClient(mockCtrl)
	mockedRpcClient.EXPECT().Close()

	mockedNetwork := driver.NewMockNetwork(mockCtrl)
	mockedNetwork.EXPECT().DialRandomRpc().Return(mockedRpcClient, nil)

	mockedApp := app.NewMockApplication(mockCtrl)
	mockedApp.EXPECT().CreateGenerator(mockedRpcClient).Return(mockedGenerator, nil).Times(workers)
	mockedApp.EXPECT().WaitUntilApplicationIsDeployed(mockedRpcClient).Return(nil)

	// app should be called 10-times to generate 10 txs
	mockedGenerator.EXPECT().GenerateTx().Return(&demoTx, nil).MinTimes(5).MaxTimes(11)
	// network should be called 10-times to send 10 txs
	mockedNetwork.EXPECT().SendTransaction(&demoTx).MinTimes(5).MaxTimes(11)

	// use constant shaper
	constantShaper := shaper.NewConstantShaper(100) // 100 txs/sec

	appController, err := NewAppController(mockedApp, constantShaper, workers, mockedNetwork)
	if err != nil {
		t.Fatal(err)
	}

	// let the app run for 100 ms - should give 10 txs
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// note: Run is supposed to run in a new thread
	err = appController.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

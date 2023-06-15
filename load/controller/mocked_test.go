package controller

import (
	"context"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestMockedTrafficGenerating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var txsChan = make(chan *types.Transaction, 100)
	var rpcClient ethclient.Client
	var demoTx types.Transaction

	workers := 2
	mockedGenerator := app.NewMockTransactionGenerator(mockCtrl)

	mockedApp := app.NewMockApplication(mockCtrl)
	mockedApp.EXPECT().CreateGenerator(&rpcClient).Return(mockedGenerator, nil).Times(workers)
	mockedApp.EXPECT().WaitUntilGeneratorsCreated(&rpcClient).Return(nil)

	// app should be called 10-times to send 10 txs
	mockedGenerator.EXPECT().GenerateTx().Return(&demoTx, nil).MinTimes(5).MaxTimes(11)

	// use constant shaper
	constantShaper := shaper.NewConstantShaper(100) // 100 txs/sec

	appController, err := NewAppController(mockedApp, constantShaper, workers, txsChan, &rpcClient)
	if err != nil {
		t.Fatal(err)
	}

	// let the app run for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// note: Run is supposed to run in a new thread
	err = appController.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(txsChan) < 5 || len(txsChan) > 11 {
		t.Errorf("invalid amount of txs in the channel: %d", len(txsChan))
	}
}

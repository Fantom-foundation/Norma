package controller

import (
	"context"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestMockedTrafficGenerating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	workers := 2

	// generator should be initialized first and then called 10-times to send 10 txs
	mockedGenerator := generator.NewMockTransactionGenerator(mockCtrl)

	mockedGeneratorFactory := generator.NewMockTransactionGeneratorFactory(mockCtrl)
	mockedGeneratorFactory.EXPECT().Create(gomock.Any()).Return(mockedGenerator, nil).Times(workers)
	mockedGenerator.EXPECT().SendTx().Return(nil).MinTimes(9).MaxTimes(11)
	mockedGenerator.EXPECT().Close().Return(nil).Times(workers)

	// use constant shaper
	constantShaper := shaper.NewConstantShaper(100) // 100 txs/sec

	sourceDriver := NewAppController(mockedGeneratorFactory, constantShaper, func() (*ethclient.Client, error) { return nil, nil }, workers)
	err := sourceDriver.Init()
	if err != nil {
		t.Fatal(err)
	}

	// let the sourceDriver run for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// note: Run is supposed to run in a new thread
	err = sourceDriver.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond) // wait for Closes
}

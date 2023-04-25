package controller

import (
	"context"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestMockedTrafficGenerating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// generator should be initialized first and then called 10-times to send 10 txs
	mockedGenerator := generator.NewMockTransactionGenerator(mockCtrl)
	gomock.InOrder(
		mockedGenerator.EXPECT().Init(nil).Return(nil),
		mockedGenerator.EXPECT().SendTx().Return(nil).MinTimes(40).MaxTimes(60),
	)

	// use constant shaper
	constantShaper := shaper.NewConstantShaper(100) // 100 txs/sec

	sourceDriver := NewAppController(mockedGenerator, constantShaper, nil)
	err := sourceDriver.Init()
	if err != nil {
		t.Fatal(err)
	}

	// let the sourceDriver run for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// note: Run is supposed to run in a new thread
	err = sourceDriver.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

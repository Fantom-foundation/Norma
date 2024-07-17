// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"context"
	"testing"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/rpc"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/Fantom-foundation/Norma/load/shaper"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/mock/gomock"
)

func TestMockedTrafficGenerating(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var demoTx types.Transaction

	numUsers := 2
	mockedGenerator := app.NewMockUser(mockCtrl)

	mockedRpcClient := rpc.NewMockRpcClient(mockCtrl)
	mockedRpcClient.EXPECT().Close()

	mockedNetwork := driver.NewMockNetwork(mockCtrl)
	mockedNetwork.EXPECT().DialRandomRpc().Return(mockedRpcClient, nil)

	mockedApp := app.NewMockApplication(mockCtrl)
	mockedApp.EXPECT().CreateUser(mockedRpcClient).Return(mockedGenerator, nil).Times(numUsers)
	mockedApp.EXPECT().WaitUntilApplicationIsDeployed(mockedRpcClient).Return(nil)

	// app should be called 10-times to generate 10 txs
	mockedGenerator.EXPECT().GenerateTx(gomock.Any()).Return(&demoTx, nil).MinTimes(5).MaxTimes(11)
	// network should be called 10-times to send 10 txs
	mockedNetwork.EXPECT().SendTransaction(&demoTx).MinTimes(5).MaxTimes(11)

	// use constant shaper
	constantShaper := shaper.NewConstantShaper(100) // 100 txs/sec

	appController, err := NewAppController(mockedApp, constantShaper, numUsers, nil, mockedNetwork)
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

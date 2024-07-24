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

package appmon

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"go.uber.org/mock/gomock"
)

func TestReceivedTransactionSensorReportsProperValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []uint64{4, 8, 16, 23}
	factory := &receivedTransactionsSensorFactory{}
	for _, expected := range tests {
		application := driver.NewMockApplication(ctrl)
		application.EXPECT().GetReceivedTransactions().Return(expected, nil)

		sensor, err := factory.CreateSensor(application)
		if err != nil {
			t.Fatalf("creation of sensor failed: %v", err)
		}
		if res, err := sensor.ReadValue(); err != nil || res != int(expected) {
			t.Errorf("sensor fetched wrong value, wanted %d, got %d, err %v", expected, res, err)
		}
	}

}

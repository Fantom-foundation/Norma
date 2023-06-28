package appmon

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/golang/mock/gomock"
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

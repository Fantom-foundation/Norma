package accountmon

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/golang/mock/gomock"
)

func TestSentTransactionSensorReportsProperValue(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		account int
		count   uint64
	}{
		{0, 4},
		{1, 5},
		{2, 3},
		{4, 0},
	}

	factory := &sentTransactionsSensorFactory{}

	for _, test := range tests {
		application := driver.NewMockApplication(ctrl)
		application.EXPECT().GetSentTransactions(test.account).Return(test.count, nil)

		sensor, err := factory.CreateSensor(application, test.account)
		if err != nil {
			t.Fatalf("creation of sensor failed: %v", err)
		}
		if res, err := sensor.ReadValue(); err != nil || res != int(test.count) {
			t.Errorf("sensor fetched wrong value, wanted %d, got %d, err %v", test.count, res, err)
		}
	}

}

package accountmon

import (
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/golang/mock/gomock"
)

func TestSentTransactionSensorReportsProperValue(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		account  int
		counts   []uint64
		expected int
	}{
		{0, []uint64{4, 5, 3}, 4},
		{1, []uint64{4, 5, 3}, 5},
		{2, []uint64{4, 5, 3}, 3},
		{4, []uint64{4, 5, 3}, 0},
	}

	factory := &sentTransactionsSensorFactory{}

	for _, test := range tests {
		application := driver.NewMockApplication(ctrl)
		application.EXPECT().GetTransactionCounts().Return(app.TransactionCounts{
			SentTxs: test.counts,
		}, nil)

		sensor, err := factory.CreateSensor(application, test.account)
		if err != nil {
			t.Fatalf("creation of sensor failed: %v", err)
		}
		if res, err := sensor.ReadValue(); err != nil || res != test.expected {
			t.Errorf("sensor fetched wrong value, wanted %d, got %d, err %v", test.expected, res, err)
		}
	}

}

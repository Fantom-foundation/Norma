package appmon

import (
	"fmt"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
)

var (
	// ReceivedTransactions is a metric capturing the number of transactions actually received by an application.
	// It captures all transactions that went into the application, i.e. this number removes all rejected transactions
	// at any level of the network.
	// The key of the series is the number of accounts sending the transactions and the value is the number received transactions.
	ReceivedTransactions = monitoring.Metric[monitoring.App, monitoring.Series[monitoring.Time, int]]{
		Name:        "ReceivedTransactions",
		Description: "The number of transactions actually received by an application over time",
	}
)

func init() {
	if err := monitoring.RegisterSource(ReceivedTransactions, newReceivedTransactionsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// newReceivedTransactionsSource is an internal factory for the ReceivedTranactions metric.
func newReceivedTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.App, monitoring.Series[monitoring.Time, int]] {
	return NewPeriodicAppDataSource[int](ReceivedTransactions, monitor, &receivedTransactionsSensorFactory{})
}

type receivedTransactionsSensorFactory struct{}

func (f *receivedTransactionsSensorFactory) CreateSensor(app driver.Application) (utils.Sensor[int], error) {
	return &receivedTransactionsSensor{
		app: app,
	}, nil
}

type receivedTransactionsSensor struct {
	app driver.Application
}

func (s *receivedTransactionsSensor) ReadValue() (int, error) {
	counts, err := s.app.GetTransactionCounts()
	if err != nil {
		return 0, err
	}
	return int(counts.ReceivedTxs), nil
}

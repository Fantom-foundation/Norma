package accountmon

import (
	"fmt"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
)

var (
	// SentTransactions is a metric capturing the number of transactions sent to an application on a per account granularity.
	// It captures the number of transactions sent to the network, i.e. all attempts to feed the application
	// with transactions, which does not necessarily mean the application has received all of them.
	// The key of the series is the number of accounts sending the transactions and the value is the number of sent transactions.
	SentTransactions = monitoring.Metric[monitoring.Account, monitoring.Series[monitoring.Time, int]]{
		Name:        "SentTransactions",
		Description: "The number of transactions attempted to be sent to an application",
	}
)

func init() {
	if err := monitoring.RegisterSource(SentTransactions, newSentTransactionsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// newSentTransactionsSource is an internal factory for the ReceivedTranactions metric.
func newSentTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.Account, monitoring.Series[monitoring.Time, int]] {
	return NewPeriodicAccountDataSource[int](SentTransactions, monitor, &sentTransactionsSensorFactory{})
}

type sentTransactionsSensorFactory struct{}

func (f *sentTransactionsSensorFactory) CreateSensor(app driver.Application, account int) (utils.Sensor[int], error) {
	return &sentTransactionsSensor{
		app:     app,
		account: account,
	}, nil
}

type sentTransactionsSensor struct {
	app     driver.Application
	account int
}

func (s *sentTransactionsSensor) ReadValue() (int, error) {
	counts, err := s.app.GetTransactionCounts()
	if err != nil {
		return 0, err
	}
	if s.account < 0 || s.account >= len(counts.SentTxs) {
		return 0, nil
	}
	return int(counts.SentTxs[s.account]), nil
}

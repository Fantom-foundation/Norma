package app

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/export"
	"github.com/Fantom-foundation/Norma/load/app"
)

var (
	// SentTransactions is a metric capturing the number of transactions sent to an application.
	// It captures the number of transactions sent to the network, i.e. all attempts to feed the application
	// with transactions, which does not necessarily mean the application has received all of them.
	// The key of the series is the number of accounts sending the transactions and the value is the number of sent transactions.
	SentTransactions = monitoring.Metric[monitoring.App, monitoring.Series[int, int]]{
		Name:        "SentTransactions",
		Description: "The number of transactions attempted to be sent to an application",
	}

	// ReceivedTransactions is a metric capturing the number of transactions actually received by an application.
	// It captures all transactions that went into the application, i.e. this number removes all rejected transactions
	// at any level of the network.
	// The key of the series is the number of accounts sending the transactions and the value is the number received transactions.
	ReceivedTransactions = monitoring.Metric[monitoring.App, monitoring.Series[int, int]]{
		Name:        "ReceivedTransactions",
		Description: "The number of transactions actually received by an application",
	}
)

func init() {
	if err := monitoring.RegisterSource(SentTransactions, newSentTransactionsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}

	if err := monitoring.RegisterSource(ReceivedTransactions, newReceivedTransactionsSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// countGetter is a function that returns a value obtained from a transaction counter
type countGetter func(app.TransactionCounts) (int, error)

// TxsCounter allows for metering the number of transactions sent or received on an application (i.e. a smart contract).
// It may happen that some applications are not able to count number of applications they have received.
// In this case, the application is not added in this metric at all.
type TxsCounter struct {
	metric       monitoring.Metric[monitoring.App, monitoring.Series[int, int]]
	monitor      *monitoring.Monitor
	getter       countGetter
	series       map[monitoring.App]*monitoring.SyncedSeries[int, int]
	applications []driver.Application
	seriesLock   *sync.Mutex
}

// NewSentTransactionsSource creates a metric that counts the number of transactions sent to an application.
// This number of transactions was not necessarily  received by the application as the transactions
// could be filtered out by any layers between the RPC endpoint and actual block processing,
// or the client was not able to process requested amount of transactions and the transactions could not reach
// the block processing.
func NewSentTransactionsSource(monitor *monitoring.Monitor) *TxsCounter {
	s := func(c app.TransactionCounts) (int, error) {
		return int(c.SentTxs), nil
	}
	res := newTxsCounterSource(monitor, s, SentTransactions)

	return res
}

// NewReceivedTransactionsSource creates a metric that counts the number of transactions received to an application.
// This number of transactions may be smaller than the number of actually sent transactions
// as the transactions could be filtered out by any layers between the RPC endpoint and actual block processing,
// or the client was not able to process requested amount of transactions and the transactions could not reach
// the block processing.
func NewReceivedTransactionsSource(monitor *monitoring.Monitor) *TxsCounter {
	s := func(c app.TransactionCounts) (int, error) {
		return int(c.ReceivedTxs), nil
	}
	res := newTxsCounterSource(monitor, s, ReceivedTransactions)

	return res
}

// newSentTransactionsSource is the same as its public counterpart, it just returns the Source type.
func newSentTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.App, monitoring.Series[int, int]] {
	return NewSentTransactionsSource(monitor)
}

// newReceivedTransactionsSource is the same as its public counterpart, it just returns the Source type.
func newReceivedTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.App, monitoring.Series[int, int]] {
	return NewReceivedTransactionsSource(monitor)
}

func newTxsCounterSource(monitor *monitoring.Monitor, sensor countGetter, metric monitoring.Metric[monitoring.App, monitoring.Series[int, int]]) *TxsCounter {
	res := &TxsCounter{
		monitor:      monitor,
		metric:       metric,
		getter:       sensor,
		series:       make(map[monitoring.App]*monitoring.SyncedSeries[int, int], 50),
		applications: make([]driver.Application, 0, 50),
		seriesLock:   &sync.Mutex{},
	}

	monitor.Network().RegisterListener(res)

	for _, app := range monitor.Network().GetActiveApplications() {
		res.AfterApplicationCreation(app)
	}

	monitor.Writer().Add(func() error {
		return export.AddAppSeriesSource[int, int](monitor.Writer(), res, export.DirectConverter[int]{}, export.DirectConverter[int]{})
	})

	return res
}

func (s *TxsCounter) GetMetric() monitoring.Metric[monitoring.App, monitoring.Series[int, int]] {
	return s.metric
}

func (s *TxsCounter) GetSubjects() []monitoring.App {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	res := make([]monitoring.App, 0, len(s.series))
	for k := range s.series {
		res = append(res, k)
	}
	return res
}

func (s *TxsCounter) GetData(node monitoring.App) (monitoring.Series[int, int], bool) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	res, exists := s.series[node]
	return res, exists
}

func (s *TxsCounter) Shutdown() error {
	// compute data before shutdown
	// alternatively we would need a listener of application Start()/Stop()
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	var errs []error
	for _, app := range s.applications {
		txs, err := app.GetTransactionCounts()
		if err == nil {
			if val, err := s.getter(txs); err == nil {
				series, exists := s.series[monitoring.App(app.Config().Name)]
				if !exists {
					series = &monitoring.SyncedSeries[int, int]{}
					s.series[monitoring.App(app.Config().Name)] = series
				}
				errs = append(errs, series.Append(app.Config().Accounts, val))
			} else {
				errs = append(errs, err)
			}
		}
	}

	s.applications = s.applications[0:0]

	return errors.Join(errs...)
}

func (s *TxsCounter) AfterNodeCreation(driver.Node) {
	// ignored
}

func (s *TxsCounter) AfterApplicationCreation(app driver.Application) {
	s.seriesLock.Lock()
	defer s.seriesLock.Unlock()

	s.applications = append(s.applications, app)
}

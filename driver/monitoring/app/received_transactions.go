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
	count, err := s.app.GetReceivedTransactions()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

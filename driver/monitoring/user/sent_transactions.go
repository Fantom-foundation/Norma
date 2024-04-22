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

package user

import (
	"fmt"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
)

var (
	// SentTransactions is a metric capturing the number of transactions sent to an application on a per user granularity.
	// It captures the number of transactions sent to the network, i.e. all attempts to feed the application
	// with transactions, which does not necessarily mean the application has received all of them.
	// The key of the series is the number of accounts sending the transactions and the value is the number of sent transactions.
	SentTransactions = monitoring.Metric[monitoring.User, monitoring.Series[monitoring.Time, int]]{
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
func newSentTransactionsSource(monitor *monitoring.Monitor) monitoring.Source[monitoring.User, monitoring.Series[monitoring.Time, int]] {
	return NewPeriodicUserDataSource[int](SentTransactions, monitor, &sentTransactionsSensorFactory{})
}

type sentTransactionsSensorFactory struct{}

func (f *sentTransactionsSensorFactory) CreateSensor(app driver.Application, user int) (utils.Sensor[int], error) {
	return &sentTransactionsSensor{
		app:  app,
		user: user,
	}, nil
}

type sentTransactionsSensor struct {
	app  driver.Application
	user int
}

func (s *sentTransactionsSensor) ReadValue() (int, error) {
	count, err := s.app.GetSentTransactions(s.user)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

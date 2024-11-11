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

package app

import (
	"github.com/Fantom-foundation/Norma/driver/rpc"
)

//go:generate mockgen -source app.go -destination app_mock.go -package app

type Application interface {
	// CreateUsers creates a number of users for this application that
	// can generate transactions for this application.
	CreateUsers(context AppContext, numUsers int) ([]User, error)

	// GetReceivedTransactions returns the total number of transactions
	// received by this application up to the current point in time.
	GetReceivedTransactions(rpcClient rpc.RpcClient) (uint64, error)
}

// User produces a stream of transactions to Generate traffic on the chain.
// Implementations are not required to be thread-safe.
type User interface {
	// SendTransaction requests this user to send a transaction to the chain
	// and wait for its completion.
	SendTransaction(rpcClient rpc.RpcClient) error

	// GetTotalNumberOfSentTransactions returns the number of successfully
	// sent transactions by this user.
	GetTotalNumberOfSentTransactions() uint64
}

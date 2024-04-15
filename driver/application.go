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

package driver

//go:generate mockgen -source application.go -destination application_mock.go -package driver

// Application is an abstraction of an application running on a Norma net.
type Application interface {
	// Start begins producing load on the network as configured for this app.
	Start() error
	// Stop terminates the load production.
	Stop() error

	// Config returns current application configuration.
	Config() *ApplicationConfig

	// GetNumberOfUsers retrieves the number of users interacting with this application.
	// This value is expected to be a constant over the life-time of an application.
	GetNumberOfUsers() int

	// GetSentTransactions returns the number of transactions sent by a given user.
	GetSentTransactions(user int) (uint64, error)

	// GetReceivedTransactions returns the number fo transactions received by the appliation
	// on the network.
	GetReceivedTransactions() (uint64, error)
}

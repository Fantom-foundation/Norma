package driver

import "github.com/Fantom-foundation/Norma/load/app"

//go:generate mockgen -source application.go -destination application_mock.go -package driver

// Application is an abstraction of an application running on a Norma net.
type Application interface {
	// Start begins producing load on the network as configured for this app.
	Start() error
	// Stop terminates the load production.
	Stop() error

	// Config returns current application configuration.
	Config() *ApplicationConfig

	// GetTransactionCounts returns information about expected and received transactions
	// if this information is available for this application.
	// If the information is not available, second argument returns false.
	GetTransactionCounts() (app.TransactionCountsProvider, bool)
}

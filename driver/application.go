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

	// GetNumberOfAccounts retrieves the number of accounts interacting with this application.
	// This value is expected to be a constant over the life-time of an application.
	GetNumberOfAccounts() int

	// GetSentTransactions returns the number of transactions send from a given account.
	GetSentTransactions(account int) (uint64, error)

	// GetReceivedTransactions returns the number fo transactions received by the appliation
	// on the network.
	GetReceivedTransactions() (uint64, error)
}

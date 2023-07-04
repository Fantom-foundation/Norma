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

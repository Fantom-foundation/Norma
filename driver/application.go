package driver

//go:generate mockgen -source application.go -destination application_mock.go -package driver

// Application is an abstraction of an application running on a Norma net.
type Application interface {
	// Start begins producing load on the network as configured for this app.
	Start() error
	// Stop terminates the load production.
	Stop() error
}

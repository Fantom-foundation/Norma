package driver

//go:generate mockgen -source host.go -destination host_mock.go -package driver

type IP string

// TODO: document
// key: this is the host running nodes
type Host interface {
	// IsRunning tests whether this host is running or has been stopped.
	IsRunning() bool
	// GetIP returns this node's IP address.
	GetIP() IP

	GetAddressForService(*ServiceDescription) AddressPort

	// Stop shuts down the services running on the node gracefully, using
	// their regular shutdown procedure (not killed). After stopping the
	// service, no more interactions are expected to succeed.
	Stop() error
	// SaveLogTo transfers the logs of the node ot the given file directory.
	SaveLogTo(directory string) error
	// Cleanup releases all underlying resources. After the cleanup no more
	// operations on this node are expected to succeed.
	Cleanup() error
}

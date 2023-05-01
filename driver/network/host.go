package network

//go:generate mockgen -source host.go -destination host_mock.go -package network

// AddressPort is a string addressing an IP and a port in the format <IP>:<port>.
type AddressPort string

// Host is an execution environment for Nodes. A host could be an actual physical
// machine with its dedicated hardware, a virtual machine with shared resources,
// or a Docker container hosting services. Hosts may be implicitly created by a
// Networks's CreateNode function.
type Host interface {
	// IsRunning tests whether this host is running or has been stopped.
	IsRunning() bool

	// GetAddressForService returns the address of a service running on this
	// host, or nil if such a service is not offered.
	GetAddressForService(*ServiceDescription) *AddressPort

	// Stop shuts down the services running on the host gracefully, using
	// their regular shutdown procedure (not killed). After stopping the
	// service, no more interactions are expected to succeed.
	Stop() error

	// SaveLogTo transfers the logs of the host to the given file directory.
	SaveLogTo(directory string) error

	// Cleanup releases all underlying resources. After the cleanup no more
	// operations on this host are expected to succeed.
	Cleanup() error
}

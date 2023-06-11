package network

import "io"

//go:generate mockgen -source host.go -destination host_mock.go -package network

// AddressPort is a string addressing an IP and a port in the format <IP>:<port>.
type AddressPort string

// Host is an execution environment for Nodes. A host could be an actual physical
// machine with its dedicated hardware, a virtual machine with shared resources,
// or a Docker container hosting services. Hosts may be implicitly created by a
// Networks's CreateNode function.
type Host interface {
	// Hostname returns the hostname of the host.
	Hostname() string

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

	// StreamLog provides a reader that is continuously providing the host log.
	// The log is tailed via the reader and blocked until next log lines are ready.
	// Before tailing new log lines, the reader can provide certain amount of previous log lines.
	// For instance, the docker host uses a ring memory of 150 lines to buffer previous lines.
	// The reader should reach its end (EOF) only when the host/container is stopped or interrupted.
	// If this method is called many times, it should dispatch the log to all returned
	// readers, i.e. all of them see the same output.
	// It is up to the caller to close the stream.
	StreamLog() (io.ReadCloser, error)

	// Cleanup releases all underlying resources. After the cleanup no more
	// operations on this host are expected to succeed.
	Cleanup() error
}

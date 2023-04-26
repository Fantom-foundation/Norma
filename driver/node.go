package driver

//go:generate mockgen -source node.go -destination node_mock.go -package driver

// Node is controlling a single node in a Norma network. It provides abstract
// control of a node, allowing it to be started (through an Environment),
// interact with the node, and shut it down.
type Node interface {
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

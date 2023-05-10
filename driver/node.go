package driver

//go:generate mockgen -source node.go -destination node_mock.go -package driver

// Node is controlling a single node in a Norma network. It provides abstract
// control of a node, allowing it to be started (through an Environment),
// interact with the node, and shut it down.
type Node interface {
	// Returns true if the node is still running, false if stopped.
	IsRunning() bool

	// GetNodeID returns an enode identifying this node within the Norma network.
	// An error shall be produced if no valid node ID could be obtained.
	GetNodeID() (NodeID, error)

	// GetRpcServiceUrl returns the URL of the RPC server running on the
	// represented node. May be nil if no such service is offered.
	GetRpcServiceUrl() *URL

	// Stop shuts down this node gracefully, using its regular shutdown
	// procedure (not killed). After stopping the service, no more interactions
	// are expected to succeed.
	Stop() error

	// Cleanup releases all underlying resources. After the cleanup no more
	// operations on this node are expected to succeed.
	Cleanup() error
}

// NodeID is a unique ID identifying each node. This identifier is used, for
// instance, to connect nodes within the network. In Opera, this ID is known
// as an 'enode' identifier.
type NodeID string

// URL is a mere alias type for a string supposed to encode a URL.
type URL string

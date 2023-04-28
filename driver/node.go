package driver

//go:generate mockgen -source node.go -destination node_mock.go -package driver

// Node is controlling a single node in a Norma network. It provides abstract
// control of a node, allowing it to be started (through an Environment),
// interact with the node, and shut it down.
type Node interface {
	// GetNodeID returns a enode identifying this node within the Norma network.
	GetNodeID() (NodeID, error)

	// GetRpcServiceUrl returns the URL of the RPC serve rrunning on the
	// represented node. May be nil if no such service is offered.
	GetRpcServiceUrl() *URL

	// GetHost retrieves the host this node is running on.
	GetHost() Host
}

// NodeID is a unique ID identifying each node. This identifier is used, for
// instance, to connect nodes within the network. In Opera, this ID is known
// as an 'enode' identifier.
type NodeID string

// URL is a mere alias type for a string supposed to encode a URL.
type URL string

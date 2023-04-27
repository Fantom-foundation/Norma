package driver

//go:generate mockgen -source node.go,host.go -destination node_mock.go -package driver

type AddressPort string

type NodeID string

// Node is controlling a single node in a Norma network. It provides abstract
// control of a node, allowing it to be started (through an Environment),
// interact with the node, and shut it down.
type Node interface {
	// TODO: document
	Host
	// GetNodeID returns a enode identifying this node within the Norma network.
	GetNodeID() (NodeID, error)

	GetAddressForService(ServiceID) AddressPort
}

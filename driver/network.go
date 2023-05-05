package driver

//go:generate mockgen -source network.go -destination network_mock.go -package driver

// Network abstracts an execution environment for running scenarios.
// Implementations may run nodes and applications locally, in docker images, or
// remotely, on actual nodes. The interface is used by the scenario driver
// to execute scenario descriptions.
type Network interface {
	// CreateNode creates a new node instance running a network client based on
	// the given configuration. It is used by the scenario executor to add
	// nodes to the network as needed.
	CreateNode(config *NodeConfig) (Node, error)

	// CreateApplication creates a new application in this network, ready to
	// produce load as defined by its configuration.
	CreateApplication(config *ApplicationConfig) (Application, error)

	// Shutdown stops all applications and nodes in the network and frees
	// any potential other resources.
	Shutdown() error
}

type NodeConfig struct {
	Name string
	// TODO: add other parameters as needed
	//  - features to include on the node
	//  - state DB configuration
	//  - EVM configuration
}

type ApplicationConfig struct {
	Name string

	// Rate defines the Tx/s the source should produce while active.
	Rate float32

	// TODO: add other parameters as needed
	//  - application type
	//  - other traffic shapes
}

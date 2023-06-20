package monitoring

import (
	"io"
	"log"
	"sync"

	"github.com/Fantom-foundation/Norma/driver"
)

//go:generate mockgen -source node_log_provider.go -destination node_log_provider_mock.go -package monitoring

// LogListener gets data of a new block every time it is occurred for a certain node.
// All listeners are executed in sequence, i.e. each processing of a block should be fast
// not to block the loop.
type LogListener interface {

	// OnBlock is triggered every time a new block is found.
	OnBlock(node Node, block Block)
}

// NodeLogProvider is an interface for registering listeners that will be notified about incoming blocks.
type NodeLogProvider interface {

	// RegisterLogListener registers the input listener to receive new blocks.
	RegisterLogListener(listener LogListener)

	// UnregisterLogListener removes the input listener from receiving new events
	UnregisterLogListener(listener LogListener)
}

// NodeLogDispatcher listens and maintains nodes of the network.
// Every time a node is added to the network, the internal list is extended.
// Log streams of all the nodes maintained in this registry are read and parsed,
// while the parsed blocks from the logs are distributed to all registered listeners.
type NodeLogDispatcher struct {
	nodes     map[Node]bool
	nodesLock sync.Mutex

	listeners     map[LogListener]bool
	listenersLock sync.Mutex

	network driver.Network
}

// NewNodeLogDispatcher creates a new instance of this registry, which is filled
// by already running nodes, and further listens to newly added nodes.
func NewNodeLogDispatcher(network driver.Network) *NodeLogDispatcher {
	res := &NodeLogDispatcher{
		network:   network,
		nodes:     make(map[Node]bool, 50),
		listeners: make(map[LogListener]bool, 50),
	}

	// listen for new Nodes
	network.RegisterListener(res)

	// get nodes that have been started before this instance creation
	for _, node := range res.network.GetActiveNodes() {
		res.AfterNodeCreation(node)
	}

	return res
}

func (n *NodeLogDispatcher) RegisterLogListener(listener LogListener) {
	n.listenersLock.Lock()
	defer n.listenersLock.Unlock()
	n.listeners[listener] = true
}

func (n *NodeLogDispatcher) UnregisterLogListener(listener LogListener) {
	n.listenersLock.Lock()
	defer n.listenersLock.Unlock()
	delete(n.listeners, listener)
}

func (n *NodeLogDispatcher) AfterNodeCreation(node driver.Node) {
	nodeId := node.GetLabel()
	n.nodesLock.Lock()
	defer n.nodesLock.Unlock()

	// open new log stream only when the node has not been in the map yet
	if _, exists := n.nodes[Node(nodeId)]; !exists {
		logStream, err := node.StreamLog()
		if err != nil {
			log.Printf("failed to obtain logs of node, will not be able to track blocks: %v", err)
			return // do not start dispatch on error
		}
		n.startDispatcher(Node(nodeId), logStream)

		n.nodes[Node(nodeId)] = true
	}
}

func (n *NodeLogDispatcher) AfterApplicationCreation(driver.Application) {
	// ignored
}

// getNodes returns all nodes so far accumulated in this registry.
func (n *NodeLogDispatcher) getNodes() []Node {
	n.nodesLock.Lock()
	defer n.nodesLock.Unlock()

	res := make([]Node, 0, len(n.nodes))
	for k := range n.nodes {
		res = append(res, k)
	}
	return res
}

// Size returns the count of nodes accumulated in this registry.
func (n *NodeLogDispatcher) getNumNodes() int {
	n.nodesLock.Lock()
	defer n.nodesLock.Unlock()

	return len(n.nodes)
}

func (n *NodeLogDispatcher) startDispatcher(node Node, reader io.ReadCloser) {
	go func() {
		defer func() {
			_ = reader.Close()
		}()
		ch := NewLogReader(reader)
		for b := range ch {
			n.listenersLock.Lock()
			for k := range n.listeners {
				k.OnBlock(node, b)
			}
			n.listenersLock.Unlock()
		}
	}()
}

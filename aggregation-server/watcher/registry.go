package watcher

import "github.com/google/uuid"

type (
	// Node represents a watcher-node.
	// watcher-nodes send file operations to the server.
	Node struct {
		Instance uuid.UUID
		port     int
	}

	// Registry stores a map of nodes that want to send file
	// operations.
	Registry struct {
		nodes map[uuid.UUID]*Node
	}
)

// NewRegistry returns an empty node registry.
func NewRegistry() *Registry {
	return &Registry{
		make(map[uuid.UUID]*Node),
	}
}

// AddNode adds a node to the registry.
func (r *Registry) AddNode(node *Node) {
	r.nodes[node.Instance] = node
}

// RemoveNode removes a node with the given id from the
// regsitry.
func (r *Registry) RemoveNode(id uuid.UUID) {
	delete(r.nodes, id)
}

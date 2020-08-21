package watcher

import (
	"sync"

	"github.com/dawsonalex/aggregator/aggregator"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	addOperation    = "add"
	removeOperation = "remove"
)

type (
	// Node represents a watcher-node.
	// watcher-nodes send file operations to the server.
	Node struct {
		Instance uuid.UUID
		seqno    int
		ops      chan aggregator.Operation
		files    map[string]struct{}
		mux      sync.RWMutex
	}

	// Operation represents an operation that a node can
	// make on a file.
	Operation struct {
		Type     string
		SeqNo    int
		Filename string
	}

	// Registry stores a map of nodes that want to send file
	// operations.
	Registry struct {
		nodes map[uuid.UUID]*Node
		mux   sync.RWMutex
	}
)

// NewRegistry returns an empty node registry.
func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[uuid.UUID]*Node),
	}
}

// AddNode adds a node to the registry.
func (r *Registry) AddNode(id uuid.UUID) {
	if _, nodeExists := r.nodes[id]; !nodeExists {
		log.WithField("ID", id).Println("Adding node")
		r.mux.Lock()
		r.nodes[id] = &Node{
			Instance: id,
			seqno:    -1,
			ops:      make(chan aggregator.Operation),
		}
		r.mux.Unlock()
	}
}

// RemoveNode removes a node with the given id from the
// regsitry.
func (r *Registry) RemoveNode(id uuid.UUID) {
	r.mux.Lock()
	if node, nodeExists := r.nodes[id]; nodeExists {
		close(node.ops)
		delete(r.nodes, id)
	}
	r.mux.Unlock()
}

// Node returns the watcher node with the given id, or nil if the node doesn't exist.
func (r *Registry) Node(id uuid.UUID) *Node {
	r.mux.RLock()
	if node, nodeExists := r.nodes[id]; nodeExists {
		return node
	}
	r.mux.RUnlock()
	return nil
}

// ListFiles returns a slice of filenames held
// by all nodes currently registered.
func (r *Registry) ListFiles() []string {
	files := make([]string, 0)
	for _, node := range r.nodes {
		files = append(files, node.ListFiles()...)
	}
	return files
}

// Do tells a node to send an operation down its operation channel.
func (n *Node) Do(op Operation) {
	// Only carry out the operation if it's the next in sequence.
	if op.SeqNo == n.seqno+1 {
		n.seqno = op.SeqNo

		switch op.Type {
		case addOperation:
			n.files[op.Filename] = struct{}{}
		case removeOperation:
			delete(n.files, op.Filename)
		}
	}
}

// ListFiles lists the files that the node is watching.
func (n *Node) ListFiles() []string {
	files := make([]string, len(n.files))
	n.mux.RLock()
	for k := range n.files {
		files = append(files, k)
	}
	n.mux.RUnlock()
	return files
}

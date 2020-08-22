package watcher

import (
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	addOperation    = "add"
	removeOperation = "remove"

	// NoSequence reqresents a nodes sequence
	// value that's not yet initialised.
	NoSequence = -1
)

type (
	// Node represents a watcher-node.
	// watcher-nodes send file operations to the server.
	Node struct {
		Instance uuid.UUID
		seqno    int
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

// AddNode adds a node to the registry. Returns true if the node
// was added and did not exist before, otherwise returns false.
func (r *Registry) AddNode(id uuid.UUID) (chan string, bool) {
	if _, nodeExists := r.nodes[id]; !nodeExists {
		log.WithField("instance-id", id).Println("Adding node")
		fileMap := make(map[string]struct{})
		node := &Node{
			Instance: id,
			seqno:    NoSequence,
			files:    fileMap,
		}
		r.mux.Lock()
		r.nodes[id] = node
		r.mux.Unlock()

		fileChan := make(chan string)
		go func() {
			for file := range fileChan {
				node.mux.Lock()
				node.files[file] = struct{}{}
				node.mux.Unlock()
			}
		}()
		return fileChan, true
	}
	return nil, false
}

// RemoveNode removes a node with the given id from the
// regsitry.
func (r *Registry) RemoveNode(id uuid.UUID) {
	r.mux.Lock()
	if _, nodeExists := r.nodes[id]; nodeExists {
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
	// Only carry out the operation if it's the next in sequence, or sequence hasn't
	// been set yet.
	if n.seqno == NoSequence || op.SeqNo == n.seqno+1 {
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
	files := make([]string, 0)
	n.mux.RLock()
	for k := range n.files {
		files = append(files, k)
	}
	n.mux.RUnlock()
	return files
}

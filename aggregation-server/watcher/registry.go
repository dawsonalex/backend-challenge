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

// Registry stores a map of nodes that want to send file
// operations.
type Registry struct {
	nodes map[uuid.UUID]*Node
	mux   sync.RWMutex
}

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

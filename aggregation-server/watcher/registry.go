package watcher

import (
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	log   *logrus.Logger
}

// NewRegistry returns an empty node registry.
func NewRegistry(logger *logrus.Logger) *Registry {
	if logger == nil {
		logger = defaultLogger()
	}
	return &Registry{
		nodes: make(map[uuid.UUID]*Node),
		log:   logger,
	}
}

func defaultLogger() *logrus.Logger {
	return log.New()
}

// AddNode adds a node to the registry. Returns true if the node
// was added and did not exist before, otherwise returns false.
// AddNode also returns a channel to directly add filenames to a node via,
// and done channel that will receive a value when all the filenames are read
// from the file channel.
func (r *Registry) AddNode(id uuid.UUID) (chan string, chan struct{}, bool) {
	if _, nodeExists := r.nodes[id]; !nodeExists {
		r.log.WithField("node-id", id).Infoln("Adding node")
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
		done := make(chan struct{})
		go func() {
			for file := range fileChan {
				node.mux.Lock()
				node.files[file] = struct{}{}
				node.mux.Unlock()
			}
			done <- struct{}{}
		}()
		return fileChan, done, true
	}
	return nil, nil, false
}

// RemoveNode removes a node with the given id from the
// regsitry.
func (r *Registry) RemoveNode(id uuid.UUID) {
	r.mux.Lock()
	if _, nodeExists := r.nodes[id]; nodeExists {
		r.log.WithField("node-id", id).Infoln("Removing node")
		delete(r.nodes, id)
	}
	r.mux.Unlock()
}

// Node returns the watcher node with the given id, or nil if the node doesn't exist.
func (r *Registry) Node(id uuid.UUID) *Node {
	r.mux.RLock()
	if node, nodeExists := r.nodes[id]; nodeExists {
		r.mux.RUnlock()
		return node
	}
	return nil
}

// ListFiles returns a slice of filenames held
// by all nodes currently registered.
func (r *Registry) ListFiles() []string {
	files := make([]string, 0)
	for _, node := range r.nodes {
		files = append(files, node.ListFiles()...)
	}
	r.log.Debugln("listing files: ", files)
	return files
}

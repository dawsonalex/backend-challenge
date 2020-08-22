package watcher

import (
	"sync"

	"github.com/google/uuid"
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
)

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

package lib

import (
	"github.com/google/uuid"
)

// FilesResponse is the type sent when a
// watcher-node requests to see the aggregated files.
type FilesResponse struct {
	Files []string `json:"files"`
}

// HelloRequest is the type received
// when a node wishes to register with the aggregator.
type HelloRequest struct {
	Instance uuid.UUID `json:"instance"`
	Port     int       `json:"port"`
}

// ByeRequest is the type received when a
// node wishes to unregister with the aggregator.
type ByeRequest struct {
	Instance string `json:"instance"`
}

// File represents a single file with a filename
type File struct {
	Filename string `json:"filename"`
}

// OperationRequest is the message sent from a watcher.
type OperationRequest struct {
	Instance uuid.UUID `json:"instance"`
	Type     string    `json:"op"`
	SeqNo    int       `json:"seqno"`
	Value    File      `json:"value"`
}

// OperationRequests is a slice of operations.
type OperationRequests []OperationRequest

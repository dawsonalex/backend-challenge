package lib

import (
	"github.com/google/uuid"
)

// File represents a single file with a filename
type File struct {
	Filename string
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

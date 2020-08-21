package lib

import (
	"github.com/dawsonalex/aggregator/aggregator"
	"github.com/google/uuid"
)

// OperationRequest is the message sent from a watcher.
type OperationRequest struct {
	Instance uuid.UUID       `json:"instance"`
	Type     string          `json:"op"`
	SeqNo    int             `json:"seqno"`
	Value    aggregator.File `json:"value"`
}

// OperationRequests is a slice of operations.
type OperationRequests []OperationRequest

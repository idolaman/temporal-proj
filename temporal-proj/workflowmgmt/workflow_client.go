package workflowmgmt

import "context"

// Any workflow engine implementation (e.g. Temporal) should satisfy this interface.
type WorkflowClient interface {
	// Close closes the workflow client connection
	Close()
	// Start starts a new workflow execution
	Start(ctx context.Context, workflowID string, task interface{}) error
	// GetResult retrieves the result of a workflow execution
	GetResult(ctx context.Context, workflowID string, result interface{}) error
}

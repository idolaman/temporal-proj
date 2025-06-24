package temporal

import "context"

type WorkflowExecutor interface {
	Start(ctx context.Context, workflowID string, task interface{}) error
	GetResult(ctx context.Context, workflowID string, result interface{}) error
}

func (tc *Client) Start(ctx context.Context, workflowID string, task interface{}) error {
	_, err := tc.c.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: tc.taskQueue,
	}, tc.workflowName, task)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}
	log.Printf("Started workflow %s", workflowID)
	return nil
}

func (tc *Client) GetResult(ctx context.Context, workflowID string, result interface{}) error {
	run := tc.c.GetWorkflow(ctx, workflowID, "")
	if err := run.Get(ctx, result); err != nil {
		return fmt.Errorf("failed to get workflow result: %w", err)
	}
	return nil
}
package temporal

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"temporal-proj/pkg/utils"
	"temporal-proj/workflowmgmt"
)

var _ workflowmgmt.WorkflowClient = (*Client)(nil)

type Client struct {
	c            client.Client
	workflowName string
	taskQueue    string
}

func NewClient(workflowName, taskQueue string) (workflowmgmt.WorkflowClient, error) {
	hostPort := utils.GetEnvOrDefault("TEMPORAL_ADDRESS", "localhost:7233")
	namespace := utils.GetEnvOrDefault("TEMPORAL_NAMESPACE", "default")

	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Temporal: %w", err)
	}

	log.Printf("Connected to Temporal server at %s (namespace: %s)", hostPort, namespace)
	return &Client{c: c, workflowName: workflowName, taskQueue: taskQueue}, nil
}

func (tc *Client) Close() {
	tc.c.Close()
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

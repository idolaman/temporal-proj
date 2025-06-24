package temporal

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"temporal-proj/pkg/utils"
)

type Client struct {
	c            client.Client
	workflowName string
	taskQueue    string
}

func NewClient(workflowName, taskQueue string) (*Client, error) {
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

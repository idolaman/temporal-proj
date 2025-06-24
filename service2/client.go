package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"nomaproj/pkg/models"
	"nomaproj/pkg/utils"
)

type TemporalClient struct {
	client client.Client
}

func NewTemporalClient() (*TemporalClient, error) {
	hostPort := utils.GetEnvOrDefault("TEMPORAL_ADDRESS", "localhost:7233")
	namespace := utils.GetEnvOrDefault("TEMPORAL_NAMESPACE", "default")

	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Temporal: %w", err)
	}

	log.Println("Connected to Temporal server")
	return &TemporalClient{client: c}, nil
}

func (tc *TemporalClient) StartScan(ctx context.Context, url string, scanID uint) error {
	task := models.ScanTask{
		URL:       url,
		Timestamp: time.Now(),
	}

	workflowID := fmt.Sprintf("scan-workflow-%d", scanID)

	_, err := tc.client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "url-scanner-task-queue",
	}, "ScanURLWorkflow", task)

	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}

	log.Printf("Started scan workflow for URL: %s (ID: %d)", url, scanID)
	return nil
}

func (tc *TemporalClient) GetScanResult(ctx context.Context, scanID uint) (*models.ScanResult, error) {
	workflowID := fmt.Sprintf("scan-workflow-%d", scanID)
	workflowRun := tc.client.GetWorkflow(ctx, workflowID, "")

	var result models.ScanResult
	err := workflowRun.Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow result: %w", err)
	}

	return &result, nil
}

func (tc *TemporalClient) Close() {
	tc.client.Close()
}

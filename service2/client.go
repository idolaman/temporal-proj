package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

// Service #1 models (copied for simplicity)
type ScanTask struct {
	URL       string    `json:"url"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ScanResult struct {
	SourceURL   string    `json:"source_url"`
	Links       []string  `json:"links"`
	TotalLinks  int       `json:"total_links"`
	ProcessedAt time.Time `json:"processed_at"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
	RequestID   string    `json:"request_id"`
}

type TemporalClient struct {
	client client.Client
}

func NewTemporalClient() (*TemporalClient, error) {
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "default",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Temporal: %w", err)
	}

	log.Println("Connected to Temporal server")
	return &TemporalClient{client: c}, nil
}

func (tc *TemporalClient) StartScan(ctx context.Context, url string, scanID uint) error {
	task := ScanTask{
		URL:       url,
		RequestID: fmt.Sprintf("scan-%d", scanID),
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

func (tc *TemporalClient) GetScanResult(ctx context.Context, scanID uint) (*ScanResult, error) {
	workflowID := fmt.Sprintf("scan-workflow-%d", scanID)
	workflowRun := tc.client.GetWorkflow(ctx, workflowID, "")

	var result ScanResult
	err := workflowRun.Get(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow result: %w", err)
	}

	return &result, nil
}

func (tc *TemporalClient) Close() {
	tc.client.Close()
}

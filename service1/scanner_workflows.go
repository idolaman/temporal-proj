package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ScannerWorkflow defines the workflow for URL scanning
type ScannerWorkflow struct{}

// ScanURLWorkflow is the main workflow function for scanning URLs
func (w *ScannerWorkflow) ScanURLWorkflow(ctx workflow.Context, task ScanTask) (ScanResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting ScanURL workflow", "url", task.URL, "request_id", task.RequestID)

	// Configure activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second, // Max time for activity to complete
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
			MaximumAttempts:    3,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Execute the ScanURL activity
	var result ScanResult
	err := workflow.ExecuteActivity(ctx, "ScanURL", task).Get(ctx, &result)
	if err != nil {
		logger.Error("ScanURL activity failed", "error", err, "request_id", task.RequestID)
		return ScanResult{
			SourceURL:   task.URL,
			RequestID:   task.RequestID,
			ProcessedAt: time.Now(),
			Success:     false,
			Error:       err.Error(),
		}, err
	}

	logger.Info("ScanURL workflow completed successfully",
		"url", task.URL,
		"links_found", result.TotalLinks,
		"request_id", task.RequestID)

	return result, nil
} 
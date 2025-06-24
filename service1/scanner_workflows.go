package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"nomaproj/pkg/models"
)

type ScannerWorkflow struct{}

// ScanURLWorkflow is the main workflow function for scanning URLs
func (w *ScannerWorkflow) ScanURLWorkflow(ctx workflow.Context, task models.ScanTask) (models.ScanResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting ScanURL workflow", "url", task.URL)

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

	var result models.ScanResult
	err := workflow.ExecuteActivity(ctx, "ScanURL", task).Get(ctx, &result)
	if err != nil {
		logger.Error("ScanURL activity failed", "error", err)
		return models.ScanResult{
			SourceURL:   task.URL,
			ProcessedAt: time.Now(),
			Success:     false,
			Error:       err.Error(),
		}, err
	}

	logger.Info("ScanURL workflow completed successfully",
		"url", task.URL,
		"links_found", result.TotalLinks)

	return result, nil
}

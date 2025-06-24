package activities

import (
	"context"
	"net/http"
	"time"

	"go.temporal.io/sdk/activity"

	"temporal-proj/pkg/models"
	"temporal-proj/service"
)

type ScannerActivities struct {
	client *http.Client
}

func NewScannerActivities() *ScannerActivities {
	return &ScannerActivities{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *ScannerActivities) ScanURL(ctx context.Context, task models.ScanTask) (models.ScanResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting URL scan", "url", task.URL)

	return service.ScanURL(ctx, s.client, task)
}

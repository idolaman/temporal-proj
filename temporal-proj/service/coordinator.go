package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"temporal-proj/temporal"
	"temporal-proj/pkg/models"
	repo "temporal-proj/repository"
)

type Coordinator struct {
	repo repo.ScanRepository
	wf   WorkflowExecutor
}

func NewCoordinator(r repo.ScanRepository, wf WorkflowExecutor) *Coordinator {
	return &Coordinator{repo: r, wf: wf}
}

func (c *Coordinator) Start(ctx context.Context, url string) (*repo.Scan, error) {
	scan := repo.Scan{URL: url, Status: "pending"}
	if err := c.repo.CreateScan(ctx, &scan); err != nil {
		return nil, err
	}

	workflowID := fmt.Sprintf("scan-workflow-%d", scan.ID)
	task := models.ScanTask{URL: url, Timestamp: time.Now()}

	if err := c.wf.Start(ctx, workflowID, task); err != nil {
		c.repo.UpdateStatus(ctx, scan.ID, "failed", -1)
		return nil, err
	}

	go c.waitForResults(scan.ID, workflowID)

	return &scan, nil
}

func (c *Coordinator) waitForResults(scanID uint, workflowID string) {
	ctx := context.Background()

	var result models.ScanResult
	err := c.wf.GetResult(ctx, workflowID, &result)
	if err != nil {
		log.Printf("Coordinator: failed to get result for scan %d: %v", scanID, err)
		_ = c.repo.UpdateStatus(ctx, scanID, "failed", -1)
		return
	}

	if result.Success {
		// Convert to Link entities
		links := make([]repo.Link, 0, len(result.Links))
		for _, u := range result.Links {
			links = append(links, repo.Link{ScanID: scanID, URL: u})
		}
		if err := c.repo.AddLinks(ctx, scanID, links); err != nil {
			log.Printf("Coordinator: failed to save links for scan %d: %v", scanID, err)
			_ = c.repo.UpdateStatus(ctx, scanID, "failed", -1)
			return
		}
		_ = c.repo.UpdateStatus(ctx, scanID, "completed", len(links))
		log.Printf("Coordinator: saved %d links for scan %d", len(links), scanID)
	} else {
		_ = c.repo.UpdateStatus(ctx, scanID, "failed", -1)
		log.Printf("Coordinator: scan %d failed: %s", scanID, result.Error)
	}
}

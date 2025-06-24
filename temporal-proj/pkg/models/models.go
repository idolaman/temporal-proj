package models

import "time"

// ScanTask represents a request to scan a URL for links.
// It is shared between the API layer and Temporal workflows.
type ScanTask struct {
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

// ScanResult represents the outcome of a scan workflow.
// Produced by the workflow and consumed by the API/service layer.
type ScanResult struct {
	SourceURL   string    `json:"source_url"`
	Links       []string  `json:"links"`
	TotalLinks  int       `json:"total_links"`
	ProcessedAt time.Time `json:"processed_at"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
}

package models

import "time"

// ScanTask represents a request to scan a URL for links
// It is shared between Service #1 (worker) and Service #2 (client/API).
type ScanTask struct {
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

// ScanResult represents the result of scanning a URL
// Produced by Service #1 and consumed by Service #2.
type ScanResult struct {
	SourceURL   string    `json:"source_url"`
	Links       []string  `json:"links"`
	TotalLinks  int       `json:"total_links"`
	ProcessedAt time.Time `json:"processed_at"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
}

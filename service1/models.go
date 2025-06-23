package main

import "time"

// ScanTask represents a request to scan a URL for links
type ScanTask struct {
	URL       string            `json:"url"`
	MaxDepth  int               `json:"max_depth,omitempty"`
	Filters   map[string]string `json:"filters,omitempty"`
	RequestID string            `json:"request_id"`
	Timestamp time.Time         `json:"timestamp"`
}

// ScanResult represents the result of scanning a URL
type ScanResult struct {
	SourceURL   string    `json:"source_url"`
	Links       []string  `json:"links"`
	TotalLinks  int       `json:"total_links"`
	ProcessedAt time.Time `json:"processed_at"`
	Success     bool      `json:"success"`
	Error       string    `json:"error,omitempty"`
	RequestID   string    `json:"request_id"`
}

// WikipediaLink represents a specific Wikipedia link with metadata
type WikipediaLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Type  string `json:"type"` // "internal", "external", "category", etc.
}

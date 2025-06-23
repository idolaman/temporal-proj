package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.temporal.io/sdk/activity"
)

// ScannerActivities contains all the activities for URL scanning
type ScannerActivities struct {
	client *http.Client
}

// NewScannerActivities creates a new instance of ScannerActivities
func NewScannerActivities() *ScannerActivities {
	return &ScannerActivities{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScanURL scans a URL and extracts all Wikipedia links
func (s *ScannerActivities) ScanURL(ctx context.Context, task ScanTask) (ScanResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting URL scan", "url", task.URL, "request_id", task.RequestID)

	result := ScanResult{
		SourceURL:   task.URL,
		RequestID:   task.RequestID,
		ProcessedAt: time.Now(),
		Success:     false,
	}

	// Validate URL
	if !isValidWikipediaURL(task.URL) {
		result.Error = "Invalid Wikipedia URL"
		return result, nil
	}

	// Make HTTP request
	resp, err := s.client.Get(task.URL)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to fetch URL: %v", err)
		return result, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("HTTP error: %d", resp.StatusCode)
		return result, nil
	}

	// Parse HTML and extract links
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to parse HTML: %v", err)
		return result, nil
	}

	links := s.extractWikipediaLinks(doc, task.URL)

	result.Links = links
	result.TotalLinks = len(links)
	result.Success = true

	logger.Info("URL scan completed", "url", task.URL, "links_found", len(links))
	return result, nil
}

// extractWikipediaLinks extracts all Wikipedia links from the document
func (s *ScannerActivities) extractWikipediaLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	linkMap := make(map[string]bool) // To avoid duplicates

	// Extract all links from the content area
	doc.Find("#mw-content-text a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		// Convert relative URLs to absolute
		absoluteURL := resolveURL(baseURL, href)

		// Filter for Wikipedia links only
		if isValidWikipediaURL(absoluteURL) && !linkMap[absoluteURL] {
			links = append(links, absoluteURL)
			linkMap[absoluteURL] = true
		}
	})

	return links
}

// isValidWikipediaURL checks if the URL is a valid Wikipedia article URL
func isValidWikipediaURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if it's a Wikipedia domain
	if !strings.Contains(parsedURL.Host, "wikipedia.org") {
		return false
	}

	// Check if it's an article URL (contains /wiki/)
	if !strings.Contains(parsedURL.Path, "/wiki/") {
		return false
	}

	// Exclude certain types of pages
	excludedPrefixes := []string{
		"/wiki/File:",
		"/wiki/Category:",
		"/wiki/Template:",
		"/wiki/Help:",
		"/wiki/Special:",
		"/wiki/User:",
		"/wiki/Talk:",
	}

	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(parsedURL.Path, prefix) {
			return false
		}
	}

	return true
}

// resolveURL resolves a relative URL against a base URL
func resolveURL(baseURL, href string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	rel, err := url.Parse(href)
	if err != nil {
		return href
	}

	return base.ResolveReference(rel).String()
} 
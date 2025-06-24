package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.temporal.io/sdk/activity"

	"nomaproj/pkg/models"
)

type ScannerActivities struct {
	client *http.Client
}

func NewScannerActivities() *ScannerActivities {
	return &ScannerActivities{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScanURL scans a URL and extracts all outgoing links (same domain or external).
func (s *ScannerActivities) ScanURL(ctx context.Context, task models.ScanTask) (models.ScanResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting URL scan", "url", task.URL)

	result := models.ScanResult{
		SourceURL:   task.URL,
		ProcessedAt: time.Now(),
		Success:     false,
	}

	// Validate URL
	if !isValidURL(task.URL) {
		result.Error = "Invalid URL"
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

	links := s.extractLinks(doc, task.URL)

	result.Links = links
	result.TotalLinks = len(links)
	result.Success = true

	logger.Info("URL scan completed", "url", task.URL, "links_found", len(links))
	return result, nil
}

// extractLinks finds all absolute links referenced in the document.
func (s *ScannerActivities) extractLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	linkMap := make(map[string]bool) // To avoid duplicates

	// Extract all anchor tags with href
	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		// Convert relative URLs to absolute
		absoluteURL := resolveURL(baseURL, href)

		// Accept only HTTP/HTTPS links and avoid duplicates
		if isValidURL(absoluteURL) && !linkMap[absoluteURL] {
			links = append(links, absoluteURL)
			linkMap[absoluteURL] = true
		}
	})

	return links
}

// isValidURL does a basic sanity check for http/https URLs.
func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
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

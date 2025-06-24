package service

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"temporal-proj/pkg/models"
	utils "temporal-proj/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

func ScanURL(ctx context.Context, client *http.Client, task models.ScanTask) (models.ScanResult, error) {
	result := models.ScanResult{
		SourceURL:   task.URL,
		ProcessedAt: time.Now(),
		Success:     false,
	}

	// Fetch and parse HTML document
	doc, err := utils.FetchDocument(ctx, client, task.URL)
	if err != nil {
		result.Error = err.Error()
		return result, nil // treat fetch errors as logical failures, not activity errors
	}

	// Extract links limited to same domain
	links := extractLinks(doc, task.URL)

	result.Links = links
	result.TotalLinks = len(links)
	result.Success = true
	return result, nil
}

func extractLinks(doc *goquery.Document, baseURL string) []string {
	var links []string
	linkMap := make(map[string]bool)

	// Determine base host once
	baseParsed, err := url.Parse(baseURL)
	var baseHost string
	if err == nil {
		baseHost = baseParsed.Host
	}

	doc.Find("a[href]").Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		abs := utils.ResolveURL(baseURL, href)
		if !utils.IsValidURL(abs) {
			return
		}

		parsed, err := url.Parse(abs)
		if err != nil || parsed.Host != baseHost {
			return
		}

		if !linkMap[abs] {
			links = append(links, abs)
			linkMap[abs] = true
		}
	})

	return links
}

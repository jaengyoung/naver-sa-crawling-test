package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// DesktopSearchURL is the Naver search URL for desktop
	DesktopSearchURL = "https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=%s"

	// Desktop CSS selectors
	desktopResultSelector = "div.nad_area ul.lst_type > li"
	desktopSiteSelector   = "a.site"
	desktopURLSelector    = "span.lnk_url_area > a.lnk_url"
	desktopTitleSelector  = "a.lnk_head span.lnk_tit"
	desktopDescSelector   = "a.link_desc"
)

// ScrapeDesktopResults scrapes Naver search results from desktop version
// This function is exported (starts with uppercase) for external use
func ScrapeDesktopResults(keyword string) ([]SearchResult, error) {
	return scrapeDesktopPage(keyword, false)
}

// ScrapeDesktopResultsWithDebug scrapes Naver search results from desktop version with debug output
// This function is exported for testing purposes
func ScrapeDesktopResultsWithDebug(keyword string) ([]SearchResult, error) {
	return scrapeDesktopPage(keyword, true)
}

// scrapeDesktopPage performs the actual scraping logic
// This function is internal (starts with lowercase)
func scrapeDesktopPage(keyword string, debug bool) ([]SearchResult, error) {
	// Build request URL
	encodedKeyword := url.QueryEscape(keyword)
	targetURL := fmt.Sprintf(DesktopSearchURL, encodedKeyword)

	if debug {
		fmt.Printf("Target URL: %s\n", targetURL)
	}

	// Create HTTP request with random headers
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add randomized headers to avoid detection
	headers := GenerateRandomHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := DesktopHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	if debug {
		fmt.Printf("Response Status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract search results
	return extractDesktopResults(doc, keyword)
}

// extractDesktopResults parses the HTML document and extracts search results
func extractDesktopResults(doc *goquery.Document, keyword string) ([]SearchResult, error) {
	var results []SearchResult

	// Find all search result items
	doc.Find(desktopResultSelector).Each(func(i int, s *goquery.Selection) {
		// Extract site name
		siteName := strings.TrimSpace(s.Find(desktopSiteSelector).Text())

		// Extract display URL
		displayURL := strings.TrimSpace(s.Find(desktopURLSelector).Text())
		displayURL = sanitizeURL(displayURL)

		// Extract title (may have multiple parts)
		var titleParts []string
		s.Find(desktopTitleSelector).Each(func(j int, titleElement *goquery.Selection) {
			titleText := strings.TrimSpace(titleElement.Text())
			if titleText != "" {
				titleParts = append(titleParts, titleText)
			}
		})
		title := strings.Join(titleParts, " · ")

		// Extract description
		description := strings.TrimSpace(s.Find(desktopDescSelector).Text())

		// Set default values for empty fields
		siteName = setDefaultValueIfEmpty(siteName, DefaultSiteName)
		displayURL = setDefaultValueIfEmpty(displayURL, DefaultURL)
		title = setDefaultValueIfEmpty(title, DefaultTitle)
		description = setDefaultValueIfEmpty(description, DefaultDescription)

		// Create search result
		result := SearchResult{
			Query:       keyword,
			Device:      DeviceDesktop,
			Rank:        i + 1,
			SiteName:    siteName,
			DisplayURL:  displayURL,
			Title:       title,
			Description: description,
		}

		results = append(results, result)
	})

	return results, nil
}

package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// MobileSearchURL is the Naver search URL for mobile
	MobileSearchURL = "https://m.search.naver.com/search.naver?sm=mtp_hty.top&where=m&query=%s"

	// Mobile CSS selectors
	mobileResultSelector = "div.api_subject_bx ul.lst_total > li"
	mobileSiteSelector   = "span.site"
	mobileURLSelector    = "span.url"
	mobileTitleSelector  = "div.tit_area span.tit"
	mobileDescSelector   = "a.desc"
)

// ScrapeMobileResults scrapes Naver search results from mobile version
// This function is exported (starts with uppercase) for external use
func ScrapeMobileResults(keyword string) ([]SearchResult, error) {
	return scrapeMobilePage(keyword, false)
}

// ScrapeMobileResultsWithDebug scrapes Naver search results from mobile version with debug output
// This function is exported for testing purposes
func ScrapeMobileResultsWithDebug(keyword string) ([]SearchResult, error) {
	return scrapeMobilePage(keyword, true)
}

// scrapeMobilePage performs the actual scraping logic
// This function is internal (starts with lowercase)
func scrapeMobilePage(keyword string, debug bool) ([]SearchResult, error) {
	// Build request URL
	encodedKeyword := url.QueryEscape(keyword)
	targetURL := fmt.Sprintf(MobileSearchURL, encodedKeyword)

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
	resp, err := MobileHTTPClient.Do(req)
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
	return extractMobileResults(doc, keyword)
}

// extractMobileResults parses the HTML document and extracts search results
func extractMobileResults(doc *goquery.Document, keyword string) ([]SearchResult, error) {
	var results []SearchResult

	// Find all search result items
	doc.Find(mobileResultSelector).Each(func(i int, s *goquery.Selection) {
		// Extract site name
		siteName := strings.TrimSpace(s.Find(mobileSiteSelector).Text())

		// Extract display URL
		displayURL := strings.TrimSpace(s.Find(mobileURLSelector).Text())
		displayURL = sanitizeURL(displayURL)

		// Extract title (may have multiple parts)
		var titleParts []string
		s.Find(mobileTitleSelector).Each(func(j int, titleElement *goquery.Selection) {
			titleText := strings.TrimSpace(titleElement.Text())
			if titleText != "" {
				titleParts = append(titleParts, titleText)
			}
		})
		title := strings.Join(titleParts, " · ")

		// Extract description
		description := strings.TrimSpace(s.Find(mobileDescSelector).Text())

		// Set default values for empty fields
		siteName = setDefaultValueIfEmpty(siteName, DefaultSiteName)
		displayURL = setDefaultValueIfEmpty(displayURL, DefaultURL)
		title = setDefaultValueIfEmpty(title, DefaultTitle)
		description = setDefaultValueIfEmpty(description, DefaultDescription)

		// Create search result
		result := SearchResult{
			Query:       keyword,
			Device:      DeviceMobile,
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

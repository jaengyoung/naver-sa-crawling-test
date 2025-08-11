# Naver Search Crawler

A Go-based web scraper for extracting search results from Naver's search engine, supporting both desktop and mobile versions.

## Architecture Overview

The project is structured using clean architecture principles with the following organization:

```
source/lambda/internal/
├── types.go           # Data structures and constants
├── http_client.go     # HTTP client configuration and headers
├── utils.go          # Utility functions
├── desktop_scraper.go # Desktop version scraping logic
├── mobile_scraper.go  # Mobile version scraping logic
└── (other files...)   # Additional functionality
```

## Core Components

### 1. Data Structures (`types.go`)

**SearchRequest**: Input structure for search operations
```go
type SearchRequest struct {
    Keyword string `json:"keyword"`
    Device  string `json:"device,omitempty"`
}
```

**SearchResult**: Output structure representing a single search result
```go
type SearchResult struct {
    Query       string `json:"query"`
    Device      string `json:"device"`  // "PC" or "Mobile"
    Rank        int    `json:"rank"`
    SiteName    string `json:"site_name"`
    DisplayURL  string `json:"display_url"`
    Title       string `json:"title"`
    Description string `json:"description"`
}
```

### 2. HTTP Client (`http_client.go`)

- **Random User Agents**: Rotates through different browser user agents
- **Random Headers**: Generates randomized HTTP headers to avoid detection
- **Connection Pooling**: Optimized HTTP clients for desktop and mobile requests

### 3. Scrapers

#### Desktop Scraper (`desktop_scraper.go`)
- **Target URL**: `https://search.naver.com/search.naver`
- **Selectors**:
  - Results container: `div.nad_area ul.lst_type > li`
  - Title: `a.lnk_head span.lnk_tit`
  - Site: `a.site`
  - URL: `span.lnk_url_area > a.lnk_url`
  - Description: `a.link_desc`

#### Mobile Scraper (`mobile_scraper.go`)
- **Target URL**: `https://m.search.naver.com/search.naver`
- **Selectors**:
  - Results container: `div.api_subject_bx ul.lst_total > li`
  - Title: `div.tit_area span.tit`
  - Site: `span.site`
  - URL: `span.url`
  - Description: `a.desc`

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "lambda/internal"
)

func main() {
    keyword := "스마트폰"
    
    // Scrape desktop results
    desktopResults, err := internal.ScrapeDesktopResults(keyword)
    if err != nil {
        log.Fatal(err)
    }
    
    // Scrape mobile results
    mobileResults, err := internal.ScrapeMobileResults(keyword)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Desktop results: %d\n", len(desktopResults))
    fmt.Printf("Mobile results: %d\n", len(mobileResults))
}
```

### Testing

Run the provided test scripts:

```bash
# Test desktop scraping
go run test_desktop.go -keyword="스마트폰"

# Test mobile scraping
go run test_mobile.go -keyword="스마트폰"
```

### With Debug Output

```go
// Get results with debug information
results, err := internal.ScrapeDesktopResultsWithDebug("스마트폰")
if err != nil {
    log.Fatal(err)
}

// This will print:
// Target URL: https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=%EC%8A%A4%EB%A7%88%ED%8A%B8%ED%8F%B0
// Response Status: 200
```

## Key Features

### 1. **Anti-Detection Mechanisms**
- Randomized User-Agent strings
- Randomized HTTP headers (Referer, Accept-Language, etc.)
- Random delays and connection pooling
- Realistic browser behavior simulation

### 2. **Robust Error Handling**
- Network timeout handling
- HTTP status code validation
- HTML parsing error recovery
- Graceful degradation with default values

### 3. **Multi-Title Support**
- Handles search results with multiple title segments
- Joins title parts with middle dot (·) separator
- Example: "LG유플러스 공식온라인스토어 · 8월한정 압도적인 이벤트혜택"

### 4. **Clean Data Output**
- Structured JSON format
- Consistent field naming
- Default values for missing data
- URL sanitization (trailing slash removal)

## Function Naming Conventions

The project follows Go's naming conventions:

- **Exported Functions** (uppercase): Can be used by external packages
  - `ScrapeDesktopResults()` - Main desktop scraping function
  - `ScrapeMobileResults()` - Main mobile scraping function
  - `ScrapeDesktopResultsWithDebug()` - Desktop scraping with debug output
  - `ScrapeMobileResultsWithDebug()` - Mobile scraping with debug output

- **Internal Functions** (lowercase): Used only within the package
  - `scrapeDesktopPage()` - Internal desktop scraping logic
  - `scrapeMobilePage()` - Internal mobile scraping logic
  - `extractDesktopResults()` - HTML parsing for desktop
  - `extractMobileResults()` - HTML parsing for mobile

## Error Handling

The scrapers return detailed error information:

```go
results, err := internal.ScrapeDesktopResults("keyword")
if err != nil {
    // Possible errors:
    // - "failed to create request: ..." 
    // - "network request failed: ..."
    // - "unexpected status code: 403"
    // - "failed to parse HTML: ..."
}
```

## Configuration

### HTTP Client Settings
- **Timeout**: 5 seconds
- **Max Idle Connections**: 100
- **Max Idle Connections Per Host**: 20
- **Idle Connection Timeout**: 90 seconds

### Default Values
- Site Name: `"-"`
- Display URL: `"-"`
- Title: `"-"`
- Description: `"-"`

## Dependencies

- `github.com/PuerkitoBio/goquery` - HTML parsing and CSS selector support
- Standard Go libraries (`net/http`, `net/url`, `strings`, etc.)
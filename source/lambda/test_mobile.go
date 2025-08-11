package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"lambda/internal"
)

func main() {
	var keyword string
	flag.StringVar(&keyword, "keyword", "", "Keyword to test mobile crawling")
	flag.Parse()

	if keyword == "" {
		fmt.Println("Usage: go run test_mobile.go -keyword=\"your keyword here\"")
		fmt.Println("Example: go run test_mobile.go -keyword=\"스마트폰\"")
		os.Exit(1)
	}

	fmt.Printf("Testing mobile scraping for keyword: %s\n", keyword)
	fmt.Println(strings.Repeat("=", 50))

	// Use the refactored function with debug output
	results, err := internal.ScrapeMobileResultsWithDebug(keyword)
	if err != nil {
		log.Printf("Crawling failed: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d results:\n", len(results))
	fmt.Println(strings.Repeat("=", 50))

	// Convert results to JSON and print
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
	fmt.Printf("\nTotal results: %d\n", len(results))
}

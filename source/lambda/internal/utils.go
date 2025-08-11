package internal

import "strings"

// setDefaultValueIfEmpty returns the defaultValue if the input string is empty or whitespace-only
func setDefaultValueIfEmpty(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

// extractMultipleTexts finds all elements matching the selector and returns their text content
func extractMultipleTexts(elements interface{}, separator string) string {
	// This would be implemented with goquery Selection
	// Keeping it simple for now, actual implementation would be in scraper files
	return ""
}

// sanitizeURL removes trailing slashes from URLs
func sanitizeURL(url string) string {
	return strings.TrimSuffix(url, "/")
}

package internal

// SearchRequest represents the input for a search crawling operation
type SearchRequest struct {
	Keyword string `json:"keyword"`
	Device  string `json:"device,omitempty"`
}

// SearchResult represents a single search result from Naver
type SearchResult struct {
	Query       string `json:"query"`
	Device      string `json:"device"`
	Rank        int    `json:"rank"`
	SiteName    string `json:"site_name"`
	DisplayURL  string `json:"display_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Device types for crawling
const (
	DeviceDesktop = "PC"
	DeviceMobile  = "MO"
)

// Default values for empty fields
const (
	DefaultSiteName    = ""
	DefaultURL         = ""
	DefaultTitle       = ""
	DefaultDescription = ""
)

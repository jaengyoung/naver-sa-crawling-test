package internal

import (
	"math/rand"
	"net/http"
	"time"
)

// HTTPClient configuration
var (
	// DesktopHTTPClient is configured for desktop web scraping
	DesktopHTTPClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// MobileHTTPClient is configured for mobile web scraping
	MobileHTTPClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}
)

// Browser User Agents
var userAgents = []string{
	// Desktop Chrome
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.7123.45 Safari/537.36",

	// Desktop Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36 Edg/137.0.0.0",

	// Desktop Firefox
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:138.0) Gecko/20100101 Firefox/138.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7; rv:138.0) Gecko/20100101 Firefox/138.0",

	// Android Chrome
	"Mozilla/5.0 (Linux; Android 14; SM-S918N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",

	// iOS Safari
	"Mozilla/5.0 (iPhone; CPU iPhone OS 18_4_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1",
}

// HTTP Headers for randomization
var (
	referers = []string{
		"https://www.naver.com/",
		"https://search.naver.com/",
		"https://www.google.com/",
		"https://www.daum.net/",
		"https://news.naver.com",
		"https://m.sports.naver.com",
		"https://map.naver.com",
		"https://shopping.naver.com",
	}

	acceptLanguages = []string{
		"ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7",
		"en-US,en;q=0.9,ko;q=0.8",
		"ko;q=0.9,en;q=0.8",
	}

	accepts = []string{
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"text/html,application/xml;q=0.9,*/*;q=0.8",
		"text/html;q=0.8,application/xhtml+xml;q=0.9,image/webp,*/*;q=0.7",
	}

	origins = []string{
		"https://www.naver.com",
		"https://search.naver.com",
		"https://www.google.com",
		"https://www.daum.net/",
		"https://news.naver.com",
		"https://m.sports.naver.com",
		"https://map.naver.com",
		"https://finance.naver.com",
		"https://comic.naver.com/index",
	}

	cookies = []string{
		"NID=abc123; NNB=xyz456",
		"SID=9876543210; LANG=ko",
		"UID=11223344; PREF=lightmode",
		"NID_SES=AAABBBCCC; NID_AUT=ZZZYYYXXX",
		"NID_JKL=JKL123456; NNB=AAABBBCCC",
		"NNB=1a2b3c4d5e; NID=abcdefg1234567",
		"SID=0011223344; PREF=darkmode; LANG=en",
		"UID=55667788; THEME=default; NID=naverid987",
	}
)

// GenerateRandomHeaders creates randomized HTTP headers to avoid detection
func GenerateRandomHeaders() map[string]string {
	headers := map[string]string{
		"User-Agent":                userAgents[rand.Intn(len(userAgents))],
		"Referer":                   referers[rand.Intn(len(referers))],
		"Accept-Language":           acceptLanguages[rand.Intn(len(acceptLanguages))],
		"Accept":                    accepts[rand.Intn(len(accepts))],
		"Connection":                "keep-alive",
		"Cache-Control":             "no-cache",
		"Upgrade-Insecure-Requests": "1",
		"DNT":                       "1",
	}

	// Randomly add optional headers
	if rand.Float64() < 0.5 {
		headers["Origin"] = origins[rand.Intn(len(origins))]
	}

	if rand.Float64() < 0.4 {
		headers["X-Requested-With"] = "XMLHttpRequest"
	}

	if rand.Float64() < 0.3 {
		headers["Cookie"] = cookies[rand.Intn(len(cookies))]
	}

	return headers
}

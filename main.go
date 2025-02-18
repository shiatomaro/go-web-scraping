package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

// Struct for user input
type ScrapeRequest struct {
	AllowedDomains []string          `json:"allowed_domains"`
	TargetURL      string            `json:"target_url"`
	Selectors      map[string]string `json:"selectors"` // Key: field name, Value: CSS selector
}

// Struct to hold scraped data
type ScrapedData map[string]string // Dynamic key-value storage

var results []ScrapedData // Store scraped data

// Scrape function
func scrapeWebsite(w http.ResponseWriter, r *http.Request) {
	var request ScrapeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	c := colly.NewCollector(
		colly.AllowedDomains(request.AllowedDomains...),
	)

	results = nil // Clear old results

	// Scrape using dynamic selectors
	c.OnHTML("body", func(h *colly.HTMLElement) {
		data := make(ScrapedData)
		for field, selector := range request.Selectors {
			data[field] = h.ChildText(selector) // Extract text by default
		}
		results = append(results, data)
	})

	// Visit user-provided URL
	err = c.Visit(request.TargetURL)
	if err != nil {
		http.Error(w, "Failed to scrape the site", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/scrape", scrapeWebsite)
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Serve frontend

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

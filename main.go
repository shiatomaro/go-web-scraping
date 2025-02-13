package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

type ScrapeRequest struct {
	AllowedDomains []string `json:"allowed_domains"`
	CSSSelector    string   `json:"css_selector"`
	TargetURL      string   `json:"target_url"`
}

type ScrapedData struct {
	Text  string `json:"text"`
	Img   string `json:"img,omitempty"`
	Price string `json:"price,omitempty"`
}

func scrapedWebsite(w http.ResponseWriter, r *http.Request) {
	var request ScrapeRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(request.AllowedDomains...),
	)

	var results []ScrapedData

	c.OnHTML(request.CSSSelector, func(h *colly.HTMLElement) {
		results = append(results, ScrapedData{
			Text:  h.Text,
			Img:   h.ChildAttr("img", "src"),
			Price: h.ChildText(".price"),
		})
	})

	err = c.Visit(request.TargetURL)
	if err != nil {
		http.Error(w, "Failed to scrape site", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "applicatiob/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/scrape", scrapedWebsite)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server running on http://localhost:8080, press ctrl+c to exist")
	http.ListenAndServe(":8080", nil)
}

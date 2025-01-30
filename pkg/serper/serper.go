package serper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// SearchResult represents the API response structure
type SearchResult struct {
	SearchParameters struct {
		Query string `json:"q"`
	} `json:"searchParameters"`
	Organic []struct {
		Title           string `json:"title"`
		Link            string `json:"link"`
		PublicationInfo string `json:"publicationInfo"`
		Snippet         string `json:"snippet"`
		Year            int    `json:"year"`
		CitedBy         int    `json:"citedBy"`
		PDFUrl          string `json:"pdfUrl,omitempty"`
		ID              string `json:"id"`
	} `json:"organic"`
	Credits   int    `json:"credits"`
	ResultStr string `json:"-"` // Holds the full JSON response as a string, but ignored in JSON marshaling
}

func VerifyClaim(claim string) (*SearchResult, error) {

	var apiKey = os.Getenv("SERPER_API_KEY")
	var apiURL = os.Getenv("SERPER_SCHOLAR_SEARCH_URL")

	// Marshal claim into JSON payload
	payload, err := json.Marshal(map[string]string{"q": claim})
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-KEY", apiKey)
	req.Header.Add("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Store full response as a string
	result.ResultStr = string(body)

	return &result, nil
}

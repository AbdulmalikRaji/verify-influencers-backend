package podchaser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GraphQLRequest represents the GraphQL request payload
type GraphQLRequest struct {
	Query string `json:"query"`
}

// PodcastResponse represents the API response structure
type PodcastResponse struct {
	Data struct {
		Podcasts struct {
			PaginatorInfo struct {
				CurrentPage  int  `json:"currentPage"`
				HasMorePages bool `json:"hasMorePages"`
				LastPage     int  `json:"lastPage"`
			} `json:"paginatorInfo"`
			Data []struct {
				ID          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				WebURL      string `json:"webUrl"`
				Episodes    struct {
					Data []struct {
						ID          string `json:"id"`
						Title       string `json:"title"`
						Description string `json:"description"`
						AirDate     string `json:"airDate"`
						AudioURL    string `json:"audioUrl"`
						Transcripts []struct {
							URL            string `json:"url"`
							Source         string `json:"source"`
							TranscriptType string `json:"transcriptType"`
							GeneratedDate  string `json:"generatedDate"`
						} `json:"transcripts"`
					} `json:"data"`
				} `json:"episodes"`
			} `json:"data"`
		} `json:"podcasts"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// FindPodcasts searches for podcasts using an influencer's name
func FindPodcasts(influencerName string) ([]PodcastResponse, error) {
	url := "https://api.podchaser.com/graphql"
	apiKey := os.Getenv("POD_DEV_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing POD_DEV_KEY environment variable")
	}

	query := fmt.Sprintf(`{
		podcasts(searchTerm: "%s") {
			paginatorInfo { currentPage, hasMorePages, lastPage }
			data {
				id, title, description, webUrl
				episodes {
					data {
						id, title, description, airDate, audioUrl
						transcripts { url, source, transcriptType, generatedDate }
					}
				}
			}
		}
	}`, influencerName)

	reqBody, _ := json.Marshal(GraphQLRequest{Query: query})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response PodcastResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0].Message)
	}

	return []PodcastResponse{response}, nil
}

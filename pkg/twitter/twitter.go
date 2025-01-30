package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Tweet struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type TwitterUserResponse struct {
	Data TwitterUser `json:"data"`
}

type TwitterUser struct {
	Name              string            `json:"name"`
	Username          string            `json:"username"`
	URL               string            `json:"url"`
	ProfileImage      string            `json:"profile_image_url"`
	Description       string            `json:"description"`
	UserPublicMetrics UserPublicMetrics `json:"public_metrics"`
}

type UserPublicMetrics struct {
	Followers int `json:"followers_count"`
}

func GetTwitterClaimsV2(username string, startTime, endTime string) ([]Tweet, error) {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		return nil, fmt.Errorf("TWITTER_BEARER_TOKEN not set")
	}

	url := fmt.Sprintf("https://api.x.com/2/tweets/search/recent?query=from:%s&start_time=%s&end_time=%s&tweet.fields=created_at", username, startTime, endTime)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweets: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch tweets: %s", resp.Status)
	}

	var result struct {
		Data []Tweet `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetTwitterUserByUsername fetches details of a Twitter user by username using Bearer Token
func GetTwitterUserByUsername(username string) (*TwitterUser, error) {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		return nil, fmt.Errorf("TWITTER_BEARER_TOKEN environment variable is not set")
	}

	// Twitter API endpoint for user lookup
	url := fmt.Sprintf("https://api.x.com/2/users/by/username/%s?user.fields=public_metrics,url,description,profile_image_url", username)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter API returned status %d", resp.StatusCode)
	}

	var response TwitterUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Return the nested user data
	return &response.Data, nil
}

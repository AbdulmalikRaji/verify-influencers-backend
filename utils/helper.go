package utils

import (
	"fmt"
	"strings"
	"time"
)

// ParseTweetTime parses a tweet's timestamp string to time.Time
func ParseTweetTime(tweetTimeStr string) (time.Time, error) {
	// Define the format for the tweet timestamp
	format := time.RFC3339Nano

	// Parse the string into time.Time
	tweetTime, err := time.Parse(format, tweetTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing tweet timestamp: %v", err)
	}

	return tweetTime, nil
}


func ConvertTimeToXFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}

func NormalizeClaim(claim string) string {
	return strings.ToLower(strings.TrimSpace(claim))
}

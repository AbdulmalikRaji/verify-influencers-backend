package dto

import "time"

type Claim struct {
	Claim     string    `json:"claim"`
	Topic     string    `json:"topic"`
	Source    int       `json:"source"`
	Raw       string    `json:"raw"`
	Category  string    `json:"category"`
	Timestamp time.Time `json:"timestamp"`
}

type FindInfluencerClaimsRequest struct {
	Username  string    `json:"username"`
	Source    int       `json:"source"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type FindInfluencerClaimsResponse struct {
	Claims   []Claim `json:"claims"`
	Username string  `json:"username"`
}

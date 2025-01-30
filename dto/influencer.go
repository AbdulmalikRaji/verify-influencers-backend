package dto

import "time"

type GetInfluencerRequest struct {
	ID       *int    `json:"id"`
	Username string `json:"username"`
}

type GetInfluencerResponse struct {
	Name       string            `json:"name"`
	Username   string            `json:"username"`
	Followers  int               `json:"followers"`
	TrustScore float64           `json:"trust_score"`
	URL        string            `json:"url"`
	Bio        string            `json:"bio"`
	Claims     []InfluencerClaim `json:"claims"`
	Topics     []string          `json:"topics"`
}

type InfluencerClaim struct {
	Claim       string    `json:"claim"`
	ClaimURL    string    `json:"claim_url"`
	Proof       string    `json:"proof"`
	ProofSource string    `json:"proof_source"`
	ProofURL    string    `json:"proof_url"`
	Topic       string    `json:"topic"`
	Score       float64   `json:"score"`
	Status      string    `json:"status"`
	ClaimedAt   time.Time `json:"claimed_at"`
}

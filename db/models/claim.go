package models

import "time"

type Claim struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	InfluencerID uint      `gorm:"not null" json:"influencer_id"`
	Content      string    `gorm:"type:varchar(500);not null" json:"content"`      // Raw content from the API (e.g., tweet, podcast)
	ParsedClaim  string    `gorm:"type:varchar(500);not null" json:"parsed_claim"` // Parsed claim after processing
	ClaimedAt    time.Time `gorm:"not null" json:"claimed_at"`                     // Timestamp of when the claim was made
	Source       string    `gorm:"type:varchar(50)" json:"source"`                 // Podcast or Tweet
	SourceURL    string    `gorm:"type:varchar(200)" json:"source_url"`            // URL of the source

	// Abstract fields
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	LastModifiedBy string    `gorm:"column:last_modified_by" json:"last_modified_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	LastModifiedAt time.Time `gorm:"column:last_modified_at" json:"last_modified_at"`
	DelFlg         bool      `gorm:"column:del_flg" json:"del_flg"` // Flag to mark deletion
}

func (Claim) TableName() string {
	return "verify_influencers.claim"
}

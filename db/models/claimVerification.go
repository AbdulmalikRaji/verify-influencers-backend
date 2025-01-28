package models

import "time"

type ClaimVerification struct {
	ID         uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	ClaimID    uint   `gorm:"not null" json:"claim_id"`
	VerifiedBy string `gorm:"type:varchar(50);not null" json:"verified_by"` // e.g., AI, Expert
	Status     string `gorm:"type:varchar(15);not null" json:"status"`      // Verified, Questionable, Debunked
	Evidence   string `gorm:"type:varchar(500)" json:"evidence"`            // Supporting evidence
	Comment    string `gorm:"type:varchar(500)" json:"comment"`             // Additional notes

	// Abstract fields
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	LastModifiedBy string    `gorm:"column:last_modified_by" json:"last_modified_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	LastModifiedAt time.Time `gorm:"column:last_modified_at" json:"last_modified_at"`
	DelFlg         bool      `gorm:"column:del_flg" json:"del_flg"`
}

func (ClaimVerification) TableName() string {
	return "verify_influencers.claim_verification"
}

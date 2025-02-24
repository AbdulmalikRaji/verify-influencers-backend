package models

import "time"

type Influencer struct {
	ID         int     `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name       string  `gorm:"type:varchar(100);not null" json:"name"`
	Username   string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Platform   string  `gorm:"type:varchar(50);not null" json:"platform"` // e.g., Twitter, Podcast
	Followers  int     `gorm:"not null" json:"followers"`
	TrustScore float64 `gorm:"type:float" json:"trust_score"`
	Category   string  `gorm:"type:varchar(100)" json:"category"`
	URL        string  `gorm:"not null" json:"url"`
	Bio        string  `gorm:"not null" json:"bio"`
	ImageURL   string  `gorm:"type:varchar(500)" json:"image_url"`

	// Abstract fields
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	LastModifiedBy string    `gorm:"column:last_modified_by" json:"last_modified_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	LastModifiedAt time.Time `gorm:"column:last_modified_at" json:"last_modified_at"`
	DelFlg         bool      `gorm:"column:del_flg" json:"del_flg"`
}

func (Influencer) TableName() string {
	return "verify_influencers.influencer"
}

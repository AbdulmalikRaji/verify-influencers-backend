package models

import "time"

type Topic struct {
	ID          int      `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`          // Name of the topic (e.g., "sleep", "hormones")
	Description string    `gorm:"type:varchar(255)" json:"description"`  // Optional description of the topic
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`   // User who created the topic
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`   // Timestamp when the topic was created
	DelFlg      bool      `gorm:"column:del_flg" json:"del_flg"`        // Flag to mark deletion
}

func (Topic) TableName() string {
	return "verify_influencers.topic"
}

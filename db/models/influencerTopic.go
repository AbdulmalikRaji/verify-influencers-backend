package models

import "time"

type InfluencerTopic struct {
	ID             int       `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	InfluencerID   int       `gorm:"not null" json:"influencer_id"` // Foreign key to Influencer
	TopicID        int       `gorm:"not null" json:"topic_id"`      // Foreign key to Topic
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	LastModifiedBy string    `gorm:"column:last_modified_by" json:"last_modified_by"`
	LastModifiedAt time.Time `gorm:"column:last_modified_at" json:"last_modified_at"`
	DelFlg         bool      `gorm:"column:del_flg" json:"del_flg"` // Flag to mark deletion
}

func (InfluencerTopic) TableName() string {
	return "verify_influencers.influencer_topic"
}

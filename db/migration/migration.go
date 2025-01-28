package migration

import (
	"log"
	"sync"

	"github.com/abdulmalikraji/verify-influencers-backend/db/models"
	"gorm.io/gorm"
)

var onlyOnce sync.Once

func Migrate(connection *gorm.DB) {

	onlyOnce.Do(func() {

		log.Println("Migrating the database...")

		if err := connection.AutoMigrate(
			&models.Influencer{},
			&models.Claim{},
			&models.ClaimVerification{},
			&models.Topic{},
			&models.InfluencerTopic{},
		); err != nil {
			log.Fatalf("Could not migrate: %v", err)
		}

		log.Println("Database migration completed successfully.")
	})
}

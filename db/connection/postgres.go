package connection

import (
	"fmt"
	"os"

	"github.com/abdulmalikraji/verify-influencers-backend/db/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initializePostgres() *gorm.DB {

	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s application_name='%s' search_path=%s sslmode=%s timezone=Europe/Istanbul",
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_DB_USER"),
		os.Getenv("POSTGRES_DB_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_DB_PORT"),
		os.Getenv("POSTGRES_DB_APP_NAME"),
		os.Getenv("POSTGRES_DB_APP_NAME"),
		os.Getenv("POSTGRES_DB_SSL_MODE"),
	)

	connection, err := gorm.Open(postgres.Open(dns))
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return nil
	}

	// Migrate the schemas for the DB
	migration.Migrate(connection)

	return connection
}

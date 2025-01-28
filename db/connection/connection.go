package connection

import (
	"sync"

	"gorm.io/gorm"
)

var once sync.Once

type Client struct {
	PostgresConnection *gorm.DB
}

var client Client

func New() Client {
	once.Do(func() {
		client = Client{
			PostgresConnection: initializePostgres(),
		}
	})

	return client
}

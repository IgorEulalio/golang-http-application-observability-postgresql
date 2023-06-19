package database

import (
	"fmt"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	if config.Config == nil {
		config.LoadConfig()
	}
}

func ConnectToDatabase() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s",
		config.Config.DatabaseUser,
		config.Config.DatabaseName,
		config.Config.DatabasePassword,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	return db, err
}

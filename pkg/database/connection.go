package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDatabase() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	return db, err
}

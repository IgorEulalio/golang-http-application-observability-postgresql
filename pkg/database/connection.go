package database

import (
	"fmt"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // or _ "github.com/jackc/pgx/v4/stdlib" for pgx
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func ConnectToDatabase() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s",
		config.Config.DatabaseUser,
		config.Config.DatabaseName,
		config.Config.DatabasePassword,
	)

	db, err := otelsql.Open("postgres", psqlInfo,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName(config.Config.DatabaseName),
	)

	if err != nil {
		return nil, err
	}

	// Wrap the *sql.DB with sqlx
	dbx := sqlx.NewDb(db, "postgres")

	return dbx, nil
}

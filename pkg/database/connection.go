package database

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // or _ "github.com/jackc/pgx/v4/stdlib" for pgx
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func ConnectToDatabase() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		config.Config.DatabaseHost,
		config.Config.DatabaseUser,
		config.Config.DatabaseName,
		config.Config.DatabasePassword,
	)

	db, err := otelsql.Open("postgres", psqlInfo,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName(config.Config.DatabaseName),
	)
	if err != nil {
		logger.Log.Error("Failed to connect to database: %v", err)
		return nil, err
	}

	// Wrap the *sql.DB with sqlx
	dbx := sqlx.NewDb(db, "postgres")

	// Run a simple query to validate the connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = dbx.PingContext(ctx)
	if err != nil {
		logger.Log.Error("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return dbx, nil
}

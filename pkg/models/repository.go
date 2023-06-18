package models

import (
	"database/sql"
	"time"
)

type Repository struct {
	ID              string         `db:"id" validate:"-"`
	Name            string         `db:"name" validate:"required"`
	Owner           string         `db:"owner" validate:"required"`
	CreationDate    time.Time      `db:"creationdate" validate:"-"`
	ConfigurationID sql.NullString `db:"configurationid" validate:"required"`
}

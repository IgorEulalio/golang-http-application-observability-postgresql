package models

import (
	"time"
)

type Repository struct {
	ID              string    `db:"id" validate:"-"`
	Name            string    `db:"name" validate:"required"`
	Owner           string    `db:"owner" validate:"required"`
	CreationDate    time.Time `db:"creationdate" validate:"-"`
	ConfigurationID string    `db:"configurationid" validate:"required"`
}

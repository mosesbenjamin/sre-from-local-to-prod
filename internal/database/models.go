// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

package repository

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID
	Result    string
	Status    string
	CreatedAT time.Time
	UpdatedAT time.Time
}

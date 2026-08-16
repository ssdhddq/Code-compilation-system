package repository

import (
	"github.com/google/uuid"
)

type Object interface {
	Get(uuid.UUID) (*Task, error)
	Save(uuid.UUID, *Task) error
	Create(*Task) error
	Delete(uuid.UUID) error
	UpdateStatus(uuid.UUID, string) error
	UpdateResult(uuid.UUID, string) error
}

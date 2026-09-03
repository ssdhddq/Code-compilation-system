package repository

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID         uuid.UUID
	Result     string
	Status     string
	Translator string
	Code       string
}

// test
type User struct {
	Id       uuid.UUID
	Login    string
	Password string
}

type Session struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	TTL       time.Time
}

package repository

import (
	"Code-compilation-system/repository/rabbit_mq"

	"github.com/google/uuid"
)

type Repository interface {
	ObjectUser
	ObjectTask
}

type ObjectTask interface {
	GetTask(uuid.UUID) (*Task, error)
	SaveTask(uuid.UUID, *Task) error
	CreateTask(*Task) error
	DeleteTask(uuid.UUID) error
	UpdateStatus(uuid.UUID, string) error
	UpdateResult(uuid.UUID, string) error
}

type ObjectUser interface {
	RegisterUser(*User) error
	AuthUser(string, string) (bool, *uuid.UUID)
	DeleteUser(uuid.UUID) error
}

type ObjectSender interface {
	Send(task rabbit_mq.TaskMessageRMQ) error
}

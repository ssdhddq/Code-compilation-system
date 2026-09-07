package postgres

import (
	"Code-compilation-system/repository"

	"github.com/google/uuid"
)

func (r *Object) GetTask(uuid uuid.UUID) (*repository.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Object) SaveTask(uuid uuid.UUID, task *repository.Task) error {
	//TODO implement me
	panic("implement me")
}

func (r *Object) CreateTask(task *repository.Task) error {
	//TODO implement me
	panic("implement me")
}

func (r *Object) DeleteTask(uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Object) UpdateStatus(uuid uuid.UUID, s string) error {
	//TODO implement me
	panic("implement me")
}

func (r *Object) UpdateResult(uuid uuid.UUID, s string) error {
	//TODO implement me
	panic("implement me")
}

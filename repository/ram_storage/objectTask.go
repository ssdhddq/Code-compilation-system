package ram_storage

import (
	"Code-compilation-system/repository"
	"sync"

	"github.com/google/uuid"
)

type ObjectTask struct {
	m    sync.RWMutex
	data map[uuid.UUID]*repository.Task
}

func NewObjectTask() *ObjectTask {
	return &ObjectTask{
		data: make(map[uuid.UUID]*repository.Task),
	}
}

func (o *ObjectTask) GetTask(id uuid.UUID) (*repository.Task, error) {
	o.m.RLock()
	defer o.m.RUnlock()
	task, exists := o.data[id]
	if !exists {
		return nil, repository.NotFound
	}
	return task, nil
}

func (o *ObjectTask) SaveTask(id uuid.UUID, task *repository.Task) error {
	o.m.Lock()
	defer o.m.Unlock()
	if task == nil {
		return repository.NilTask
	}
	o.data[id] = task
	return nil
}

func (o *ObjectTask) CreateTask(task *repository.Task) error {
	o.m.Lock()
	defer o.m.Unlock()
	if task == nil {
		return repository.NilTask
	}
	if _, exists := o.data[task.ID]; exists {
		return repository.KeyExists
	}
	o.data[task.ID] = task
	return nil
}

func (o *ObjectTask) DeleteTask(id uuid.UUID) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exists := o.data[id]; !exists {
		return repository.NotFound
	}
	delete(o.data, id)
	return nil
}

func (o *ObjectTask) UpdateStatus(id uuid.UUID, status string) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exists := o.data[id]; !exists {
		return repository.NotFound
	}
	o.data[id].Status = status

	return nil
}

func (o *ObjectTask) UpdateResult(id uuid.UUID, result string) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exists := o.data[id]; !exists {
		return repository.NotFound
	}
	o.data[id].Result = result

	return nil
}

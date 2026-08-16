package ram_storage

import (
	"Code-compilation-system/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Object struct {
	m    sync.RWMutex
	data map[uuid.UUID]*repository.Task
}

func NewObject() *Object {
	return &Object{
		data: make(map[uuid.UUID]*repository.Task),
	}
}

func (o *Object) Get(id uuid.UUID) (*repository.Task, error) {
	o.m.RLock()
	defer o.m.RUnlock()
	task, exists := o.data[id]
	if !exists {
		return nil, repository.NotFound
	}
	return task, nil
}

func (o *Object) Save(id uuid.UUID, task *repository.Task) error {
	o.m.Lock()
	defer o.m.Unlock()
	if task == nil {
		return repository.NilTask
	}
	o.data[id] = task
	return nil
}

func (o *Object) Create(task *repository.Task) error {
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

func (o *Object) Delete(id uuid.UUID) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exists := o.data[id]; !exists {
		return repository.NotFound
	}
	delete(o.data, id)
	return nil
}

func (o *Object) UpdateStatus(id uuid.UUID, status string) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exists := o.data[id]; !exists {
		return repository.NotFound
	}
	o.data[id].Status = status
	o.data[id].UpdatedAT = time.Now()
	return nil
}

func (o *Object) UpdateResult(id uuid.UUID, result string) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exists := o.data[id]; !exists {
		return repository.NotFound
	}
	o.data[id].Result = result
	o.data[id].UpdatedAT = time.Now()
	return nil
}

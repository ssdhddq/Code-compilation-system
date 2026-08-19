package ram_storage

import (
	"Code-compilation-system/repository"
	"sync"

	"github.com/google/uuid"
)

type ObjectUser struct {
	m      sync.RWMutex
	data   map[uuid.UUID]*repository.User
	logins map[string]uuid.UUID
}

func NewObjectUser() *ObjectUser {
	return &ObjectUser{
		data:   make(map[uuid.UUID]*repository.User),
		logins: make(map[string]uuid.UUID),
	}
}

func (o *ObjectUser) RegisterUser(user *repository.User) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exist := o.data[user.Id]; exist {
		return repository.KeyExists
	}
	o.data[user.Id] = user
	o.logins[user.Login] = user.Id
	return nil
}

func (o *ObjectUser) AuthUser(login, password string) bool {
	o.m.RLock()
	defer o.m.RUnlock()
	if id, exist := o.logins[login]; !exist {
		return false
	} else if password == o.data[id].Password {
		return true
	}
	return false
}

func (o *ObjectUser) DeleteUser(u uuid.UUID) error {
	o.m.Lock()
	defer o.m.Unlock()
	if user, exist := o.data[u]; exist {
		return repository.NotFound
	} else {
		delete(o.data, u)
		delete(o.logins, user.Login)
		return nil
	}
}

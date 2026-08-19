package ram_storage

import (
	"Code-compilation-system/repository"
	"log"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	log.Printf("RegisterUser: login='%s', password='%s'", user.Login, user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	o.data[user.Id] = user
	o.logins[user.Login] = user.Id
	return nil
}

func (o *ObjectUser) AuthUser(login, password string) (bool, *uuid.UUID) {
	o.m.RLock()
	defer o.m.RUnlock()
	log.Printf("AuthUser: login='%s', password='%s'", login, password)
	if id, exist := o.logins[login]; !exist {
		return false, nil
	} else if err := bcrypt.CompareHashAndPassword([]byte(o.data[id].Password), []byte(password)); err == nil {
		return true, &id
	}
	return false, nil
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

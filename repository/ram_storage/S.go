package ram_storage

import (
	"Code-compilation-system/repository"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ObjectSession struct {
	m    sync.RWMutex
	data map[uuid.UUID]*repository.Session
}

func NewObjectSession() *ObjectSession {
	return &ObjectSession{
		data: make(map[uuid.UUID]*repository.Session),
	}
}

func (o *ObjectSession) CreateSession(session *repository.Session) error {
	o.m.Lock()
	defer o.m.Unlock()
	if session == nil {
		return repository.NilSession
	}
	o.data[session.SessionID] = session
	return nil
}

func (o *ObjectSession) ValidateToken(u uuid.UUID) bool {
	o.m.RLock()
	defer o.m.RUnlock()
	if _, exist := o.data[u]; !exist {
		return false
	}
	return time.Now().Before(o.data[u].TTL)
}

func (o *ObjectSession) DeleteToken(u uuid.UUID) error {
	o.m.Lock()
	defer o.m.Unlock()
	if _, exist := o.data[u]; !exist {
		return repository.NotFound
	}
	delete(o.data, u)
	return nil
}

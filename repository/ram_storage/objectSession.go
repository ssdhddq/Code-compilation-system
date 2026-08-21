package ram_storage

import (
	"Code-compilation-system/repository"
	"container/list"
	"sync"
	"time"
)

var provider = &Provider{
	lock:     sync.Mutex{},
	sessions: make(map[string]*list.Element),
	list:     list.New(),
}

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (s *SessionStore) Set(key, value interface{}) error {
	s.value[key] = value
	err := provider.SessionUpdate(s.sid)
	return err
}

func (s *SessionStore) Get(key interface{}) interface{} {
	err := provider.SessionUpdate(s.sid)
	if err != nil {
		return nil
	}
	if v, ok := s.value[key]; ok {
		return v
	}
	return nil
}

func (s *SessionStore) Delete(key interface{}) error {
	if _, exist := s.value[key]; !exist {
		return repository.NotFound
	}
	delete(s.value, key)
	return nil
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

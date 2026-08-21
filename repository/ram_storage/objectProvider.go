package ram_storage

import (
	"Code-compilation-system/repository"
	"Code-compilation-system/session"
	"container/list"
	"sync"
	"time"
)

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func NewProvider() *Provider {
	return &Provider{
		lock:     sync.Mutex{},
		sessions: make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (provider *Provider) SessionInit(sid string) (session.Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	v := make(map[interface{}]interface{})
	newSess := &SessionStore{
		sid:          sid,
		timeAccessed: time.Now(),
		value:        v,
	}
	element := provider.list.PushBack(newSess)
	provider.sessions[sid] = element
	return newSess, nil
}

func (provider *Provider) SessionRead(sid string) (session.Session, error) {
	if v, ok := provider.sessions[sid]; ok {
		return v.Value.(*SessionStore), nil
	}
	sess, err := provider.SessionInit(sid)
	return sess, err
}

func (provider *Provider) SessionUpdate(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	if element, ok := provider.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		provider.list.MoveToFront(element)
		return nil
	}
	return repository.NotFound
}

func (provider *Provider) SessionDestroy(sid string) error {
	if element, ok := provider.sessions[sid]; ok {
		delete(provider.sessions, sid)
		provider.list.Remove(element)
		return nil
	}
	return repository.NotFound
}

func (provider *Provider) SessionGC(maxLifeTime int64) {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	for {
		element := provider.list.Back()
		if element == nil {
			return
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			provider.list.Remove(element)
			delete(provider.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}

	}
}

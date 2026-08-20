package session

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
)

var provides = make(map[string]Provider)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime int64
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("unknown provide: %s", provideName)
	}
	return &Manager{
		provider:    provider,
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
	}, nil
}

func RegisterProvider(name string, provider Provider) {
	if provider == nil {
		panic("register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("register called twice for provider " + name)
	}
	provides[name] = provider
}

func (manager *Manager) sessionId() string {
	return uuid.New().String()
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, err = manager.provider.SessionInit(sid)
		if err != nil {
			return
		}
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			return
		}
		session, err = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if err := manager.provider.SessionDestroy(cookie.Value); err != nil {
		return
	}
	expiration := time.Now()
	cookieNew := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
	http.SetCookie(w, &cookieNew)
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}

package ram_storage

import "time"

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

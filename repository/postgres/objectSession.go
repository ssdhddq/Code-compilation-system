package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[any]any
	pool         *pgxpool.Pool
}

func (s *SessionStore) Set(key, value any) error {
	s.value[key] = value
	return s.updateData()
}

func (s *SessionStore) Get(key any) any {
	if v, ok := s.value[key]; ok {
		return v
	}
	return nil
}

func (s *SessionStore) Delete(key any) error {
	delete(s.value, key)
	return s.updateData()
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

func (s *SessionStore) updateData() error {
	jsonData := make(map[string]any)
	for k, v := range s.value {
		if str, ok := k.(string); ok {
			jsonData[str] = v
		}
	}
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	query := `UPDATE sessions SET data = $1 WHERE session_id = $2`
	_, err = s.pool.Exec(context.Background(), query, bytes, s.sid)
	return err
}

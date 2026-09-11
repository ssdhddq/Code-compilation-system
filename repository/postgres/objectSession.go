package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStore struct {
	sid          string
	userID       uuid.UUID
	timeAccessed time.Time
	value        map[any]any
	pool         *pgxpool.Pool
}

func (s *SessionStore) Set(key, value interface{}) error {
	if keyStr, ok := key.(string); ok && keyStr == "userID" {
		if userIDStr, ok := value.(string); ok {
			id, err := uuid.Parse(userIDStr)
			if err != nil {
				return err
			}
			s.userID = id
			query := `UPDATE sessions SET user_id = $1 WHERE session_id = $2`
			_, err = s.pool.Exec(context.Background(), query, id, s.sid)
			return err
		}
		return fmt.Errorf("invalid userID value")
	}
	return nil
}

func (s *SessionStore) Get(key interface{}) interface{} {
	if keyStr, ok := key.(string); ok && keyStr == "userID" {
		return s.userID.String()
	}
	return nil
}

func (s *SessionStore) Delete(key interface{}) error {
	if keyStr, ok := key.(string); ok && keyStr == "userID" {
		s.userID = uuid.Nil
		query := `UPDATE sessions SET user_id = NULL WHERE session_id = $1`
		_, err := s.pool.Exec(context.Background(), query, s.sid)
		return err
	}
	return nil
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

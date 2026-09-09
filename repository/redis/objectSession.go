package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	client *redis.Client
	userID string
	sid    string
	data   map[string]any
	ttl    time.Duration
}

type sessionData struct {
	UserID    string         `json:"user_id"`
	ExpiresAt time.Time      `json:"expires_at"`
	Data      map[string]any `json:"data"`
}

func (s *SessionStore) Set(key, value any) error {
	keyStr, ok := key.(string)
	if !ok {
		return fmt.Errorf("key not string")
	}
	if keyStr == "userID" {
		userIDStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("userID not string")
		}
		s.userID = userIDStr
	} else {
		s.data[keyStr] = value
	}
	return s.saveToRedis()
}

func (s *SessionStore) Get(key any) any {
	keyStr, ok := key.(string)
	if !ok {
		return nil
	}
	return s.data[keyStr]
}

func (s *SessionStore) Delete(key any) error {
	keyStr, ok := key.(string)
	if !ok {
		return fmt.Errorf("key not string")
	}
	delete(s.data, keyStr)
	return s.saveToRedis()
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

func (s *SessionStore) saveToRedis() error {
	ctx := context.Background()
	jsonData, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, s.sid, jsonData, s.ttl).Err()
}

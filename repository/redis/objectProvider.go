package redis

import (
	"Code-compilation-system/session"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Provider struct {
	client *redis.Client
	ttl    time.Duration
}

func NewProvider(addr, password string, db int, ttl time.Duration) *Provider {
	newClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Provider{
		client: newClient,
		ttl:    ttl,
	}
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	store := &SessionStore{
		client: p.client,
		sid:    sid,
		userID: uuid.Nil.String(),
		data:   make(map[string]any),
		ttl:    p.ttl,
	}
	if err := store.saveToRedis(); err != nil {
		return nil, err
	}
	return store, nil
}

func (p *Provider) SessionRead(sid string) (session.Session, error) {
	ctx := context.Background()
	val, err := p.client.Get(ctx, sid).Result()
	if err == redis.Nil {
		return p.SessionInit(sid)
	}
	if err != nil {
		return nil, err
	}
	var data struct {
		UserID string         `json:"user_id"`
		Data   map[string]any `json:"data"`
	}

	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	log.Printf("SessionRead: sid=%s, userID='%s'", sid, data.UserID)
	store := &SessionStore{
		client: p.client,
		sid:    sid,
		userID: data.UserID,
		data:   data.Data,
		ttl:    p.ttl,
	}

	p.client.Expire(ctx, sid, p.ttl)
	return store, nil
}

func (p *Provider) SessionUpdate(sid string) error {
	ctx := context.Background()
	return p.client.Expire(ctx, sid, p.ttl).Err()
}

func (p *Provider) SessionDestroy(sid string) error {
	ctx := context.Background()
	return p.client.Del(ctx, sid).Err()
}

func (p *Provider) SessionGC(maxLifeTime int64) {
	//В редисе есть GC
}

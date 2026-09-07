package postgres

import (
	"Code-compilation-system/session"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider struct {
	pool *pgxpool.Pool
}

func NewProvider(pool *pgxpool.Pool) *Provider {
	return &Provider{pool: pool}
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	data := make(map[interface{}]interface{})
	jsonData := []byte("{}")
	query := `INSERT INTO sessions (session_id, user_id, expires_at, data) VALUES ($1, $2, $3, $4)`
	// Пока user_id не знаем, ставим нулевой, expires_at через maxLifeTime будет установлен отдельно
	// Но мы не знаем maxLifeTime здесь – он передаётся в SessionGC, а не в Init.
	// Можно сохранить время жизни по умолчанию (например, 24h) или обновлять через SessionUpdate.
	expiresAt := time.Now().Add(24 * time.Hour) // временно
	_, err := p.pool.Exec(context.Background(), query, sid, uuid.Nil, expiresAt, jsonData)
	if err != nil {
		return nil, err
	}
	store := &PostgresSessionStore{
		sid:       sid,
		data:      data,
		pool:      p.pool,
		expiresAt: expiresAt,
	}
	return store, nil
}

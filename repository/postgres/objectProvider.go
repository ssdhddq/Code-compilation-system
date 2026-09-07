package postgres

import (
	"Code-compilation-system/session"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider struct {
	pool *pgxpool.Pool
}

func NewProvider(pool *pgxpool.Pool) *Provider {
	return &Provider{pool: pool}
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	timeAccessed := time.Now().Add(24 * time.Hour)
	query := `INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := p.pool.Exec(context.Background(), query, sid, nil, timeAccessed)
	if err != nil {
		return nil, err
	}
	store := &SessionStore{
		sid:          sid,
		userID:       uuid.Nil,
		timeAccessed: timeAccessed,
		pool:         p.pool,
	}
	return store, nil
}

func (p *Provider) SessionUpdate(sid string) error {
	query := `UPDATE sessions SET expires_at = NOW() + INTERVAL '1 day' WHERE session_id = $1`
	_, err := p.pool.Exec(context.Background(), query, sid)
	return err
}

func (p *Provider) SessionRead(sid string) (session.Session, error) {
	var expiresAt time.Time
	var userID uuid.UUID
	query := `SELECT user_id, expires_at FROM sessions WHERE session_id = $1`
	err := p.pool.QueryRow(context.Background(), query, sid).Scan(&userID, &expiresAt)
	if err != nil {
		return nil, err
	}
	store := &SessionStore{
		sid:          sid,
		userID:       userID,
		timeAccessed: expiresAt,
		pool:         p.pool,
	}
	return store, nil
}

func (p *Provider) SessionDestroy(sid string) error {
	query := `DELETE FROM sessions WHERE session_id = $1`
	_, err := p.pool.Exec(context.Background(), query, sid)
	return err
}

func (p *Provider) SessionGC(maxLifeTime int64) {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, _ = p.pool.Exec(context.Background(), query)
}

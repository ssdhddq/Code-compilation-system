package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider struct {
	pool *pgxpool.Pool
}

func NewProvider(pool *pgxpool.Pool) *Provider {
	return &Provider{pool: pool}
}

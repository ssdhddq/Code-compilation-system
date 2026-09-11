package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Object struct {
	Pool *pgxpool.Pool
}

func NewObject(conn string) (*Object, error) {
	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, err
	}
	return &Object{Pool: pool}, nil
}

func (r *Object) Close() {
	r.Pool.Close()
}

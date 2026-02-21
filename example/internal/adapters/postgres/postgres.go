package postgres

import (
	"context"
	"fmt"
	"goserve/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ db.Database = &Postgres{}

type Postgres struct {
	pool *pgxpool.Pool
}

func New(connectionUrl string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), connectionUrl)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		pool: pool,
	}, nil
}

// Connect implements [db.Database].
func (p *Postgres) Connect(ctx context.Context) error {
	if _, err := p.pool.Acquire(ctx); err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	return nil
}

// Disconnect implements [db.Database].
func (p *Postgres) Disconnect(ctx context.Context) error {
	p.pool.Close()
	return nil
}

// Ping implements [db.Database].
func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

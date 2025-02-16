package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	dbc *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context, dsn string) (*pgClient, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		dbc: dbc,
	}, nil
}

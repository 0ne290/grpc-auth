package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"grpc-auth/internal/core/services"
)

type postgresUnitOfWorkStarter struct {
	pool *pgxpool.Pool
}

func NewPostgresUnitOfWorkStarter(pool *pgxpool.Pool) services.UnitOfWorkStarter {
	return &postgresUnitOfWorkStarter{pool}
}

func (uows *postgresUnitOfWorkStarter) Start(ctx context.Context) (services.UnitOfWork, error) {
	transaction, err := uows.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPostgresUnitOfWork(transaction), nil
}

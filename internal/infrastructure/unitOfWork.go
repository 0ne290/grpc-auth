package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"grpc-auth/internal/core"
)

type PostgresUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPostgresUnitOfWork(pool *pgxpool.Pool) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{pool}
}

func (uow *PostgresUnitOfWork) Begin(ctx context.Context) (core.Repository, error) {
	transaction, err := uow.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPosgresRepository(transaction), nil
}

func (*PostgresUnitOfWork) Save(ctx context.Context, repository core.Repository) error {
	postgresRepository := repository.(*PosgresRepository)

	return postgresRepository.transaction.Commit(ctx)
}

func (*PostgresUnitOfWork) Rollback(ctx context.Context, repository core.Repository) error {
	postgresRepository := repository.(*PosgresRepository)

	return postgresRepository.transaction.Rollback(ctx)
}

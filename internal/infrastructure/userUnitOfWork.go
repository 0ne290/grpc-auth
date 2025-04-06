package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"grpc-auth/internal/core/services"
)

type PostgresUserUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewPostgresUserUnitOfWork(pool *pgxpool.Pool) *PostgresUserUnitOfWork {
	return &PostgresUserUnitOfWork{pool}
}

func (uow *PostgresUserUnitOfWork) Begin(ctx context.Context) (services.UserRepository, error) {
	transaction, err := uow.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPosgresUserRepository(transaction), nil
}

func (*PostgresUserUnitOfWork) Save(ctx context.Context, repository services.UserRepository) error {
	postgresRepository := repository.(*PosgresUserRepository)

	return postgresRepository.transaction.Commit(ctx)
}

func (*PostgresUserUnitOfWork) Rollback(ctx context.Context, repository services.UserRepository) error {
	postgresRepository := repository.(*PosgresUserRepository)

	return postgresRepository.transaction.Rollback(ctx)
}

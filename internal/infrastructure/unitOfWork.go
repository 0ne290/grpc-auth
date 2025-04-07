package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v5"
	"grpc-auth/internal/core/services"
)

type postgresUnitOfWork struct {
	transaction       pgx.Tx
	userRepository    *PosgresUserRepository
	sessionRepository *PosgresSessionRepository
}

func newPostgresUnitOfWork(transaction pgx.Tx) *postgresUnitOfWork {
	return &postgresUnitOfWork{transaction, newPosgresUserRepository(transaction), newPosgresSessionRepository(transaction)}
}

func (uow *postgresUnitOfWork) UserRepository() services.UserRepository {
	return uow.userRepository
}

func (uow *postgresUnitOfWork) SessionRepository() services.SessionRepository {
	return uow.sessionRepository
}

func (uow *postgresUnitOfWork) Save(ctx context.Context) error {
	return uow.transaction.Commit(ctx)
}

func (uow *postgresUnitOfWork) Rollback(ctx context.Context) error {
	return uow.transaction.Rollback(ctx)
}

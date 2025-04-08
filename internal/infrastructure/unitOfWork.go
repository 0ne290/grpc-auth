package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
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

type MockUnitOfWork struct {
	mock.Mock
}

func NewMockUnitOfWork() *MockUnitOfWork {
	return &MockUnitOfWork{}
}

func (uow *MockUnitOfWork) UserRepository() services.UserRepository {
	args := uow.Called()
	return args.Get(0).(services.UserRepository)
}

func (uow *MockUnitOfWork) SessionRepository() services.SessionRepository {
	args := uow.Called()
	return args.Get(0).(services.SessionRepository)
}

func (uow *MockUnitOfWork) Save(ctx context.Context) error {
	args := uow.Called(ctx)
	return args.Error(0)
}

func (uow *MockUnitOfWork) Rollback(ctx context.Context) error {
	args := uow.Called(ctx)
	return args.Error(0)
}

package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
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

type MockUnitOfWorkStarter struct {
	mock.Mock
}

func NewMockUnitOfWorkStarter() *MockUnitOfWorkStarter {
	return &MockUnitOfWorkStarter{}
}

func (uows *MockUnitOfWorkStarter) Start(ctx context.Context) (services.UnitOfWork, error) {
	args := uows.Called(ctx)
	return args.Get(0).(services.UnitOfWork), args.Error(1)
}

package infrastructure

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"grpc-auth/internal/core/entities"
)

type PosgresUserRepository struct {
	transaction pgx.Tx
}

func newPosgresUserRepository(transaction pgx.Tx) *PosgresUserRepository {
	return &PosgresUserRepository{transaction}
}

func (r *PosgresUserRepository) TryCreate(ctx context.Context, user *entities.User) (bool, error) {
	const query string = "INSERT INTO users VALUES ($1, $2, $3, $4)"

	_, err := r.transaction.Exec(ctx, query, user.Uuid, user.CreatedAt, user.Name, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Check unique_violation PostgreSQL error
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *PosgresUserRepository) TryUpdate(ctx context.Context, user *entities.User) (bool, error) {
	const query string = "UPDATE users SET name = $2, password = $3 WHERE uuid = $1"

	_, err := r.transaction.Exec(ctx, query, user.Uuid, user.Name, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Check unique_violation PostgreSQL error
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *PosgresUserRepository) TryGetByName(ctx context.Context, name string) (*entities.User, error) {
	const query string = "SELECT * FROM users WHERE name = $1 FOR UPDATE"

	user := &entities.User{}

	err := r.transaction.QueryRow(ctx, query, name).Scan(&user.Uuid, &user.CreatedAt, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *PosgresUserRepository) TryGetByUuid(ctx context.Context, userUuid uuid.UUID) (*entities.User, error) {
	const query string = "SELECT * FROM users WHERE uuid = $1 FOR UPDATE"

	user := &entities.User{}

	err := r.transaction.QueryRow(ctx, query, userUuid).Scan(&user.Uuid, &user.CreatedAt, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *PosgresUserRepository) TryDelete(ctx context.Context, userUuid uuid.UUID) (bool, error) {
	const query string = "DELETE FROM users WHERE uuid = $1"

	commandTag, err := r.transaction.Exec(ctx, query, userUuid)
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() != 0, nil
}

func (r *PosgresUserRepository) Exists(ctx context.Context, userUuid uuid.UUID) (bool, error) {
	const query string = "SELECT EXISTS(SELECT 1 FROM users WHERE uuid = $1)"

	var exists bool
	err := r.transaction.QueryRow(ctx, query, userUuid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *MockUserRepository { return &MockUserRepository{} }

func (r *MockUserRepository) TryCreate(ctx context.Context, user *entities.User) (bool, error) {
	args := r.Called(ctx, user)
	return args.Bool(0), args.Error(1)
}

func (r *MockUserRepository) TryGetByName(ctx context.Context, name string) (*entities.User, error) {
	args := r.Called(ctx, name)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (r *MockUserRepository) TryDelete(ctx context.Context, userUuid uuid.UUID) (bool, error) {
	args := r.Called(ctx, userUuid)
	return args.Bool(0), args.Error(1)
}

func (r *MockUserRepository) Exists(ctx context.Context, userUuid uuid.UUID) (bool, error) {
	args := r.Called(ctx, userUuid)
	return args.Bool(0), args.Error(1)
}

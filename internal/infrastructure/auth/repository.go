package auth

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core "grpc-auth/internal/core/auth"
)

type PosgresRepository struct {
	transaction pgx.Tx
}

func newPosgresRepository(transaction pgx.Tx) *PosgresRepository {
	return &PosgresRepository{transaction}
}

func (r *PosgresRepository) TryCreate(ctx context.Context, user *core.User) (bool, error) {
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

func (r *PosgresRepository) TryGetByName(ctx context.Context, name string) (*core.User, error) {
	const query string = "SELECT * FROM users WHERE name = $1 FOR UPDATE"

	user := &core.User{}

	err := r.transaction.QueryRow(ctx, query, name).Scan(&user.Uuid, &user.CreatedAt, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

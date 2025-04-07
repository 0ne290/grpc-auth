package infrastructure

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"grpc-auth/internal/core/entities"
)

type PosgresSessionRepository struct {
	transaction pgx.Tx
}

func newPosgresSessionRepository(transaction pgx.Tx) *PosgresSessionRepository {
	return &PosgresSessionRepository{transaction}
}

func (r *PosgresSessionRepository) Create(ctx context.Context, session *entities.Session) error {
	const query string = "INSERT INTO sessions VALUES ($1, $2, $3)"

	_, err := r.transaction.Exec(ctx, query, session.RefreshToken, session.UserUuid, session.ExpirationAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PosgresSessionRepository) TryGetByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*entities.Session, error) {
	const query string = "SELECT * FROM sessions WHERE refresh_token = $1 FOR UPDATE"

	session := &entities.Session{}

	err := r.transaction.QueryRow(ctx, query, refreshToken).Scan(&session.RefreshToken, &session.UserUuid, &session.ExpirationAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}

func (r *PosgresSessionRepository) DeleteByRefreshToken(ctx context.Context, refreshToken uuid.UUID) error {
	const query string = "DELETE FROM sessions WHERE refresh_token = $1"

	_, err := r.transaction.Exec(ctx, query, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

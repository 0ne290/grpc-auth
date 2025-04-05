package auth

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (Repository, error)
	Save(ctx context.Context, repository Repository) error
	Rollback(ctx context.Context, repository Repository) error
}

type Repository interface {
	TryCreate(ctx context.Context, user *User) (bool, error)
	TryGetByName(ctx context.Context, name string) (*User, error)
}

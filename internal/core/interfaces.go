package core

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type UuidProvider interface {
	Random() uuid.UUID
}

type TimeProvider interface {
	Now() time.Time
}

type Salter interface {
	Salt(uuid uuid.UUID, createdAt time.Time, name, password string) string
}

type Hasher interface {
	Hash(saltedPassword string) string
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Repository, error)
	Save(ctx context.Context, repository Repository) error
	Rollback(ctx context.Context, repository Repository) error
}

type Repository interface {
	TryCreate(ctx context.Context, user *User) (bool, error)
	TryGetByName(ctx context.Context, name string) (*User, error)
}

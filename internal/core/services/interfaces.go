package services

import (
	"context"
	"github.com/google/uuid"
	"grpc-auth/internal/core/entities"
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

type UserUnitOfWork interface {
	Begin(ctx context.Context) (UserRepository, error)
	Save(ctx context.Context, repository UserRepository) error
	Rollback(ctx context.Context, repository UserRepository) error
}

type UserRepository interface {
	TryCreate(ctx context.Context, user *entities.User) (bool, error)
	TryGetByName(ctx context.Context, name string) (*entities.User, error)
}

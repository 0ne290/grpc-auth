package services

import (
	"context"
	"github.com/google/uuid"
	"grpc-auth/internal/core/entities"
	"grpc-auth/internal/core/valueObjects"
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

type UnitOfWorkStarter interface {
	Start(ctx context.Context) (UnitOfWork, error)
}

type UnitOfWork interface {
	UserRepository() UserRepository
	SessionRepository() SessionRepository

	Save(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UserRepository interface {
	TryCreate(ctx context.Context, user *entities.User) (bool, error)
	TryGetByName(ctx context.Context, name string) (*entities.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *entities.Session) error
	TryGetByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*entities.Session, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken uuid.UUID) error
}

type JwtManager interface {
	Generate(info *valueObjects.AuthInfo) (string, error)
	Parse(tokenString string) *valueObjects.AuthInfo
}

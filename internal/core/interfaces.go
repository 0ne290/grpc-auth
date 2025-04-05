package core

import (
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

package entities

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	RefreshToken uuid.UUID
	UserUuid     uuid.UUID
	ExpirationAt time.Time
}

func NewSession(refreshToken, userUuid uuid.UUID, expirationAt time.Time) *Session {
	return &Session{refreshToken, userUuid, expirationAt}
}

package infrastructure

import (
	"github.com/google/uuid"
	"time"
)

type RealSalter struct {
	staticSalt string
}

func NewRealSalter(staticSalt string) *RealSalter {
	return &RealSalter{staticSalt}
}

func (s *RealSalter) Salt(uuid uuid.UUID, createdAt time.Time, name, password string) string {
	uuidString := uuid.String()
	createdAtString := createdAt.String()
	return createdAtString + s.staticSalt + name + s.staticSalt + password + createdAtString + uuidString + uuidString
}

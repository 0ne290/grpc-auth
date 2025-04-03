package core

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Uuid      uuid.UUID
	CreatedAt time.Time
	Name      string
	Password  string
}

func newUser(uuid uuid.UUID, createdAt time.Time, name, password string) *User {
	return &User{uuid, createdAt, name, password}
}

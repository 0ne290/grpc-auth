package infrastructure

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
)

const staticSalt string = "-\\S?-bPGZO{n!]o6&8VvL2;oR*E7f~~pQe-;b*Z9qKkZ]HB<zLYC*PP1q>=Y^{gT"

type RealSalter struct{}

func NewRealSalter() *RealSalter {
	return &RealSalter{}
}

func (s *RealSalter) Salt(uuid uuid.UUID, createdAt time.Time, name, password string) string {
	uuidString := uuid.String()
	createdAtString := createdAt.String()
	return createdAtString + staticSalt + name + staticSalt + password + createdAtString + uuidString + uuidString
}

type MockSalter struct {
	mock.Mock
}

func NewMockSalter() *MockSalter {
	return &MockSalter{}
}

func (s *MockSalter) Salt(uuid uuid.UUID, createdAt time.Time, name, password string) string {
	args := s.Called(uuid, createdAt, name, password)
	return args.String(0)
}

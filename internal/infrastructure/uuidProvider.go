package infrastructure

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type RealUuidProvider struct{}

func NewRealUuidProvider() *RealUuidProvider {
	return &RealUuidProvider{}
}

func (*RealUuidProvider) Random() uuid.UUID {
	return uuid.New()
}

type MockUuidProvider struct {
	mock.Mock
}

func NewMockUuidProvider() *MockUuidProvider {
	return &MockUuidProvider{}
}

func (up *MockUuidProvider) Random() uuid.UUID {
	args := up.Called()
	return args.Get(0).(uuid.UUID)
}

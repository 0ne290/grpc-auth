package infrastructure

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type RealTimeProvider struct{}

func NewRealTimeProvider() *RealTimeProvider {
	return &RealTimeProvider{}
}

func (*RealTimeProvider) Now() time.Time {
	return time.Now().Round(time.Second).UTC()
}

type MockTimeProvider struct {
	mock.Mock
}

func NewMockTimeProvider() *MockTimeProvider {
	return &MockTimeProvider{}
}

func (tp *MockTimeProvider) Now() time.Time {
	args := tp.Called()
	return args.Get(0).(time.Time)
}

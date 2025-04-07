package infrastructure

import "time"

type RealTimeProvider struct{}

func NewRealTimeProvider() *RealTimeProvider {
	return &RealTimeProvider{}
}

func (*RealTimeProvider) Now() time.Time {
	return time.Now().Round(time.Second).UTC()
}

package infrastructure

import "github.com/google/uuid"

type RealUuidProvider struct{}

func NewRealUuidProvider() *RealUuidProvider {
	return &RealUuidProvider{}
}

func (*RealUuidProvider) Random() uuid.UUID {
	return uuid.New()
}

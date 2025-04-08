package infrastructure

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/stretchr/testify/mock"
)

type Sha512Hasher struct{}

func NewSha512Hasher() *Sha512Hasher {
	return &Sha512Hasher{}
}

func (h *Sha512Hasher) Hash(saltedPassword string) string {
	checksum := sha512.Sum512([]byte(saltedPassword))

	return hex.EncodeToString(checksum[:])
}

type MockHasher struct {
	mock.Mock
}

func NewMockHasher() *MockHasher {
	return &MockHasher{}
}

func (h *MockHasher) Hash(saltedPassword string) string {
	args := h.Called(saltedPassword)
	return args.String(0)
}

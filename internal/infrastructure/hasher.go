package infrastructure

import (
	"crypto/sha512"
	"encoding/hex"
)

type Sha512Hasher struct{}

func NewSha512Hasher() *Sha512Hasher {
	return &Sha512Hasher{}
}

func (h *Sha512Hasher) Hash(saltedPassword string) string {
	checksum := sha512.Sum512([]byte(saltedPassword))

	return hex.EncodeToString(checksum[:])
}

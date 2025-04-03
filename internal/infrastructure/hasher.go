package infrastructure

import "hash"

type RealHasher struct {
	hasher hash.Hash
}

func NewRealHasher(hasher hash.Hash) *RealHasher {
	return &RealHasher{hasher}
}

func (h *RealHasher) Hash(saltedPassword string) string {
	ret := string(h.hasher.Sum([]byte(saltedPassword)))

	h.hasher.Reset()

	return ret
}

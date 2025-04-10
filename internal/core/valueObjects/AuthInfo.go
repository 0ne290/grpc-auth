package valueObjects

import (
	"github.com/google/uuid"
	"time"
)

type AuthInfo struct {
	UserUuid     uuid.UUID
	ExpirationAt time.Time
}

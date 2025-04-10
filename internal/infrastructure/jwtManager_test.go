package infrastructure_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"grpc-auth/internal/core/valueObjects"
	"grpc-auth/internal/infrastructure"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	// Arrange
	key := []byte("123_secret_321")
	userUuid, _ := uuid.Parse("e631182f-2be6-4b24-84a9-339881d1c89b")
	expirationAt := time.Date(1986, time.April, 26, 1, 23, 47, 0, time.UTC)
	info := &valueObjects.AuthInfo{UserUuid: userUuid, ExpirationAt: expirationAt}
	manager := infrastructure.NewRealJwtManager(key)

	// Act
	token, err := manager.Generate(info)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	t.Log("token: ", token)
}

func TestParse(t *testing.T) {
	// Arrange
	key := []byte("123_secret_321")
	userUuid, _ := uuid.Parse("e631182f-2be6-4b24-84a9-339881d1c89b")
	expirationAt := time.Date(1986, time.April, 26, 1, 23, 47, 0, time.UTC)
	expectedInfo := &valueObjects.AuthInfo{UserUuid: userUuid, ExpirationAt: expirationAt}
	manager := infrastructure.NewRealJwtManager(key)
	token, _ := manager.Generate(expectedInfo)

	// Act
	actualInfo := manager.Parse(token)

	// Assert
	assert.NotEmpty(t, actualInfo)

	t.Log("expected info: ", expectedInfo)
	t.Log("actual info: ", actualInfo)

	assert.Equal(t, *expectedInfo, *actualInfo)
}

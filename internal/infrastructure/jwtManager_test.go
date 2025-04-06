package infrastructure_test

import (
	"github.com/google/uuid"
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
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

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
	actualInfo, err := manager.Parse(token)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if actualInfo == nil {
		t.Fatal("actual info is nil")
	}

	t.Log("expected info: ", expectedInfo)
	t.Log("actual info: ", actualInfo)

	if *actualInfo != *expectedInfo {
		t.Fatal("actual info not equal expected info")
	}
}

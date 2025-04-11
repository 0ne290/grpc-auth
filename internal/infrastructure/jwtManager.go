package infrastructure

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"grpc-auth/internal/core/value-objects"
	"time"
)

type RealJwtManager struct {
	key []byte
}

func NewRealJwtManager(key []byte) *RealJwtManager {
	return &RealJwtManager{key}
}

func (jm *RealJwtManager) Generate(info *value_objects.AuthInfo) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"userUuid":     info.UserUuid,
			"expirationAt": info.ExpirationAt,
		},
	)

	signedToken, err := token.SignedString(jm.key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (jm *RealJwtManager) TryParse(tokenString string) *value_objects.AuthInfo {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jm.key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
	if err != nil || !token.Valid {
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)

	userUuidAny, ok := claims["userUuid"]
	if !ok {
		return nil
	}
	userUuidString, ok := userUuidAny.(string)
	if !ok {
		return nil
	}
	userUuid, err := uuid.Parse(userUuidString)
	if err != nil {
		return nil
	}

	expirationAtAny, ok := claims["expirationAt"]
	if !ok {
		return nil
	}
	expirationAtString, ok := expirationAtAny.(string)
	if !ok {
		return nil
	}
	expirationAt, err := time.Parse(time.RFC3339, expirationAtString)
	if err != nil {
		return nil
	}

	return &value_objects.AuthInfo{UserUuid: userUuid, ExpirationAt: expirationAt}
}

type MockJwtManager struct {
	mock.Mock
}

func NewMockJwtManager() *MockJwtManager {
	return &MockJwtManager{}
}

func (jm *MockJwtManager) Generate(info *value_objects.AuthInfo) (string, error) {
	args := jm.Called(info)
	return args.String(0), args.Error(1)
}

func (jm *MockJwtManager) TryParse(tokenString string) *value_objects.AuthInfo {
	args := jm.Called(tokenString)
	return args.Get(0).(*value_objects.AuthInfo)
}

package infrastructure

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"grpc-auth/internal/core/valueObjects"
	"time"
)

type RealJwtManager struct {
	key []byte
}

func NewRealJwtManager(key []byte) *RealJwtManager {
	return &RealJwtManager{key}
}

func (jm *RealJwtManager) Generate(info *valueObjects.AuthInfo) (string, error) {
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

func (jm *RealJwtManager) Parse(tokenString string) (*valueObjects.AuthInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jm.key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))

	if err != nil || !token.Valid {
		return nil, nil
	}

	claims := token.Claims.(jwt.MapClaims)
	userUuid, err := uuid.Parse(claims["userUuid"].(string))
	if err != nil {
		return nil, err
	}
	expirationAt, err := time.Parse(time.RFC3339, claims["expirationAt"].(string))
	if err != nil {
		return nil, err
	}

	return &valueObjects.AuthInfo{UserUuid: userUuid, ExpirationAt: expirationAt}, nil
}

type MockJwtManager struct {
	mock.Mock
}

func NewMockJwtManager() *MockJwtManager {
	return &MockJwtManager{}
}

func (jm *MockJwtManager) Generate(info *valueObjects.AuthInfo) (string, error) {
	args := jm.Called(info)
	return args.String(0), args.Error(1)
}

func (jm *MockJwtManager) Parse(tokenString string) (*valueObjects.AuthInfo, error) {
	args := jm.Called(tokenString)
	return args.Get(0).(*valueObjects.AuthInfo), args.Error(1)
}

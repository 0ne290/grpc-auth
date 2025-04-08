package auth_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"grpc-auth/internal/core/entities"
	"grpc-auth/internal/core/services/auth"
	"grpc-auth/internal/core/valueObjects"
	"grpc-auth/internal/infrastructure"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	// Arrange
	const accessTokenLifetime time.Duration = 0
	const refreshTokenLifetime time.Duration = 0
	unitOfWorkStarter := infrastructure.NewMockUnitOfWorkStarter()
	unitOfWork := infrastructure.NewMockUnitOfWork()
	userRepository := infrastructure.NewMockUserRepository()
	timeProvider := infrastructure.NewMockTimeProvider()
	uuidProvider := infrastructure.NewMockUuidProvider()
	hasher := infrastructure.NewMockHasher()
	salter := infrastructure.NewMockSalter()
	jwtManager := infrastructure.NewMockJwtManager()

	password := "password"
	saltedPassword := password + "salt"
	userUuid := uuid.Nil
	userCreatedAt := time.Date(2025, 4, 8, 14, 39, 0, 0, time.UTC)
	userName := "Name"
	userPassword := saltedPassword + "hash"
	user := entities.NewUser(userUuid, userCreatedAt, userName, userPassword)
	ctx := context.TODO()

	unitOfWorkStarter.On("Start", ctx).Return(unitOfWork, nil)
	unitOfWork.On("UserRepository").Return(userRepository)
	unitOfWork.On("Save", ctx).Return(nil)
	userRepository.On("TryCreate", ctx, user).Return(true, nil)
	timeProvider.On("Now").Return(userCreatedAt)
	uuidProvider.On("Random").Return(userUuid)
	hasher.On("Hash", saltedPassword).Return(userPassword)
	salter.On("Salt", userUuid, userCreatedAt, userName, password).Return(saltedPassword)

	request := &auth.RegisterRequest{Name: userName, Password: password}
	service := auth.NewRealService(accessTokenLifetime, refreshTokenLifetime, unitOfWorkStarter, timeProvider, uuidProvider, hasher, salter, jwtManager)

	// Act
	response, err := service.Register(ctx, request)
	t.Log(response)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	unitOfWorkStarter.AssertCalled(t, "Start", ctx)
	unitOfWork.AssertCalled(t, "UserRepository")
	uuidProvider.AssertCalled(t, "Random")
	timeProvider.AssertCalled(t, "Now")
	salter.AssertCalled(t, "Salt", userUuid, userCreatedAt, userName, password)
	hasher.AssertCalled(t, "Hash", saltedPassword)
	userRepository.AssertCalled(t, "TryCreate", ctx, user)
	unitOfWork.AssertCalled(t, "Save", ctx)
}

func TestLogin(t *testing.T) {
	// Arrange
	const accessTokenLifetime time.Duration = 0
	const refreshTokenLifetime time.Duration = 0
	unitOfWorkStarter := infrastructure.NewMockUnitOfWorkStarter()
	unitOfWork := infrastructure.NewMockUnitOfWork()
	userRepository := infrastructure.NewMockUserRepository()
	sessionRepository := infrastructure.NewMockSessionRepository()
	timeProvider := infrastructure.NewMockTimeProvider()
	uuidProvider := infrastructure.NewMockUuidProvider()
	hasher := infrastructure.NewMockHasher()
	salter := infrastructure.NewMockSalter()
	jwtManager := infrastructure.NewMockJwtManager()

	password := "password"
	saltedPassword := password + "salt"
	fakeUuid := uuid.Nil
	fakeNow := time.Date(2025, 4, 8, 14, 39, 0, 0, time.UTC)
	userName := "Name"
	userPassword := saltedPassword + "hash"
	user := entities.NewUser(fakeUuid, fakeNow, userName, userPassword)
	session := entities.NewSession(fakeUuid, fakeUuid, fakeNow)
	authInfo := &valueObjects.AuthInfo{UserUuid: fakeUuid, ExpirationAt: fakeNow}
	accessToken := "Fake access token"
	ctx := context.TODO()

	unitOfWorkStarter.On("Start", ctx).Return(unitOfWork, nil)
	unitOfWork.On("UserRepository").Return(userRepository)
	unitOfWork.On("SessionRepository").Return(sessionRepository)
	unitOfWork.On("Save", ctx).Return(nil)
	userRepository.On("TryGetByName", ctx, userName).Return(user, nil)
	sessionRepository.On("Create", ctx, session).Return(nil)
	timeProvider.On("Now").Return(fakeNow)
	uuidProvider.On("Random").Return(fakeUuid)
	hasher.On("Hash", saltedPassword).Return(userPassword)
	salter.On("Salt", fakeUuid, fakeNow, userName, password).Return(saltedPassword)
	jwtManager.On("Generate", authInfo).Return(accessToken, nil)

	request := &auth.LoginRequest{Name: userName, Password: password}
	service := auth.NewRealService(accessTokenLifetime, refreshTokenLifetime, unitOfWorkStarter, timeProvider, uuidProvider, hasher, salter, jwtManager)

	// Act
	response, err := service.Login(ctx, request)
	t.Log(response)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	unitOfWorkStarter.AssertCalled(t, "Start", ctx)
	unitOfWork.AssertCalled(t, "UserRepository")
	unitOfWork.AssertCalled(t, "SessionRepository")
	userRepository.AssertCalled(t, "TryGetByName", ctx, userName)
	salter.AssertCalled(t, "Salt", fakeUuid, fakeNow, userName, password)
	hasher.AssertCalled(t, "Hash", saltedPassword)
	timeProvider.AssertCalled(t, "Now")
	jwtManager.AssertCalled(t, "Generate", authInfo)
	uuidProvider.AssertCalled(t, "Random")
	sessionRepository.AssertCalled(t, "Create", ctx, session)
	unitOfWork.AssertCalled(t, "Save", ctx)
}

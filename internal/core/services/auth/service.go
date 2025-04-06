package auth

import (
	"context"
	"grpc-auth/internal"
	"grpc-auth/internal/core/entities"
	"grpc-auth/internal/core/services"
	"grpc-auth/internal/core/valueObjects"
	"time"
)

const timeDay = time.Hour * 24

type RealService struct {
	accessTokenLifetimeInHours time.Duration
	refreshTokenLifetimeInDays time.Duration
	unitOfWork                 services.UserUnitOfWork
	timeProvider               services.TimeProvider
	uuidProvider               services.UuidProvider
	hasher                     services.Hasher
	salter                     services.Salter
	jwtManager                 services.JwtManager
}

func NewRealService(authConfig internal.AuthConfig, unitOfWork services.UserUnitOfWork, timeProvider services.TimeProvider, uuidProvider services.UuidProvider, hasher services.Hasher, salter services.Salter, jwtManager services.JwtManager) *RealService {
	return &RealService{time.Hour * authConfig.AccessTokenLifetimeInHours, timeDay * authConfig.RefreshTokenLifetimeInDays, unitOfWork, timeProvider, uuidProvider, hasher, salter, jwtManager}
}

func (s *RealService) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	userUuid := s.uuidProvider.Random()
	createdAt := s.timeProvider.Now()

	saltedPassword := s.salter.Salt(userUuid, createdAt, request.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)

	user := entities.NewUser(userUuid, createdAt, request.Name, hashOfSaltedPassword)

	ok, err := repository.TryCreate(ctx, user)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}
	if !ok {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{"user created"}, nil
}

func (s *RealService) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	repository, err := s.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}

	user, err := repository.TryGetByName(ctx, request.Name)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}
	if user == nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	now := s.timeProvider.Now()

	authInfo := &valueObjects.AuthInfo{user.Uuid, now.Add(s.accessTokenLifetimeInHours)}
	accessToken, err := s.jwtManager.Generate(authInfo)
	if err != nil {
		return nil, err
	}

	refreshToken := s.uuidProvider.Random()
	session := entities.NewSession(refreshToken, user.Uuid, now.Add(s.refreshTokenLifetimeInDays))

	/* TODO:

	expirationAt := s.timeProvider.Now() + Days(s.authConfig.RefreshTokenLifetimeInDays)


	Сохранить сессию в БД

	Сгенерировать токен доступа

	Вернуть токены обновления и доступа

	*/

	return &LoginResponse{"stub"}, nil
}

func (s *RealService) CheckToken(request *CheckTokenRequest) (*CheckTokenResponse, error) {
	if request.Token != "stub" {
		return nil, &services.InvariantViolationError{Message: "permission denied"}
	}

	return &CheckTokenResponse{"permission granted"}, nil
}

package auth

import (
	"context"
	"github.com/google/uuid"
	"grpc-auth/internal/core/entities"
	"grpc-auth/internal/core/services"
	"grpc-auth/internal/core/value-objects"
	"time"
)

type RealService struct {
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
	unitOfWorkStarter    services.UnitOfWorkStarter
	timeProvider         services.TimeProvider
	uuidProvider         services.UuidProvider
	hasher               services.Hasher
	salter               services.Salter
	jwtManager           services.JwtManager
}

func NewRealService(accessTokenLifetime, refreshTokenLifetime time.Duration, unitOfWorkStarter services.UnitOfWorkStarter, timeProvider services.TimeProvider, uuidProvider services.UuidProvider, hasher services.Hasher, salter services.Salter, jwtManager services.JwtManager) *RealService {
	return &RealService{accessTokenLifetime, refreshTokenLifetime, unitOfWorkStarter, timeProvider, uuidProvider, hasher, salter, jwtManager}
}

func (s *RealService) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()

	userUuid := s.uuidProvider.Random()
	createdAt := s.timeProvider.Now()

	saltedPassword := s.salter.Salt(userUuid, createdAt, request.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)

	user := entities.NewUser(userUuid, createdAt, request.Name, hashOfSaltedPassword)

	ok, err := userRepository.TryCreate(ctx, user)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if !ok {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{"user created"}, nil
}

func (s *RealService) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()
	sessionRepository := unitOfWork.SessionRepository()

	user, err := userRepository.TryGetByName(ctx, request.Name)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if user == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	now := s.timeProvider.Now()

	authInfo := &value_objects.AuthInfo{UserUuid: user.Uuid, ExpirationAt: now.Add(s.accessTokenLifetime)}
	accessToken, err := s.jwtManager.Generate(authInfo)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	refreshToken := s.uuidProvider.Random()
	session := entities.NewSession(refreshToken, user.Uuid, now.Add(s.refreshTokenLifetime))

	err = sessionRepository.Create(ctx, session)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{refreshToken.String(), accessToken}, nil
}

func (s *RealService) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*DeleteUserResponse, error) {
	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()

	user, err := userRepository.TryGetByName(ctx, request.Name)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if user == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	err = userRepository.Delete(ctx, user.Uuid)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &DeleteUserResponse{"user deleted"}, nil
}

func (s *RealService) DeleteSession(ctx context.Context, request *DeleteSessionRequest) (*DeleteSessionResponse, error) {
	refreshToken, err := uuid.Parse(request.RefreshToken)
	if err != nil {
		return nil, &services.InvariantViolationError{Message: "refresh token format is invalid"}
	}

	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()
	sessionRepository := unitOfWork.SessionRepository()

	session, err := sessionRepository.TryGetByRefreshToken(ctx, refreshToken)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if session == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "refresh token does not exists"}
	}

	user, err := userRepository.TryGetByName(ctx, request.Name)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if user == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword || user.Uuid != session.UserUuid {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	err = sessionRepository.Delete(ctx, refreshToken)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &DeleteSessionResponse{"session deleted"}, nil
}

func (s *RealService) ChangeName(ctx context.Context, request *ChangeNameRequest) (*ChangeNameResponse, error) {
	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()

	user, err := userRepository.TryGetByName(ctx, request.Name)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if user == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	user.Name = request.NewName
	saltedPassword = s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword = s.hasher.Hash(saltedPassword)
	user.Password = hashOfSaltedPassword

	updated, err := userRepository.TryUpdate(ctx, user)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if !updated {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &ChangeNameResponse{"name updated"}, nil
}

func (s *RealService) ChangePassword(ctx context.Context, request *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()

	user, err := userRepository.TryGetByName(ctx, request.Name)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if user == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "login or/and password is invalid"}
	}

	saltedPassword = s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.NewPassword)
	hashOfSaltedPassword = s.hasher.Hash(saltedPassword)
	user.Password = hashOfSaltedPassword

	_, err = userRepository.TryUpdate(ctx, user)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{"password updated"}, nil
}

func (s *RealService) RefreshTokens(ctx context.Context, request *RefreshTokensRequest) (*RefreshTokensResponse, error) {
	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	sessionRepository := unitOfWork.SessionRepository()

	refreshToken, err := uuid.Parse(request.RefreshToken)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "refresh token format is invalid"}
	}

	session, err := sessionRepository.TryGetByRefreshToken(ctx, refreshToken)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}
	if session == nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "refresh token does not exists"}
	}

	err = sessionRepository.Delete(ctx, refreshToken)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	now := s.timeProvider.Now()

	if session.ExpirationAt.Before(now) {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "refresh token expired"}
	}

	refreshToken = s.uuidProvider.Random()

	session = entities.NewSession(refreshToken, session.UserUuid, now.Add(s.refreshTokenLifetime))

	err = sessionRepository.Create(ctx, session)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	authInfo := &value_objects.AuthInfo{UserUuid: session.UserUuid, ExpirationAt: now.Add(s.accessTokenLifetime)}
	accessToken, err := s.jwtManager.Generate(authInfo)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &RefreshTokensResponse{refreshToken.String(), accessToken}, nil
}

func (s *RealService) CheckAccessToken(ctx context.Context, request *CheckAccessTokenRequest) (*CheckAccessTokenResponse, error) {
	authInfo := s.jwtManager.TryParse(request.AccessToken)
	if authInfo == nil {
		return nil, &services.InvariantViolationError{Message: "access token is invalid"}
	}

	if authInfo.ExpirationAt.Before(s.timeProvider.Now()) {
		return &CheckAccessTokenResponse{false}, nil
	}

	unitOfWork, err := s.unitOfWorkStarter.Start(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := unitOfWork.UserRepository()

	exists, err := userRepository.Exists(ctx, authInfo.UserUuid)
	if err != nil {
		_ = unitOfWork.Rollback(ctx)

		return nil, err
	}

	if !exists {
		_ = unitOfWork.Rollback(ctx)

		return nil, &services.InvariantViolationError{Message: "user not found"}
	}

	err = unitOfWork.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &CheckAccessTokenResponse{true}, nil
}

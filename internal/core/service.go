package core

import "context"

type RealService struct {
	unitOfWork   UnitOfWork
	timeProvider TimeProvider
	uuidProvider UuidProvider
	hasher       Hasher
	salter       Salter
}

func NewRealService(unitOfWork UnitOfWork, timeProvider TimeProvider, uuidProvider UuidProvider, hasher Hasher, salter Salter) *RealService {
	return &RealService{unitOfWork, timeProvider, uuidProvider, hasher, salter}
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

	user := newUser(userUuid, createdAt, request.Name, hashOfSaltedPassword)

	ok, err := repository.TryCreate(ctx, user)
	if err != nil {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, err
	}
	if !ok {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, &InvariantViolationError{"login or/and password is invalid"}
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

		return nil, &InvariantViolationError{"login or/and password is invalid"}
	}

	saltedPassword := s.salter.Salt(user.Uuid, user.CreatedAt, user.Name, request.Password)
	hashOfSaltedPassword := s.hasher.Hash(saltedPassword)
	if user.Password != hashOfSaltedPassword {
		_ = s.unitOfWork.Rollback(ctx, repository)

		return nil, &InvariantViolationError{"login or/and password is invalid"}
	}

	err = s.unitOfWork.Save(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{"stub"}, nil
}

func (s *RealService) CheckToken(request *CheckTokenRequest) (*CheckTokenResponse, error) {
	if request.Token != "stub" {
		return nil, &InvariantViolationError{"permission denied"}
	}

	return &CheckTokenResponse{"permission granted"}, nil
}

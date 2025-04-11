package auth

import (
	"context"
	service "grpc-auth/internal/core/services/auth"
)

type Service interface {
	Register(ctx context.Context, request *service.RegisterRequest) (*service.RegisterResponse, error)
	Login(ctx context.Context, request *service.LoginRequest) (*service.LoginResponse, error)
	DeleteUser(ctx context.Context, request *service.DeleteUserRequest) (*service.DeleteUserResponse, error)
	RefreshTokens(ctx context.Context, request *service.RefreshTokensRequest) (*service.RefreshTokensResponse, error)
	CheckAccessToken(ctx context.Context, request *service.CheckAccessTokenRequest) (*service.CheckAccessTokenResponse, error)
}

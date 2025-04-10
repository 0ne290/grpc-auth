package auth

import (
	"context"
	service "grpc-auth/internal/core/services/auth"
)

type Service interface {
	Register(ctx context.Context, request *service.RegisterRequest) (*service.RegisterResponse, error)
	Login(ctx context.Context, request *service.LoginRequest) (*service.LoginResponse, error)
	RefreshTokens(ctx context.Context, request *service.RefreshTokensRequest) (*service.RefreshTokensResponse, error)
	CheckAccessToken(request *service.CheckAccessTokenRequest) (*service.CheckAccessTokenResponse, error)
}

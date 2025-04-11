package auth

import (
	"context"
	service "grpc-auth/internal/core/services/auth"
)

type Service interface {
	Register(ctx context.Context, request *service.RegisterRequest) (*service.RegisterResponse, error)
	Login(ctx context.Context, request *service.LoginRequest) (*service.LoginResponse, error)
	DeleteUser(ctx context.Context, request *service.DeleteUserRequest) (*service.DeleteUserResponse, error)
	DeleteSession(ctx context.Context, request *service.DeleteSessionRequest) (*service.DeleteSessionResponse, error)
	ChangeName(ctx context.Context, request *service.ChangeNameRequest) (*service.ChangeNameResponse, error)
	ChangePassword(ctx context.Context, request *service.ChangePasswordRequest) (*service.ChangePasswordResponse, error)
	RefreshTokens(ctx context.Context, request *service.RefreshTokensRequest) (*service.RefreshTokensResponse, error)
	CheckAccessToken(ctx context.Context, request *service.CheckAccessTokenRequest) (*service.CheckAccessTokenResponse, error)
}

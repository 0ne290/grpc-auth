package auth

import (
	"context"
	auth2 "grpc-auth/internal/core/services/auth"
)

type Service interface {
	Register(ctx context.Context, request *auth2.RegisterRequest) (*auth2.RegisterResponse, error)
	Login(ctx context.Context, request *auth2.LoginRequest) (*auth2.LoginResponse, error)
	CheckToken(request *auth2.CheckTokenRequest) (*auth2.CheckTokenResponse, error)
}

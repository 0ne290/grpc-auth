package auth

import (
	"context"
	core "grpc-auth/internal/core/auth"
)

type Service interface {
	Register(ctx context.Context, request *core.RegisterRequest) (*core.RegisterResponse, error)
	Login(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error)
	CheckToken(request *core.CheckTokenRequest) (*core.CheckTokenResponse, error)
}

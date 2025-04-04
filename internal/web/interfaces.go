package web

import (
	"context"
	"grpc-auth/internal/core"
)

type Service interface {
	Register(ctx context.Context, request *core.RegisterRequest) (*core.RegisterResponse, error)
	Login(ctx context.Context, request *core.LoginRequest) (*core.LoginResponse, error)
	CheckToken(request *core.CheckTokenRequest) (*core.CheckTokenResponse, error)
}

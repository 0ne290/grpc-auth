package auth

import (
	"context"
	"google.golang.org/grpc"
	"grpc-auth/grpc/gen"
	auth2 "grpc-auth/internal/core/services/auth"
)

type Controller struct {
	auth.UnimplementedAuthServer
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func RegisterController(grpcServer *grpc.Server, controller *Controller) {
	auth.RegisterAuthServer(grpcServer, controller)
}

func (s *Controller) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	ret, err := s.service.Register(ctx, mapRegisterRequest(req))

	return mapRegisterResponse(ret), err
}

func mapRegisterRequest(source *auth.RegisterRequest) *auth2.RegisterRequest {
	if source == nil {
		return nil
	}

	return &auth2.RegisterRequest{Name: source.Username, Password: source.Password}
}

func mapRegisterResponse(source *auth2.RegisterResponse) *auth.RegisterResponse {
	if source == nil {
		return nil
	}

	return &auth.RegisterResponse{Message: source.Message}
}

func (s *Controller) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	ret, err := s.service.Login(ctx, mapLoginRequest(req))

	return mapLoginResponse(ret), err
}

func mapLoginRequest(source *auth.LoginRequest) *auth2.LoginRequest {
	if source == nil {
		return nil
	}

	return &auth2.LoginRequest{Name: source.Username, Password: source.Password}
}

func mapLoginResponse(source *auth2.LoginResponse) *auth.LoginResponse {
	if source == nil {
		return nil
	}

	return &auth.LoginResponse{Token: source.Token}
}

func (s *Controller) CheckToken(_ context.Context, req *auth.CheckTokenRequest) (*auth.CheckTokenResponse, error) {
	ret, err := s.service.CheckToken(mapCheckTokenRequest(req))

	return mapCheckTokenResponse(ret), err
}

func mapCheckTokenRequest(source *auth.CheckTokenRequest) *auth2.CheckTokenRequest {
	if source == nil {
		return nil
	}

	return &auth2.CheckTokenRequest{Token: source.Token}
}

func mapCheckTokenResponse(source *auth2.CheckTokenResponse) *auth.CheckTokenResponse {
	if source == nil {
		return nil
	}

	return &auth.CheckTokenResponse{Message: source.Message}
}

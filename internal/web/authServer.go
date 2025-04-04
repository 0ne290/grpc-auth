package web

import (
	"context"
	"google.golang.org/grpc"
	"grpc-auth/grpc/gen"
	"grpc-auth/internal/core"
)

type Server struct {
	auth.UnimplementedAuthServer
	service Service
}

func NewServer(service Service) *Server {
	return &Server{service: service}
}

func RegisterServer(grpcServer *grpc.Server, server *Server) {
	auth.RegisterAuthServer(grpcServer, server)
}

func (s *Server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	ret, err := s.service.Register(ctx, mapRegisterRequest(req))

	return mapRegisterResponse(ret), err
}

func mapRegisterRequest(source *auth.RegisterRequest) *core.RegisterRequest {
	if source == nil {
		return nil
	}

	return &core.RegisterRequest{Name: source.Username, Password: source.Password}
}

func mapRegisterResponse(source *core.RegisterResponse) *auth.RegisterResponse {
	if source == nil {
		return nil
	}

	return &auth.RegisterResponse{Message: source.Message}
}

func (s *Server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	ret, err := s.service.Login(ctx, mapLoginRequest(req))

	return mapLoginResponse(ret), err
}

func mapLoginRequest(source *auth.LoginRequest) *core.LoginRequest {
	if source == nil {
		return nil
	}

	return &core.LoginRequest{Name: source.Username, Password: source.Password}
}

func mapLoginResponse(source *core.LoginResponse) *auth.LoginResponse {
	if source == nil {
		return nil
	}

	return &auth.LoginResponse{Token: source.Token}
}

func (s *Server) CheckToken(_ context.Context, req *auth.CheckTokenRequest) (*auth.CheckTokenResponse, error) {
	ret, err := s.service.CheckToken(mapCheckTokenRequest(req))

	return mapCheckTokenResponse(ret), err
}

func mapCheckTokenRequest(source *auth.CheckTokenRequest) *core.CheckTokenRequest {
	if source == nil {
		return nil
	}

	return &core.CheckTokenRequest{Token: source.Token}
}

func mapCheckTokenResponse(source *core.CheckTokenResponse) *auth.CheckTokenResponse {
	if source == nil {
		return nil
	}

	return &auth.CheckTokenResponse{Message: source.Message}
}

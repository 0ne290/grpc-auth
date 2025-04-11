package auth

import (
	"context"
	"google.golang.org/grpc"
	"grpc-auth/grpc/gen"
	service "grpc-auth/internal/core/services/auth"
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

func mapRegisterRequest(source *auth.RegisterRequest) *service.RegisterRequest {
	if source == nil {
		return nil
	}

	return &service.RegisterRequest{Name: source.Username, Password: source.Password}
}

func mapRegisterResponse(source *service.RegisterResponse) *auth.RegisterResponse {
	if source == nil {
		return nil
	}

	return &auth.RegisterResponse{Message: source.Message}
}

func (s *Controller) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	ret, err := s.service.Login(ctx, mapLoginRequest(req))

	return mapLoginResponse(ret), err
}

func mapLoginRequest(source *auth.LoginRequest) *service.LoginRequest {
	if source == nil {
		return nil
	}

	return &service.LoginRequest{Name: source.Username, Password: source.Password}
}

func mapLoginResponse(source *service.LoginResponse) *auth.LoginResponse {
	if source == nil {
		return nil
	}

	return &auth.LoginResponse{RefreshToken: source.RefreshToken, AccessToken: source.AccessToken}
}

func (s *Controller) DeleteUser(ctx context.Context, req *auth.DeleteUserRequest) (*auth.DeleteUserResponse, error) {
	ret, err := s.service.DeleteUser(ctx, mapDeleteUserRequest(req))

	return mapDeleteUserResponse(ret), err
}

func mapDeleteUserRequest(source *auth.DeleteUserRequest) *service.DeleteUserRequest {
	if source == nil {
		return nil
	}

	return &service.DeleteUserRequest{AccessToken: source.AccessToken}
}

func mapDeleteUserResponse(source *service.DeleteUserResponse) *auth.DeleteUserResponse {
	if source == nil {
		return nil
	}

	return &auth.DeleteUserResponse{Message: source.Message}
}

func (s *Controller) DeleteSession(ctx context.Context, req *auth.DeleteSessionRequest) (*auth.DeleteSessionResponse, error) {
	ret, err := s.service.DeleteSession(ctx, mapDeleteSessionRequest(req))

	return mapDeleteSessionResponse(ret), err
}

func mapDeleteSessionRequest(source *auth.DeleteSessionRequest) *service.DeleteSessionRequest {
	if source == nil {
		return nil
	}

	return &service.DeleteSessionRequest{RefreshToken: source.RefreshToken}
}

func mapDeleteSessionResponse(source *service.DeleteSessionResponse) *auth.DeleteSessionResponse {
	if source == nil {
		return nil
	}

	return &auth.DeleteSessionResponse{Message: source.Message}
}

func (s *Controller) ChangeName(ctx context.Context, req *auth.ChangeNameRequest) (*auth.ChangeNameResponse, error) {
	ret, err := s.service.ChangeName(ctx, mapChangeNameRequest(req))

	return mapChangeNameResponse(ret), err
}

func mapChangeNameRequest(source *auth.ChangeNameRequest) *service.ChangeNameRequest {
	if source == nil {
		return nil
	}

	return &service.ChangeNameRequest{AccessToken: source.AccessToken, NewName: source.NewName}
}

func mapChangeNameResponse(source *service.ChangeNameResponse) *auth.ChangeNameResponse {
	if source == nil {
		return nil
	}

	return &auth.ChangeNameResponse{Message: source.Message}
}

func (s *Controller) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	ret, err := s.service.ChangePassword(ctx, mapChangePasswordRequest(req))

	return mapChangePasswordResponse(ret), err
}

func mapChangePasswordRequest(source *auth.ChangePasswordRequest) *service.ChangePasswordRequest {
	if source == nil {
		return nil
	}

	return &service.ChangePasswordRequest{AccessToken: source.AccessToken, NewPassword: source.NewPassword}
}

func mapChangePasswordResponse(source *service.ChangePasswordResponse) *auth.ChangePasswordResponse {
	if source == nil {
		return nil
	}

	return &auth.ChangePasswordResponse{Message: source.Message}
}

func (s *Controller) RefreshTokens(ctx context.Context, req *auth.RefreshTokensRequest) (*auth.RefreshTokensResponse, error) {
	ret, err := s.service.RefreshTokens(ctx, mapRefreshTokensRequest(req))

	return mapRefreshTokensResponse(ret), err
}

func mapRefreshTokensRequest(source *auth.RefreshTokensRequest) *service.RefreshTokensRequest {
	if source == nil {
		return nil
	}

	return &service.RefreshTokensRequest{RefreshToken: source.RefreshToken}
}

func mapRefreshTokensResponse(source *service.RefreshTokensResponse) *auth.RefreshTokensResponse {
	if source == nil {
		return nil
	}

	return &auth.RefreshTokensResponse{RefreshToken: source.RefreshToken, AccessToken: source.AccessToken}
}

func (s *Controller) CheckAccessToken(ctx context.Context, req *auth.CheckAccessTokenRequest) (*auth.CheckAccessTokenResponse, error) {
	ret, err := s.service.CheckAccessToken(ctx, mapCheckAccessTokenRequest(req))

	return mapCheckAccessTokenResponse(ret), err
}

func mapCheckAccessTokenRequest(source *auth.CheckAccessTokenRequest) *service.CheckAccessTokenRequest {
	if source == nil {
		return nil
	}

	return &service.CheckAccessTokenRequest{AccessToken: source.AccessToken}
}

func mapCheckAccessTokenResponse(source *service.CheckAccessTokenResponse) *auth.CheckAccessTokenResponse {
	if source == nil {
		return nil
	}

	return &auth.CheckAccessTokenResponse{IsActive: source.IsActive}
}

package auth

type RegisterRequest struct {
	Name, Password string
}

type LoginRequest struct {
	Name, Password string
}

type DeleteUserRequest struct {
	AccessToken string
}

type DeleteSessionRequest struct {
	RefreshToken string
}

type ChangeNameRequest struct {
	AccessToken, NewName string
}

type ChangePasswordRequest struct {
	AccessToken, NewPassword string
}

type RefreshTokensRequest struct {
	RefreshToken string
}

type CheckAccessTokenRequest struct {
	AccessToken string
}

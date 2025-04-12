package auth

type RegisterRequest struct {
	Name, Password string
}

type LoginRequest struct {
	Name, Password string
}

type DeleteUserRequest struct {
	Name, Password string
}

type DeleteSessionRequest struct {
	Name, Password, RefreshToken string
}

type ChangeNameRequest struct {
	Name, Password, NewName string
}

type ChangePasswordRequest struct {
	Name, Password, NewPassword string
}

type RefreshTokensRequest struct {
	RefreshToken string
}

type CheckAccessTokenRequest struct {
	AccessToken string
}

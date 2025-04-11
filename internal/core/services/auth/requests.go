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

type RefreshTokensRequest struct {
	RefreshToken string
}

type CheckAccessTokenRequest struct {
	AccessToken string
}

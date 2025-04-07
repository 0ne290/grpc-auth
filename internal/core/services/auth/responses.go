package auth

type RegisterResponse struct {
	Message string
}

type LoginResponse struct {
	RefreshToken, AccessToken string
}

type RefreshTokensResponse struct {
	RefreshToken, AccessToken string
}

type CheckAccessTokenResponse struct {
	IsActive bool
}

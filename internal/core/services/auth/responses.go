package auth

type RegisterResponse struct {
	Message string
}

type LoginResponse struct {
	Token string
}

type CheckTokenResponse struct {
	Message string
}

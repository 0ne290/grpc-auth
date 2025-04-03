package core

type RegisterRequest struct {
	Name, Password string
}

type LoginRequest struct {
	Name, Password string
}

type CheckTokenRequest struct {
	Token string
}

package service

// LoginRequest contains the credentials submitted during login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse contains the generated JWT and account details.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Username    string `json:"username"`
	Role        string `json:"role"`
}

// AuthService defines authentication operations.
type AuthService interface {
	Login(request LoginRequest) (*LoginResponse, error)
}

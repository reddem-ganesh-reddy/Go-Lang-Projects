package service

import (
	"fmt"
	"time"

	"EMS/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// account contains an in-memory account's authentication details.
type account struct {
	username     string
	passwordHash []byte
	role         string
}

// AuthServiceImpl implements the authentication service.
type AuthServiceImpl struct {
	jwtSecret []byte
	expiresIn time.Duration
	accounts  map[string]account
}

// NewAuthService creates the authentication service using environment configuration.
func NewAuthService(cfg config.Config) (*AuthServiceImpl, error) {
	if cfg.AuthAdminUsername == "" || cfg.AuthAdminPassword == "" {
		return nil, fmt.Errorf("admin credentials are missing")
	}

	if cfg.AuthEmployeeUsername == "" || cfg.AuthEmployeePassword == "" {
		return nil, fmt.Errorf("employee credentials are missing")
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT secret must contain at least 32 characters")
	}

	if cfg.JWTExpirationMinutes <= 0 {
		return nil, fmt.Errorf("JWT expiration must be greater than zero")
	}

	adminHash, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.AuthAdminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin credentials: %w", err)
	}

	employeeHash, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.AuthEmployeePassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize employee credentials: %w", err)
	}

	// Prevent duplicate usernames from creating ambiguous accounts.
	if cfg.AuthAdminUsername == cfg.AuthEmployeeUsername {
		return nil, fmt.Errorf("admin and employee usernames must be different")
	}

	return &AuthServiceImpl{
		jwtSecret: []byte(cfg.JWTSecret),
		expiresIn: time.Duration(cfg.JWTExpirationMinutes) * time.Minute,
		accounts: map[string]account{
			cfg.AuthAdminUsername: {
				username:     cfg.AuthAdminUsername,
				passwordHash: adminHash,
				role:         "admin",
			},
			cfg.AuthEmployeeUsername: {
				username:     cfg.AuthEmployeeUsername,
				passwordHash: employeeHash,
				role:         "employee",
			},
		},
	}, nil
}

// Login verifies credentials and generates a JWT token.
func (s *AuthServiceImpl) Login(request LoginRequest) (*LoginResponse, error) {
	user, exists := s.accounts[request.Username]
	if !exists {
		return nil, fmt.Errorf("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword(
		user.passwordHash,
		[]byte(request.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	now := time.Now()
	expiresAt := now.Add(s.expiresIn)

	// Create JWT claims containing the username and role.
	claims := jwt.MapClaims{
		"username": user.username,
		"role":     user.role,
		"iat":      now.Unix(),
		"exp":      expiresAt.Unix(),
	}

	// Sign the JWT using the configured secret.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &LoginResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.expiresIn.Seconds()),
		Username:    user.username,
		Role:        user.role,
	}, nil
}

package security

import (
	"fmt"
	"time"

	"jarvis/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims defines the custom claims for JWT tokens.
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new JWT access token.
func GenerateAccessToken(cfg *config.SecurityConfig, userID uuid.UUID, username, email string) (string, error) {
	expirationTime := time.Now().Add(cfg.JWTExpiresIn)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenString, nil
}

// GenerateRefreshToken generates a new JWT refresh token.
func GenerateRefreshToken(cfg *config.SecurityConfig, userID uuid.UUID, username, email string) (string, error) {
	expirationTime := time.Now().Add(cfg.JWTRefreshExpiresIn)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims.
func ValidateToken(cfg *config.SecurityConfig, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

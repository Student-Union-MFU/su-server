package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	expiry time.Duration
}

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	expiryHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil {
		expiryHours = 24
	}

	return &JWTService{
		secret: []byte(secret),
		expiry: time.Duration(expiryHours) * time.Hour,
	}
}

// Generate creates a new JWT token for a user
func (s *JWTService) Generate(userID int, userType string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

// Validate parses and validates a JWT token
func (s *JWTService) Validate(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

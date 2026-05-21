package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/north-fy/talker/services/user/internal/config"

	"time"
)

type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret)
}

func ValidateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return config.JwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Извлекаем claims
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

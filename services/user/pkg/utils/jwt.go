package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/north-fy/talker/services/user/internal/config"
	"github.com/north-fy/talker/services/user/internal/domain/models"

	"time"
)

type UserClaims struct {
	User models.User `json:"user"`
	Role string      `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &UserClaims{
		User: user,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   strconv.Itoa(int(user.UID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
}

func ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Извлекаем claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

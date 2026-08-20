package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ctxUserID        = "user_id"
	ctxToken         = "auth_token"
	authHeaderKey    = "Authorization"
	authHeaderPrefix = "Bearer "
)

type Claims struct {
	UserID int64
	jwt.RegisteredClaims
}

// ParseToken валидирует JWT и извлекает user_id из claims.
func ParseToken(tokenString, secret string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return strconv.ParseInt(claims.Subject, 10, 64)
}

// Auth вернёт middleware, требующее валидный Bearer JWT.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authHeaderKey)
		if !strings.HasPrefix(header, authHeaderPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid authorization header",
				"code":  http.StatusUnauthorized,
			})

			return
		}

		token := strings.TrimPrefix(header, authHeaderPrefix)
		userID, err := ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
				"code":  http.StatusUnauthorized,
			})

			return
		}

		c.Set(ctxUserID, userID)
		c.Set(ctxToken, token)
		c.Request = c.Request.WithContext(WithAuthToken(c.Request.Context(), token))
		c.Next()
	}
}

func UserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return 0, false
	}

	id, ok := v.(int64)
	return id, ok
}

func Token(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxToken)
	if !ok {
		return "", false
	}

	t, ok := v.(string)
	return t, ok
}

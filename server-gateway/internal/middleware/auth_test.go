package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func generateToken(t *testing.T, secret string, userID int64) string {
	t.Helper()

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "auth-service",
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return s
}

func newTestEngine(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	e := gin.New()
	e.Use(Auth(secret))
	e.GET("/protected", func(c *gin.Context) {
		id, _ := UserID(c)
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})

	return e
}

func TestAuth_ValidToken(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(authHeaderKey, "Bearer "+generateToken(t, testSecret, 42))

	newTestEngine(testSecret).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body map[string]int64
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if body["user_id"] != 42 {
		t.Fatalf("user_id = %d, want 42", body["user_id"])
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	newTestEngine(testSecret).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(authHeaderKey, "Bearer not-a-jwt")

	newTestEngine(testSecret).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuth_WrongSecret(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set(authHeaderKey, "Bearer "+generateToken(t, "other-secret", 42))

	newTestEngine(testSecret).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestParseToken(t *testing.T) {
	t.Parallel()

	token := generateToken(t, testSecret, 99)
	userID, err := ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != 99 {
		t.Fatalf("userID = %d, want 99", userID)
	}

	if _, err := ParseToken("invalid", testSecret); err == nil {
		t.Fatal("expected error for invalid token")
	}
}

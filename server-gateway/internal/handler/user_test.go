package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newUserTestEngine(h *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	e := gin.New()
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	return e
}

func TestUserHandler_Register(t *testing.T) {
	h := NewUserHandler(zap.NewNop(), &fakeUserClient{
		registerFn: func() (*userv1.RegisterResponse, error) {
			return &userv1.RegisterResponse{UserId: 7}, nil
		},
	})

	engine := newUserTestEngine(h)
	body := []byte(`{"first_name":"John","last_name":"Doe","email":"john@example.com","password":"password123"}`)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp map[string]int64
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp["user_id"] != 7 {
		t.Fatalf("user_id = %d, want 7", resp["user_id"])
	}
}

func TestUserHandler_Register_ValidationError(t *testing.T) {
	h := NewUserHandler(zap.NewNop(), &fakeUserClient{})

	engine := newUserTestEngine(h)
	body := []byte(`{"email":"john@example.com"}`)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestUserHandler_Login(t *testing.T) {
	h := NewUserHandler(zap.NewNop(), &fakeUserClient{
		loginFn: func() (*userv1.LoginResponse, error) {
			return &userv1.LoginResponse{UserId: 7, Token: "jwt-token"}, nil
		},
	})

	engine := newUserTestEngine(h)
	body := []byte(`{"email":"john@example.com","password":"password123"}`)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestUserHandler_Login_ValidationError(t *testing.T) {
	h := NewUserHandler(zap.NewNop(), &fakeUserClient{})

	engine := newUserTestEngine(h)
	body := []byte(`{"email":"not-an-email","password":"short"}`)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestUserHandler_Login_GRPCError(t *testing.T) {
	h := NewUserHandler(zap.NewNop(), &fakeUserClient{
		loginFn: func() (*userv1.LoginResponse, error) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		},
	})

	engine := newUserTestEngine(h)
	body := []byte(`{"email":"john@example.com","password":"password123"}`)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp["code"] != float64(http.StatusUnauthorized) {
		t.Fatalf("unexpected error body: %+v", resp)
	}
}

// fakeUserClient — мок userv1.UserServiceClient.
type fakeUserClient struct {
	registerFn func() (*userv1.RegisterResponse, error)
	loginFn    func() (*userv1.LoginResponse, error)
}

func (f *fakeUserClient) Register(_ context.Context, _ *userv1.RegisterRequest, _ ...grpc.CallOption) (*userv1.RegisterResponse, error) {
	if f.registerFn != nil {
		return f.registerFn()
	}
	return &userv1.RegisterResponse{UserId: 1}, nil
}

func (f *fakeUserClient) Login(_ context.Context, _ *userv1.LoginRequest, _ ...grpc.CallOption) (*userv1.LoginResponse, error) {
	if f.loginFn != nil {
		return f.loginFn()
	}
	return &userv1.LoginResponse{UserId: 1, Token: "token"}, nil
}

func (f *fakeUserClient) GetMe(_ context.Context, _ *userv1.GetMeRequest, _ ...grpc.CallOption) (*userv1.GetMeResponse, error) {
	return &userv1.GetMeResponse{}, nil
}

func (f *fakeUserClient) GetUsers(_ context.Context, _ *userv1.GetUsersRequest, _ ...grpc.CallOption) (*userv1.GetUsersResponse, error) {
	return &userv1.GetUsersResponse{}, nil
}

func (f *fakeUserClient) ValidateToken(_ context.Context, _ *userv1.ValidateTokenRequest, _ ...grpc.CallOption) (*userv1.ValidateTokenResponse, error) {
	return &userv1.ValidateTokenResponse{}, nil
}

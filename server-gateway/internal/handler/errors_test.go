package handler

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHTTPStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "unauthenticated", err: status.Error(codes.Unauthenticated, "no token"), want: http.StatusUnauthorized},
		{name: "permission denied", err: status.Error(codes.PermissionDenied, "denied"), want: http.StatusForbidden},
		{name: "not found", err: status.Error(codes.NotFound, "not found"), want: http.StatusNotFound},
		{name: "invalid argument", err: status.Error(codes.InvalidArgument, "bad"), want: http.StatusBadRequest},
		{name: "already exists", err: status.Error(codes.AlreadyExists, "exists"), want: http.StatusConflict},
		{name: "deadline exceeded", err: status.Error(codes.DeadlineExceeded, "timeout"), want: http.StatusGatewayTimeout},
		{name: "unavailable", err: status.Error(codes.Unavailable, "down"), want: http.StatusServiceUnavailable},
		{name: "plain error", err: &errPlain{}, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := HTTPStatus(tt.err); got != tt.want {
				t.Fatalf("HTTPStatus() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	t.Parallel()

	if got := ErrorMessage(status.Error(codes.NotFound, "chat not found")); got != "chat not found" {
		t.Fatalf("ErrorMessage() = %q, want %q", got, "chat not found")
	}

	if got := ErrorMessage(&errPlain{}); got != "internal server error" {
		t.Fatalf("ErrorMessage() = %q, want %q", got, "internal server error")
	}
}

type errPlain struct{}

func (e *errPlain) Error() string { return "plain" }

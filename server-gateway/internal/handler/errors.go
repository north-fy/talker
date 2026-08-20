package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HTTPStatus(err error) int {
	s, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError
	}

	switch s.Code() {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.InvalidArgument, codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unavailable, codes.Aborted:
		return http.StatusServiceUnavailable
	case codes.OK:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}

func ErrorMessage(err error) string {
	s, ok := status.FromError(err)
	if !ok || s.Message() == "" {
		return "internal server error"
	}

	return s.Message()
}

func WriteError(c *gin.Context, err error) {
	code := HTTPStatus(err)

	c.AbortWithStatusJSON(code, gin.H{
		"error": ErrorMessage(err),
		"code":  code,
	})
}

func WriteClientError(c *gin.Context, httpCode int, err error) {
	msg := "invalid request"
	if err != nil {
		msg = err.Error()
	}

	c.AbortWithStatusJSON(httpCode, gin.H{
		"error": msg,
		"code":  httpCode,
	})
}

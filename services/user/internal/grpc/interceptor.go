package grpc

import (
	"context"
	"strings"

	"github.com/north-fy/talker/services/user/internal/domain/models"
	"github.com/north-fy/talker/services/user/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var AuthFunc = func(ctx context.Context) (context.Context, error) {
	publicMethods := map[string]bool{
		"/user.v1.UserService/Register": true,
		"/user.v1.UserService/Login":    true,
	}

	method, ok := grpc.Method(ctx)
	if ok && publicMethods[method] {
		return ctx, nil
	}

	// Извлекаем metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")

	_, err := utils.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	ctx = context.WithValue(ctx, models.Token, token)
	return ctx, nil
}

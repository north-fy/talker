package grpc

import (
	"context"
	"strings"

	"github.com/north-fy/talker/services/user/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var AuthFunc = func(ctx context.Context) (context.Context, error) {
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

	// Отправляем string токен для дальнейшей работы сервиса ValidateToken
	ctx = context.WithValue(ctx, "token", token)
	return ctx, nil
}

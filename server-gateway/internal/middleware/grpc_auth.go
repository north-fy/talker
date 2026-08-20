package middleware

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	grpcAuthHeader = "authorization"
	grpcAuthPrefix = "Bearer "
)

// WithAuthToken добавляет Bearer-токен в исходящий gRPC контекст.
func WithAuthToken(ctx context.Context, token string) context.Context {
	if token == "" {
		return ctx
	}

	md := metadata.Pairs(grpcAuthHeader, grpcAuthPrefix+token)
	if existing, ok := metadata.FromOutgoingContext(ctx); ok {
		md = metadata.Join(existing, md)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

package app

import (
	"context"

	"github.com/north-fy/talker/services/chat/internal/app/grpc"
	"github.com/north-fy/talker/services/chat/internal/config"
	grpc2 "github.com/north-fy/talker/services/chat/internal/grpc"
	"github.com/north-fy/talker/services/chat/internal/service"
	"github.com/north-fy/talker/services/chat/internal/storage/postgres"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpc.App
}

func New(ctx context.Context, log *zap.Logger, grpcPort int, cfgDB config.PostgresCfg) *App {
	db := postgres.NewStorage(ctx, cfgDB)

	chatStorage := postgres.NewChatStorage(db)
	memberStorage := postgres.NewMemberStorage(db)
	inviteStorage := postgres.NewInviteStorage(db)

	serv := service.NewChatService(log, chatStorage, memberStorage, inviteStorage)

	app := grpc.New(log, serv, grpcPort)

	return &App{
		GRPCSrv: app,
	}
}

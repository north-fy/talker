package app

import (
	"context"

	grpcapp "github.com/north-fy/talker/services/user/internal/app/grpc"
	"github.com/north-fy/talker/services/user/internal/config"
	"github.com/north-fy/talker/services/user/internal/service"
	"github.com/north-fy/talker/services/user/internal/storage/postgres"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(ctx context.Context, log *zap.Logger, grpcPort int, cfgDB config.PostgresCfg) *App {
	db := postgres.NewStorage(ctx, cfgDB)
	serv := service.NewService(log, db)

	app := grpcapp.New(log, serv, grpcPort)

	return &App{
		GRPCSrv: app,
	}
}

package app

import (
	"context"

	"github.com/north-fy/talker/services/message/internal/app/grpc"
	"github.com/north-fy/talker/services/message/internal/config"
	grpc2 "github.com/north-fy/talker/services/message/internal/grpc"
	"github.com/north-fy/talker/services/message/internal/service"
	"github.com/north-fy/talker/services/message/internal/storage/postgres"
	"github.com/north-fy/talker/services/message/internal/storage/redis"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpc.App
}

func New(ctx context.Context, log *zap.Logger, grpcPort int, cfgDB config.PostgresCfg, cfgBus config.RedisCfg) *App {
	db := postgres.NewStorage(ctx, cfgDB)
	bus := redis.NewStorage(ctx, cfgBus)
	serv := RegisterService(log, db, bus)

	// WebSocket
	wsServ := service.NewWebSocketService(log, bus)
	ws := grpc2.NewServerWebSocket(wsServ)

	app := grpc.New(log, serv, *ws, grpcPort)

	return &App{
		GRPCSrv: app,
	}
}

type Service struct {
	*service.MessageFuncService
	*service.ReactionService
	*service.FeatureService
}

func RegisterService(log *zap.Logger, db *postgres.Storage, bus *redis.Storage) *Service {
	msgServ := service.NewMessageFuncService(log, db, bus)
	reactServ := service.NewReactionService(log, db, bus)
	featureServ := service.NewFeatureService(log, db, bus)

	return &Service{
		msgServ,
		reactServ,
		featureServ,
	}
}


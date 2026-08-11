package app

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/services/chat/internal/app/grpc"
	"github.com/north-fy/talker/services/chat/internal/clients"
	"github.com/north-fy/talker/services/chat/internal/config"
	"github.com/north-fy/talker/services/chat/internal/service"
	"github.com/north-fy/talker/services/chat/internal/storage/postgres"
	"github.com/north-fy/talker/services/chat/internal/storage/redis"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpc.App
}

func New(ctx context.Context, log *zap.Logger, grpcPort int, cfgDB config.PostgresCfg, cfgCache config.RedisCfg, userAddr, messageAddr string) *App {
	db := postgres.NewStorage(ctx, cfgDB)
	cache := redis.NewStorage(ctx, cfgCache)

	userClient, err := clients.NewUserClient(userAddr)
	if err != nil {
		panic(err)
	}

	messageClient, err := clients.NewMessageClient(messageAddr)
	if err != nil {
		panic(err)
	}

	serv := registerService(log, db, cache, userClient, messageClient)

	app := grpc.New(log, serv, grpcPort)

	return &App{
		GRPCSrv: app,
	}
}

type Service struct {
	*service.ChatFuncService
	*service.MemberService
	*service.FeatureService
	*service.InternalService
	*service.InviteService
}

func registerService(log *zap.Logger, db *postgres.Storage, cache *redis.Storage, userClient userv1.UserServiceClient, messageClient messagev1.MessageServiceClient) *Service {
	chatServ := service.NewChatFuncService(log, db, cache)
	memberServ := service.NewMemberService(log, db, userClient, cache)
	featureServ := service.NewFeatureService(log, db, messageClient, cache)
	internalServ := service.NewInternalService(log, db)
	inviteServ := service.NewInviteService(log, db, cache)

	return &Service{
		ChatFuncService: chatServ,
		MemberService:   memberServ,
		FeatureService:  featureServ,
		InternalService: internalServ,
		InviteService:   inviteServ,
	}
}

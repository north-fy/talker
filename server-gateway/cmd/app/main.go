package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/north-fy/talker/server-gateway/internal/app"
	"github.com/north-fy/talker/server-gateway/internal/clients"
	"github.com/north-fy/talker/server-gateway/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := &config.Config{}
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	userClient, err := clients.NewUserClient(cfg.UserSrvCfg.Addr)
	if err != nil {
		logger.Fatal("failed to create user client", zap.Error(err))
	}

	messageClient, err := clients.NewMessageClient(cfg.MessageSrvCfg.Addr)
	if err != nil {
		logger.Fatal("failed to create message client", zap.Error(err))
	}

	chatClient, err := clients.NewChatClient(cfg.ChatSrvCfg.Addr)
	if err != nil {
		logger.Fatal("failed to create chat client", zap.Error(err))
	}

	server := app.New(
		context.Background(),
		logger,
		cfg.HTTPCfg.Port,
		cfg.UserSrvCfg.Addr,
		cfg.MessageSrvCfg.Addr,
		cfg.ChatSrvCfg.Addr,
		cfg.JWTSecret,
		userClient,
		messageClient,
		chatClient,
	)
	go server.MustRun()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutting down by signal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Stop(ctx)
}

package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/north-fy/talker/services/user/internal/app"
	"github.com/north-fy/talker/services/user/internal/config"
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

	ctx, cancel := context.WithCancel(context.Background())

	server := app.New(ctx, logger, cfg.GRPCCfg.Port, cfg.PostgresCfg)
	server.GRPCSrv.MustRun()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		select {
		case <-ctx.Done():
			cancel()
			server.GRPCSrv.Stop()
			logger.Error("server stopped", zap.Error(ctx.Err()))
		case <-sigChan:
			cancel()
			server.GRPCSrv.Stop()
			logger.Info("server stopped by sigChan")
		}
	}()
}

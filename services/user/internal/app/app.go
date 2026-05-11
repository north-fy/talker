package app

import (
	"log/slog"

	grpcapp "github.com/north-fy/talker/services/user/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

type ConfigDatabase struct {
	Host     string
	User     string
	Password string
	SSLmode  string
}

func New(log *slog.Logger, grpcPort int, cfgDB ConfigDatabase) *App {
	// TODO: init db + service
	panic("implement me")
}

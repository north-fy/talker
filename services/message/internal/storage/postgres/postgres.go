package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/north-fy/talker/services/message/internal/config"
)

const CountTriesPing = 3

type Storage struct {
	conn *pgxpool.Pool
}

func NewStorage(ctx context.Context, cfg config.PostgresCfg) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		panic(err)
	}

	for range CountTriesPing {
		err = conn.Ping(ctx)
		if err != nil {
			panic(err)
		}
	}

	return &Storage{conn: conn}
}


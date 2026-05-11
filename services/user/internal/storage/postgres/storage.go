package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/north-fy/talker/services/user/internal/config"
	"github.com/north-fy/talker/services/user/internal/domain/models"
)

const CountTriesPing = 3

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(cfg config.PostgresCfg) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	conn, err := pgx.Connect(context.TODO(), url)
	if err != nil {
		panic(err)
	}

	for range CountTriesPing {
		err = conn.Ping(context.TODO())
		if err != nil {
			panic(err)
		}
	}

	return &Storage{conn: conn}
}

func (s *Storage) Close() error {
	return s.conn.Close(context.TODO())
}

func (s *Storage) InsertUser(ctx context.Context, user models.User) (int64, error) {
	panic("implement")
}

func (s *Storage) SelectUser(ctx context.Context, user models.User) (models.Session, error) {
	panic("implement")
}

func (s *Storage) SelectUserByToken(ctx context.Context, token string) (models.User, error) {
	panic("implement")
}

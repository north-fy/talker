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

func NewStorage(ctx context.Context, cfg config.PostgresCfg) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	conn, err := pgx.Connect(ctx, url)
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

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) InsertUser(ctx context.Context, user models.User) (int64, error) {
	query := `
	INSERT INTO users (first_name, last_name, email, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	var id int64
	if err := s.conn.QueryRow(ctx, query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) SelectUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `
	SELECT id, first_name, last_name, email, password_hash
	FROM users
	WHERE email = $1
	`

	var user models.User
	if err := s.conn.QueryRow(ctx, query, email).Scan(&user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

package service

import (
	"context"

	"github.com/north-fy/talker/services/user/internal/domain/models"
	"go.uber.org/zap"
)

type Service struct {
	log     *zap.Logger
	storage Storage
}

type Storage interface {
	InsertUser(ctx context.Context, user models.User) (int64, error)
	SelectUser(ctx context.Context, user models.User) (models.Session, error)
	SelectUserByToken(ctx context.Context, token string) (models.User, error)
}

func NewService(log *zap.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s *Service) Register(ctx context.Context, user models.User) (int64, error) {
	s.log = s.log.With(zap.Any("model.User", user))

	return s.storage.InsertUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, user models.User) (models.Session, error) {
	s.log = s.log.With(zap.Any("model.User", user))

	return s.storage.SelectUser(ctx, user)
}

func (s *Service) GetMe(ctx context.Context, token string) (models.User, error) {
	s.log = s.log.With(zap.String("token", token))

	panic(1)
}

func (s *Service) ValidateToken(ctx context.Context, token string) (bool, error) {
	s.log = s.log.With(zap.String("token", token))

	panic(1)
}

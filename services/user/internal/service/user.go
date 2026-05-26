package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/north-fy/talker/services/user/internal/config"
	"github.com/north-fy/talker/services/user/internal/domain/models"
	"github.com/north-fy/talker/services/user/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	log     *zap.Logger
	storage Storage
}

type Storage interface {
	InsertUser(ctx context.Context, user models.User) (int64, error)
	SelectUserByEmail(ctx context.Context, email string) (models.User, error)
}

func NewService(log *zap.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s *Service) Register(ctx context.Context, user models.User) (int64, error) {
	s.log = s.log.With(zap.Any("model.User", user))
	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		s.log.Error("validation failed", zap.Error(err))
		return 0, err
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), config.CostHash)
	if err != nil {
		s.log.Error("generate hash password failed", zap.Error(err))
		return 0, err
	}

	user.Password = string(hashPass)

	id, err := s.storage.InsertUser(ctx, user)
	if err != nil {
		s.log.Error("insert user failed", zap.Error(err))
		return 0, err
	}

	s.log.Info("user registered", zap.Int64("user_id", id))

	return id, nil
}

func (s *Service) Login(ctx context.Context, user models.UserLogin) (models.Session, error) {
	s.log = s.log.With(zap.Any("model.User", user))

	userResp, err := s.storage.SelectUserByEmail(ctx, user.Email)
	if err != nil {
		s.log.Error("select user failed", zap.Error(err))
		return models.Session{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userResp.Password), []byte(user.Password)); err != nil {
		s.log.Error("passwords do not match", zap.Error(err))
		return models.Session{}, err
	}

	token, err := utils.GenerateToken(userResp, "user")
	if err != nil {
		s.log.Error("token doesn't generate", zap.Error(err))
		return models.Session{}, err
	}

	ctx = context.WithValue(ctx, "token", token)

	return models.Session{
		UID:   userResp.UID,
		Token: token,
	}, nil
}

func (s *Service) GetMe(ctx context.Context, token string) (models.User, error) {
	s.log = s.log.With(zap.String("token", token))

	userToken, ok := ctx.Value("token").(string)
	if !ok {
		s.log.Error("user not authenticated")
		return models.User{}, fmt.Errorf("user not authenticated")
	}

	claims, err := utils.ValidateToken(userToken)
	if err != nil {
		s.log.Error("token is not valid", zap.Error(err))
		return models.User{}, fmt.Errorf("token is not valid - %w", err)
	}

	// TODO: delete token in args GetMe
	_ = token

	return models.User{
		UID:       claims.User.UID,
		FirstName: claims.User.FirstName,
		LastName:  claims.User.LastName,
		Email:     claims.User.Email,
	}, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (bool, error) {
	s.log = s.log.With(zap.String("token", token))

	userToken, ok := ctx.Value("token").(string)
	if !ok {
		s.log.Error("user not authenticated")
		return false, fmt.Errorf("user not authenticated")
	}

	_, err := utils.ValidateToken(userToken)
	if err != nil {
		s.log.Error("token is not valid", zap.Error(err))
		return false, fmt.Errorf("token is not valid - %w", err)
	}

	// TODO: delete token in args GetMe
	_ = token

	return true, nil
}

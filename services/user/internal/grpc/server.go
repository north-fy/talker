package grpc

import (
	"context"

	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/services/user/internal/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	Register(ctx context.Context, user models.User) (int64, error)
	Login(ctx context.Context, user models.UserLogin) (models.Session, error)
	GetMe(ctx context.Context, token string) (models.User, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type serverAPI struct {
	userv1.UnimplementedUserServiceServer
	user UserService
}

func Register(gRPC *grpc.Server, service UserService) {
	userv1.RegisterUserServiceServer(gRPC, &serverAPI{user: service})
}

func (s *serverAPI) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	// TODO: обработка ошибок err

	user := models.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
	}

	id, err := s.user.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return &userv1.RegisterResponse{UserId: id}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	// TODO: обработка ошибок err

	user := models.UserLogin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	session, err := s.user.Login(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return &userv1.LoginResponse{UserId: session.UID, Token: session.Token}, nil
}

func (s *serverAPI) GetMe(ctx context.Context, req *userv1.GetMeRequest) (*userv1.GetMeResponse, error) {
	// TODO: обработка ошибок err

	token := req.GetToken()

	userInfo, err := s.user.GetMe(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return &userv1.GetMeResponse{
		UserId:    userInfo.UID,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Email:     userInfo.Email,
	}, nil
}

func (s *serverAPI) ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	// TODO: обработка ошибок err
	token := req.GetToken()

	isValid, err := s.user.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	return &userv1.ValidateTokenResponse{IsValid: isValid}, nil
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/server-gateway/internal/middleware"
	"go.uber.org/zap"
)

type UserHandler struct {
	log    *zap.Logger
	client userv1.UserServiceClient
}

func NewUserHandler(log *zap.Logger, client userv1.UserServiceClient) *UserHandler {
	return &UserHandler{log: log, client: client}
}

type registerRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.Register(c.Request.Context(), &userv1.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		h.log.Error("register failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.Login(c.Request.Context(), &userv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.log.Error("login failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": resp.UserId,
		"token":   resp.Token,
	})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	token, _ := middleware.Token(c)

	resp, err := h.client.GetMe(c.Request.Context(), &userv1.GetMeRequest{
		Token: token,
	})
	if err != nil {
		h.log.Error("get me failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    resp.UserId,
		"first_name": resp.FirstName,
		"last_name":  resp.LastName,
		"email":      resp.Email,
	})
}

func (h *UserHandler) ValidateToken(c *gin.Context) {
	token, _ := middleware.Token(c)

	resp, err := h.client.ValidateToken(c.Request.Context(), &userv1.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		h.log.Error("validate token failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"is_valid": resp.IsValid})
}

package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	userv1 "github.com/north-fy/talker/pkg/protos/user"
	"github.com/north-fy/talker/server-gateway/internal/handler"
	"github.com/north-fy/talker/server-gateway/internal/middleware"
	"go.uber.org/zap"
)

type App struct {
	log        *zap.Logger
	httpServer *http.Server
}

func New(
	ctx context.Context,
	log *zap.Logger,
	httpPort int,
	userAddr, messageAddr, chatAddr, jwtSecret string,
	userClient userv1.UserServiceClient,
	messageClient messagev1.MessageServiceClient,
	chatClient chatv1.ChatServiceClient,
) *App {
	engine := newRouter(log, userAddr, messageAddr, chatAddr, jwtSecret, userClient, messageClient, chatClient)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", httpPort),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &App{
		log:        log,
		httpServer: server,
	}
}

func newRouter(
	log *zap.Logger,
	userAddr, messageAddr, chatAddr, jwtSecret string,
	userClient userv1.UserServiceClient,
	messageClient messagev1.MessageServiceClient,
	chatClient chatv1.ChatServiceClient,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(
		middleware.CORS(),
		middleware.Logger(log),
		middleware.Recovery(log),
	)

	engine.GET("/healthz", handler.Health)

	userH := handler.NewUserHandler(log, userClient)
	chatH := handler.NewChatHandler(log, chatClient)
	messageH := handler.NewMessageHandler(log, messageClient)

	api := engine.Group("/api/v1")

	// auth
	api.POST("/auth/register", userH.Register)
	api.POST("/auth/login", userH.Login)

	api.GET("/auth/me", middleware.Auth(jwtSecret), userH.GetMe)
	api.POST("/auth/validate", middleware.Auth(jwtSecret), userH.ValidateToken)

	// chats
	chats := api.Group("/chats", middleware.Auth(jwtSecret))
	chats.POST("", chatH.CreateChat)
	chats.GET("", chatH.GetChats)
	chats.GET("/:chat_id", chatH.GetChat)
	chats.PATCH("/:chat_id", chatH.UpdateChat)

	// members
	chats.POST("/:chat_id/members", chatH.AddMember)
	chats.GET("/:chat_id/members", chatH.GetMembers)
	chats.DELETE("/:chat_id/members/:user_id", chatH.RemoveMember)
	chats.PATCH("/:chat_id/members/:user_id/role", chatH.UpdateMemberRole)
	chats.GET("/:chat_id/members/:user_id", chatH.GetMember)

	// invites
	chats.POST("/:chat_id/invites", chatH.CreateInviteLink)
	chats.DELETE("/:chat_id/invites/:invite_id", chatH.RevokeInviteLink)

	api.GET("/users/me/chats", middleware.Auth(jwtSecret), chatH.GetUserChats)
	api.POST("/invites/join", middleware.Auth(jwtSecret), chatH.JoinChatByInvite)

	// messages
	messages := api.Group("/messages", middleware.Auth(jwtSecret))
	messages.POST("", messageH.SendMessage)
	messages.GET("/search", messageH.SearchMessages)
	messages.GET("/:message_id", messageH.GetMessage)
	messages.PATCH("/:message_id", messageH.EditMessage)
	messages.DELETE("/:message_id", messageH.DeleteMessage)
	messages.POST("/:message_id/reactions", messageH.AddReaction)
	messages.DELETE("/:message_id/reactions/:reaction", messageH.RemoveReaction)

	chatMessages := api.Group("/chats/:chat_id/messages", middleware.Auth(jwtSecret))
	chatMessages.GET("", messageH.GetMessages)
	chatMessages.POST("", messageH.SendMessage)

	chatRead := api.Group("/chats/:chat_id", middleware.Auth(jwtSecret))
	chatRead.POST("/read", messageH.MarkAsRead)
	chatRead.GET("/unread-count", messageH.GetUnreadCount)

	return engine
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "restapi.Run"

	log := a.log.With(zap.String("op", op))

	log.Info("HTTP server is running", zap.String("addr", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {
	const op = "restapi.Stop"

	a.log.With(zap.String("op", op)).Info("stopping HTTP server")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to shutdown HTTP server", zap.Error(err))
	}
}

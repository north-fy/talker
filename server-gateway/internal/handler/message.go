package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/server-gateway/internal/middleware"
	"go.uber.org/zap"
)

type MessageHandler struct {
	log    *zap.Logger
	client messagev1.MessageServiceClient
}

func NewMessageHandler(log *zap.Logger, client messagev1.MessageServiceClient) *MessageHandler {
	return &MessageHandler{log: log, client: client}
}

type sendMessageRequest struct {
	ChatID      int64                 `json:"chat_id" binding:"required"`
	Content     string                `json:"content" binding:"required"`
	Type        messagev1.MessageType `json:"type"`
	ReplyTo     int64                 `json:"reply_to"`
	Attachments []int64               `json:"attachments"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.SendMessage(c.Request.Context(), &messagev1.SendMessageRequest{
		ChatId:      req.ChatID,
		Content:     req.Content,
		Type:        req.Type,
		ReplyTo:     req.ReplyTo,
		Attachments: req.Attachments,
	})
	if err != nil {
		h.log.Error("send message failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	req := &messagev1.GetMessagesRequest{ChatId: chatID}

	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		req.Limit = int32(v)
	}
	if v, err := strconv.ParseInt(c.Query("before"), 10, 64); err == nil {
		req.Before = v
	}
	if v, err := strconv.ParseInt(c.Query("after"), 10, 64); err == nil {
		req.After = v
	}

	resp, err := h.client.GetMessages(c.Request.Context(), req)
	if err != nil {
		h.log.Error("get messages failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MessageHandler) GetMessage(c *gin.Context) {
	messageID, err := strconv.ParseInt(c.Param("message_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.GetMessage(c.Request.Context(), &messagev1.GetMessageRequest{
		MessageId: messageID,
	})
	if err != nil {
		h.log.Error("get message failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

type editMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *MessageHandler) EditMessage(c *gin.Context) {
	messageID, err := strconv.ParseInt(c.Param("message_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req editMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.EditMessage(c.Request.Context(), &messagev1.EditMessageRequest{
		MessageId: messageID,
		Content:   req.Content,
	})
	if err != nil {
		h.log.Error("edit message failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

type deleteMessageRequest struct {
	ForEveryone bool  `json:"for_everyone"`
	ChatID      int64 `json:"chat_id"`
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	messageID, err := strconv.ParseInt(c.Param("message_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req deleteMessageRequest
	_ = c.ShouldBindJSON(&req)

	_, err = h.client.DeleteMessage(c.Request.Context(), &messagev1.DeleteMessageRequest{
		MessageId:   messageID,
		ForEveryone: req.ForEveryone,
		ChatId:      req.ChatID,
	})
	if err != nil {
		h.log.Error("delete message failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type addReactionRequest struct {
	Reaction string `json:"reaction" binding:"required"`
}

func (h *MessageHandler) AddReaction(c *gin.Context) {
	messageID, err := strconv.ParseInt(c.Param("message_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req addReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.AddReaction(c.Request.Context(), &messagev1.AddReactionRequest{
		MessageId: messageID,
		Reaction:  req.Reaction,
	})
	if err != nil {
		h.log.Error("add reaction failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MessageHandler) RemoveReaction(c *gin.Context) {
	messageID, err := strconv.ParseInt(c.Param("message_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	reaction := c.Param("reaction")

	_, err = h.client.RemoveReaction(c.Request.Context(), &messagev1.RemoveReactionRequest{
		MessageId: messageID,
		Reaction:  reaction,
	})
	if err != nil {
		h.log.Error("remove reaction failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *MessageHandler) SearchMessages(c *gin.Context) {
	req := &messagev1.SearchMessagesRequest{
		Query: c.Query("query"),
	}

	if v, err := strconv.ParseInt(c.Query("chat_id"), 10, 64); err == nil {
		req.ChatId = v
	}
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		req.Limit = int32(v)
	}
	if v, err := strconv.ParseInt(c.Query("before"), 10, 64); err == nil {
		req.Before = v
	}

	resp, err := h.client.SearchMessages(c.Request.Context(), req)
	if err != nil {
		h.log.Error("search messages failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

type markAsReadRequest struct {
	UpToMessageID int64 `json:"up_to_message_id" binding:"required"`
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req markAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	userID, _ := middleware.UserID(c)

	_, err = h.client.MarkAsRead(c.Request.Context(), &messagev1.MarkAsReadRequest{
		ChatId:        chatID,
		UserId:        userID,
		UpToMessageId: req.UpToMessageID,
	})
	if err != nil {
		h.log.Error("mark as read failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	userID, _ := middleware.UserID(c)

	resp, err := h.client.GetUnreadCount(c.Request.Context(), &messagev1.GetUnreadCountRequest{
		ChatId: chatID,
		UserId: userID,
	})
	if err != nil {
		h.log.Error("get unread count failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

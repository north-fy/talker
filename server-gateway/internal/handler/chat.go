package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	chatv1 "github.com/north-fy/talker/pkg/protos/chat"
	"github.com/north-fy/talker/server-gateway/internal/middleware"
	"go.uber.org/zap"
)

type ChatHandler struct {
	log    *zap.Logger
	client chatv1.ChatServiceClient
}

func NewChatHandler(log *zap.Logger, client chatv1.ChatServiceClient) *ChatHandler {
	return &ChatHandler{log: log, client: client}
}

// ==================== CHATS ====================

type createChatRequest struct {
	Name         string          `json:"name" binding:"required,max=255"`
	Type         chatv1.ChatType `json:"type" binding:"required"`
	MemberIDs    []int64         `json:"member_ids" binding:"required,min=1"`
	AvatarBase64 string          `json:"avatar_base64"`
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req createChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	if userID, ok := middleware.UserID(c); ok {
		req.MemberIDs = append([]int64{userID}, req.MemberIDs...)
	}

	resp, err := h.client.CreateChat(c.Request.Context(), &chatv1.CreateChatRequest{
		Name:         req.Name,
		Type:         req.Type,
		MemberIds:    req.MemberIDs,
		AvatarBase64: req.AvatarBase64,
	})
	if err != nil {
		h.log.Error("create chat failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.GetChat(c.Request.Context(), &chatv1.GetChatRequest{
		ChatId: chatID,
	})
	if err != nil {
		h.log.Error("get chat failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	filter := &chatv1.ChatFilter{
		Search: c.Query("search"),
	}

	if t, err := strconv.Atoi(c.Query("type")); err == nil {
		filter.Type = chatv1.ChatType(t)
	}

	if includeArchived, err := strconv.ParseBool(c.Query("include_archived")); err == nil {
		filter.IncludeArchived = includeArchived
	}

	resp, err := h.client.GetChats(c.Request.Context(), &chatv1.GetChatsRequest{Filter: filter})
	if err != nil {
		h.log.Error("get chats failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

type updateChatRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=255"`
	AvatarBase64 *string `json:"avatar_base64"`
}

func (h *ChatHandler) UpdateChat(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req updateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var name, avatar string
	if req.Name != nil {
		name = *req.Name
	}
	if req.AvatarBase64 != nil {
		avatar = *req.AvatarBase64
	}

	resp, err := h.client.UpdateChat(c.Request.Context(), &chatv1.UpdateChatRequest{
		ChatId:       chatID,
		Name:         &name,
		AvatarBase64: &avatar,
	})
	if err != nil {
		h.log.Error("update chat failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) GetUserChats(c *gin.Context) {
	userID, _ := middleware.UserID(c)

	resp, err := h.client.GetUserChats(c.Request.Context(), &chatv1.GetUserChatsRequest{
		UserId: userID,
	})
	if err != nil {
		h.log.Error("get user chats failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== MEMBERS ====================

type addMemberRequest struct {
	UserID int64       `json:"user_id" binding:"required"`
	Role   chatv1.Role `json:"role" binding:"required"`
}

func (h *ChatHandler) AddMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req addMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	invitedBy, _ := middleware.UserID(c)

	resp, err := h.client.AddMember(c.Request.Context(), &chatv1.AddMemberRequest{
		ChatId:    chatID,
		UserId:    req.UserID,
		Role:      req.Role,
		InvitedBy: invitedBy,
	})
	if err != nil {
		h.log.Error("add member failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) RemoveMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	_, err = h.client.RemoveMember(c.Request.Context(), &chatv1.RemoveMemberRequest{
		ChatId: chatID,
		UserId: userID,
	})
	if err != nil {
		h.log.Error("remove member failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ChatHandler) GetMembers(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	filter := &chatv1.MemberFilter{Search: c.Query("search")}
	if r, err := strconv.Atoi(c.Query("role")); err == nil {
		filter.Role = chatv1.Role(r)
	}

	resp, err := h.client.GetMembers(c.Request.Context(), &chatv1.GetMembersRequest{
		ChatId: chatID,
		Filter: filter,
	})
	if err != nil {
		h.log.Error("get members failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

type updateMemberRoleRequest struct {
	Role chatv1.Role `json:"role" binding:"required"`
}

func (h *ChatHandler) UpdateMemberRole(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req updateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.UpdateMemberRole(c.Request.Context(), &chatv1.UpdateMemberRoleRequest{
		ChatId: chatID,
		UserId: userID,
		Role:   req.Role,
	})
	if err != nil {
		h.log.Error("update member role failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) GetMember(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.GetMember(c.Request.Context(), &chatv1.GetMemberRequest{
		ChatId: chatID,
		UserId: userID,
	})
	if err != nil {
		h.log.Error("get member failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== INVITES ====================

type createInviteLinkRequest struct {
	MaxUses int32 `json:"max_uses"`
}

func (h *ChatHandler) CreateInviteLink(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	var req createInviteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	resp, err := h.client.CreateInviteLink(c.Request.Context(), &chatv1.CreateInviteLinkRequest{
		ChatId:  chatID,
		MaxUses: req.MaxUses,
	})
	if err != nil {
		h.log.Error("create invite link failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

type joinChatByInviteRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

func (h *ChatHandler) JoinChatByInvite(c *gin.Context) {
	var req joinChatByInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	userID, _ := middleware.UserID(c)

	resp, err := h.client.JoinChatByInvite(c.Request.Context(), &chatv1.JoinChatByInviteRequest{
		InviteCode: req.InviteCode,
		UserId:     userID,
	})
	if err != nil {
		h.log.Error("join chat by invite failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) RevokeInviteLink(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Param("chat_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	inviteID, err := strconv.ParseInt(c.Param("invite_id"), 10, 64)
	if err != nil {
		WriteClientError(c, http.StatusBadRequest, err)

		return
	}

	_, err = h.client.RevokeInviteLink(c.Request.Context(), &chatv1.RevokeInviteLinkRequest{
		ChatId:   chatID,
		InviteId: inviteID,
	})
	if err != nil {
		h.log.Error("revoke invite link failed", zap.Error(err))
		WriteError(c, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

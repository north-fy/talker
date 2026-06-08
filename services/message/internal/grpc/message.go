package grpc

import (
	"context"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/dto"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"github.com/north-fy/talker/services/message/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const workers = 5

type MessageFuncService interface {
	SendMessage(ctx context.Context, req dto.SendMessageRequest) (models.Message, error)
	GetMessages(ctx context.Context, req dto.GetMessagesRequest) (dto.GetMessagesResponse, error)
	EditMessage(ctx context.Context, req dto.EditMessageRequest) (models.Message, error)
	DeleteMessage(ctx context.Context, req dto.DeleteMessageRequest) (bool, error)
	GetMessage(ctx context.Context, req dto.GetMessageRequest) (models.Message, error)
}

func (s *serverAPI) SendMessage(ctx context.Context, req *messagev1.SendMessageRequest) (*messagev1.Message, error) {
	msgType := dto.MessageType(req.Type.Number())

	msgReq := dto.SendMessageRequest{
		ChatID:      req.GetChatId(),
		Content:     req.GetContent(),
		MessageType: msgType,
		ReplyTo:     req.GetReplyTo(),
		Attachments: req.GetAttachments(),
	}

	message, err := s.msg.SendMessage(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid argument")
	}

	ReplyInfoTmp := &messagev1.ReplyInfo{
		MessageId:      message.ReplyInfoMsg.MessageID,
		SenderName:     message.ReplyInfoMsg.SenderName,
		ContentPreview: message.ReplyInfoMsg.ContentPreview,
	}

	// convert msgType to messagev1.MessageType
	var msgTypes messagev1.MessageType
	tmpType := int32(msgType)

	enumDesc := messagev1.MessageType(0).Descriptor()
	if enumValue := enumDesc.Values().ByNumber(protoreflect.EnumNumber(tmpType)); enumValue != nil {
		msgTypes = messagev1.MessageType(enumValue.Number())
	} else {
		return nil, status.Error(codes.InvalidArgument, "invalid message type")
	}

	return &messagev1.Message{
		Id:          message.ID,
		ChatId:      message.ChatID,
		SenderId:    message.SenderID,
		Content:     message.Content,
		Type:        msgTypes,
		CreatedAt:   timestamppb.New(message.CreatedAt),
		UpdatedAt:   timestamppb.New(message.UpdatedAt),
		IsEdited:    message.IsEdited,
		IsDeleted:   message.IsDeleted,
		ReplyTo:     message.ReplyTo,
		ReplyInfo:   ReplyInfoTmp,
		Attachments: message.Attachments,
		Reactions:   message.Reactions,
	}, nil
}

func (s *serverAPI) GetMessages(ctx context.Context, req *messagev1.GetMessagesRequest) (*messagev1.GetMessagesResponse, error) {
	msgReq := dto.GetMessagesRequest{
		ChatID: req.GetChatId(),
		Limit:  req.GetLimit(),
		Before: req.GetBefore(),
		After:  req.GetAfter(),
	}

	msgResp, err := s.msg.GetMessages(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	messages := utils.ParallelConvert(msgResp.Messages, workers)

	return &messagev1.GetMessagesResponse{
		Messages:   messages,
		HasMore:    msgResp.HasMore,
		TotalCount: msgResp.TotalCount,
	}, nil
}

func (s *serverAPI) EditMessage(ctx context.Context, req *messagev1.EditMessageRequest) (*messagev1.Message, error) {
	msgReq := dto.EditMessageRequest{
		MessageID: req.GetMessageId(),
		Content:   req.GetContent(),
	}

	msg, err := s.msg.EditMessage(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return utils.ConvertToProtoMessage(&msg), nil
}

func (s *serverAPI) DeleteMessage(ctx context.Context, req *messagev1.DeleteMessageRequest) (*messagev1.Empty, error) {
	msgReq := dto.DeleteMessageRequest{
		MessageID:   req.GetMessageId(),
		ForEveryone: req.GetForEveryone(),
	}

	isDeleted, err := s.msg.DeleteMessage(ctx, msgReq)
	if err != nil || !isDeleted {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return nil, nil
}

func (s *serverAPI) GetMessage(ctx context.Context, req *messagev1.GetMessageRequest) (*messagev1.Message, error) {
	msgReq := dto.GetMessageRequest{
		MessageID: req.GetMessageId(),
	}

	msg, err := s.msg.GetMessage(ctx, msgReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return utils.ConvertToProtoMessage(&msg), nil
}

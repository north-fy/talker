package utils

import (
	"strconv"
	"sync"

	messagev1 "github.com/north-fy/talker/pkg/protos/message"
	"github.com/north-fy/talker/services/message/internal/domain/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertMessageType convert dto.MessageType to messagev1.MessageType
func ConvertMessageType(messageType string) messagev1.MessageType {
	tmpType, _ := strconv.Atoi(messageType)
	return messagev1.MessageType(tmpType)
}

// ConvertToProtoMessage convert models.Message to messagev1.Message
func ConvertToProtoMessage(msg *models.Message) *messagev1.Message {
	return &messagev1.Message{
		Id:        msg.ID,
		ChatId:    msg.ChatID,
		SenderId:  msg.SenderID,
		Content:   msg.Content,
		Type:      ConvertMessageType(msg.MessageType),
		CreatedAt: timestamppb.New(msg.CreatedAt),
		UpdatedAt: timestamppb.New(msg.UpdatedAt),
		IsEdited:  msg.IsEdited,
		IsDeleted: msg.IsDeleted,
		ReplyTo:   msg.ReplyTo,
		ReplyInfo: &messagev1.ReplyInfo{
			MessageId:      msg.ReplyInfoMsg.MessageID,
			SenderName:     msg.ReplyInfoMsg.SenderName,
			ContentPreview: msg.ReplyInfoMsg.ContentPreview,
		},
		Attachments: msg.Attachments,
		Reactions:   msg.Reactions,
	}
}

// ParallelConvert convert []*models.Message to []*messagev1.Message
func ParallelConvert(messages []*models.Message, workers int) []*messagev1.Message {
	wg := &sync.WaitGroup{}
	sema := make(chan struct{}, workers)
	res := make([]*messagev1.Message, len(messages))

	for i, msg := range messages {

		wg.Add(1)
		go func(i int, msg *models.Message) {
			defer wg.Done()
			sema <- struct{}{}
			res[i] = ConvertToProtoMessage(msg)
			<-sema
		}(i, msg)
	}

	wg.Wait()
	return res
}

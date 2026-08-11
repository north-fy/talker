package service

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/north-fy/talker/services/chat/internal/domain"
)

var Validator = validator.New()

// validateStruct валидирует структуру и маппит ошибки на доменные.
func validateStruct(ctx context.Context, req any) error {
	if err := Validator.StructCtx(ctx, req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fe := range validationErrors {
				return mapValidationError(fe)
			}
		}

		return domain.ErrInvalidStruct
	}

	return nil
}

func mapValidationError(fe validator.FieldError) error {
	switch fe.Field() {
	case "Name":
		if fe.Tag() == "required" {
			return domain.ErrChatNameEmpty
		}

		return domain.ErrChatNameTooLong
	case "Role":
		return domain.ErrInvalidRole
	case "RequiredPermission":
		return domain.ErrInvalidPermission
	default:
		return domain.ErrInvalidStruct
	}
}

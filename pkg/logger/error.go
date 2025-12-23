package logger

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type errorDto struct {
	ErrorDetails map[string]any `json:"error_details"`
	Error        string         `json:"error"`
	Description  string         `json:"description"`
	Code         string         `json:"code"`
	Id           string         `json:"id"`
}

type baseError interface {
	Error() string
	Code() string
	ErrorDetails() map[string]any
	Description() string
}

func errorModel(ctx context.Context, err error, messages ...string) *errorDto {
	internalMessage := "Something went wrong, please try again in a few seconds"
	description := ""
	if len(messages) > 0 {
		internalMessage = messages[0]
	}
	if len(messages) > 1 {
		description = messages[1]
	}

	var correlationId string
	correlationIdFromEnv := ctx.Value("correlationId")
	if correlationIdFromEnv != nil && correlationIdFromEnv.(string) != "" {
		correlationId = correlationIdFromEnv.(string)
	} else {
		correlationId = uuid.NewString()
	}

	var e baseError
	switch {
	case err != nil && errors.As(err, &e):
		return &errorDto{
			Id:           correlationId,
			Error:        e.Error(),
			Description:  e.Description(),
			Code:         e.Code(),
			ErrorDetails: e.ErrorDetails(),
		}
	case err != nil:
		return &errorDto{
			Id:           correlationId,
			Error:        err.Error(),
			Description:  internalMessage,
			Code:         "IS-00500",
			ErrorDetails: nil,
		}
	default:
		return &errorDto{
			Id:           correlationId,
			Error:        internalMessage,
			Description:  description,
			Code:         "IS-00500",
			ErrorDetails: nil,
		}
	}
}
